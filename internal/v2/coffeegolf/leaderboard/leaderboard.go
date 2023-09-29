package leaderboard

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ryansheppard/morningjuegos/internal/v2/cache"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
)

type Leaderboard struct {
	ctx   context.Context
	query *database.Queries
	cache *cache.Cache
}

func New(ctx context.Context, query *database.Queries, cache *cache.Cache) *Leaderboard {
	return &Leaderboard{
		ctx:   ctx,
		query: query,
		cache: cache,
	}
}

func (l *Leaderboard) getLeaderboardCacheKey(guildID int64) string {
	return fmt.Sprintf("leaderboard:%d", guildID)
}

// TODO: Chunk this function up
func (l *Leaderboard) GenerateLeaderboard(guildIDString string) string {
	guildID, err := strconv.ParseInt(guildIDString, 10, 64)
	if err != nil {
		slog.Error("Failed to parse guildID", "guild", guildIDString, "error", err)
		return "Could not generate a leaderboard for this discord server"
	}

	slog.Info("Generating leaderboard", "guild", guildID)

	tournament, err := l.query.GetActiveTournament(l.ctx, guildID)
	if err == sql.ErrNoRows {
		return "Could not find a tournament for this discord server"
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return "Error getting a tournament for this discord server"
	}

	cacheKey := l.getLeaderboardCacheKey(guildID)
	var cached interface{}
	if l.cache != nil {
		cached, err = l.cache.GetKey(cacheKey)
		if err != nil {
			slog.Error("Failed to get leaderboard from cache", "guild", guildID, "error", err)
		}
	}

	// TODO figure out what happens if we get an error before this
	if cached != nil {
		slog.Info("Returning cached leaderboard", "guild", guildID)
		return cached.(string)
	}

	startDate := tournament.StartTime.Format("Jan 2, 2006")
	endDate := tournament.EndTime.Format("Jan 2, 2006")

	tournamentString := fmt.Sprintf("Current Tournament: %s - %s", startDate, endDate)

	strokeLeaders, err := l.query.GetLeaders(l.ctx, tournament.ID)
	if err != nil {
		slog.Error("Failed to get leaders", "guild", guildID, "error", err)
	}

	if len(strokeLeaders) == 0 {
		return "No one has played yet!"
	}

	leaderStrings := []string{}
	notYetPlayed := []string{}
	skipCounter := 0
	for i, leader := range strokeLeaders {
		_, err := l.query.HasPlayedToday(l.ctx, database.HasPlayedTodayParams{
			PlayerID:     leader.PlayerID,
			TournamentID: tournament.ID,
		})
		hasPlayed := true
		if err == sql.ErrNoRows {
			hasPlayed = false
		} else if err != nil {
			slog.Error("Failed to check if player has played today", "guild", guildID, "player", leader.PlayerID, "error", err)
			continue
		}

		if hasPlayed {
			leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%d> - %d Total Strokes", i+1-skipCounter, leader.PlayerID, leader.TotalStrokes))
		} else {
			notYetPlayed = append(notYetPlayed, fmt.Sprintf("<@%d> - %d Total Strokes", leader.PlayerID, leader.TotalStrokes))
			skipCounter++
		}
	}

	leaderString := "Leaders\n" + strings.Join(leaderStrings, "\n")

	notYetPlayedString := ""
	if len(notYetPlayed) > 0 {
		rand.Shuffle(len(notYetPlayed), func(i, j int) { notYetPlayed[i], notYetPlayed[j] = notYetPlayed[j], notYetPlayed[i] })

		notYetPlayedString = "\nNot Played Yet\n" + strings.Join(notYetPlayed, "\n") + "\n"
	}

	hole, err := l.query.GetHardestHole(l.ctx, tournament.ID)
	if err != nil {
		slog.Error("Failed to get hardest hole", "tournament", tournament.ID, "error", err)
	}
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %0.2f strokes", hole.Color, hole.Strokes)

	// TODO: handle errors here
	firstMost, _ := l.query.GetMostCommonHoleForNumber(l.ctx, database.GetMostCommonHoleForNumberParams{
		TournamentID: tournament.ID,
		HoleNumber:   0,
	})
	lastMost, _ := l.query.GetMostCommonHoleForNumber(l.ctx, database.GetMostCommonHoleForNumberParams{
		TournamentID: tournament.ID,
		HoleNumber:   4,
	})
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost.Color, lastMost.Color)

	worstRound, _ := l.query.GetWorstRound(l.ctx, tournament.ID)
	worstRoundString := fmt.Sprintf("The worst round was %d strokes by <@%d>", worstRound.TotalStrokes, worstRound.PlayerID)

	statsStr := "\n\nStats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString + "\n" + worstRoundString

	all := tournamentString + "\n\n" + leaderString + notYetPlayedString + statsStr

	if l.cache != nil {
		l.cache.SetKey(cacheKey, all, 3600)
	}

	return all
}
