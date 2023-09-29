package leaderboard

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/cache"
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
func (l *Leaderboard) generateLeaderboard(guildID int64) string {
	tournament, err := l.query.GetActiveTournament(l.ctx, guildID)
	if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, err)
		return "Could not find a tournament for this discord server"
	}
	if tournament == (database.Tournament{}) {
		return "No active tournament"
	}

	cacheKey := l.getLeaderboardCacheKey(guildID)
	cached, err := l.cache.GetKey(cacheKey)
	if err != nil {
		slog.Error("Failed to get leaderboard from cache", "guild", guildID, err)
	}

	// TODO figure out what happens if we get an error before this
	if cached != nil {
		return cached.(string)
	}

	startDate := tournament.StartTime.Format("Jan 2, 2006")
	endDate := tournament.EndTime.Format("Jan 2, 2006")

	tournamentString := fmt.Sprintf("Current Tournament: %s - %s", startDate, endDate)

	strokeLeaders, err := l.query.GetLeaders(l.ctx, tournament.ID)
	if err != nil {
		slog.Error("Failed to get leaders", "guild", guildID, err)
	}

	if len(strokeLeaders) == 0 {
		return "No one has played yet!"
	}

	leaderStrings := []string{}
	notYetPlayed := []string{}
	skipCounter := 0
	for i, leader := range strokeLeaders {
		hasPlayedToday, err := l.query.HasPlayedToday(l.ctx, database.HasPlayedTodayParams{
			PlayerID:     leader.PlayerID,
			TournamentID: tournament.ID,
		})
		if err != nil {
			slog.Error("Failed to check if player has played today", "guild", guildID, "player", leader.PlayerID, err)
		}
		if hasPlayedToday != (database.Round{}) {
			leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%d> - %d Total Strokes", i+1-skipCounter, leader.PlayerID, leader.Strokes))
		} else {
			notYetPlayed = append(notYetPlayed, fmt.Sprintf("<@%d> - %d Total Strokes", leader.PlayerID, leader.Strokes))
			skipCounter++
		}
	}

	leaderString := "Leaders\n" + strings.Join(leaderStrings, "\n") + "\n"

	notYetPlayedString := ""
	if len(notYetPlayed) > 0 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(notYetPlayed), func(i, j int) { notYetPlayed[i], notYetPlayed[j] = notYetPlayed[j], notYetPlayed[i] })

		notYetPlayedString = "Not Played Yet\n" + strings.Join(notYetPlayed, "\n")
	}

	hole, err := l.query.GetHardestHole(l.ctx, tournament.ID)
	if err != nil {
		slog.Error("Failed to get hardest hole", "tournament", tournament.ID, err)
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

	statsStr := "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString + "\n" + worstRoundString

	all := tournamentString + "\n\n" + leaderString + "\n" + notYetPlayedString + "\n" + statsStr

	l.cache.SetKey(cacheKey, all, 3600)

	return all
}
