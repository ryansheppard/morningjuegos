package leaderboard

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

var (
	colors = map[string]string{
		"blue":   "🟦",
		"red":    "🟥",
		"yellow": "🟨",
		"green":  "🟩",
		"purple": "🟪",
	}
	placements = map[int]string{
		1: "🥇",
		2: "🥈",
		3: "🥉",
	}
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

type GenerateLeaderboardParams struct {
	GuildID string
	Date    string
}

func (p *GenerateLeaderboardParams) SetDate(date string) {
	p.Date = date
}

func (l *Leaderboard) GenerateLeaderboard(params GenerateLeaderboardParams) string {
	// todo: find better way to do this
	startTime := time.Now().Add(11 * -24 * time.Hour)
	endTime := time.Now().Add(24 * time.Hour)
	// includeEmoji gets set false when parsing a date specific leaderboard, so it gets used in a lot of places
	includeEmoji := true
	if params.Date != "" {
		newYork, err := time.LoadLocation("America/New_York")
		if err != nil {
			slog.Error("Failed to load location", "error", err)
			return "Could not parse the given date, try using the format 2023-01-01 (yyyy-mm-dd)"
		}
		parsedTime, err := time.ParseInLocation("2006-01-02", params.Date, newYork)
		if err != nil {
			slog.Error("Failed to parse date", "date", params.Date, "error", err)
			return "Could not parse the given date, try using the format 2023-01-01 (yyyy-mm-dd)"
		}

		startTime = parsedTime
		endTime = parsedTime.Add(24 * time.Hour).Add(-1 * time.Second)
		includeEmoji = false
	}

	guildID, err := strconv.ParseInt(params.GuildID, 10, 64)
	if err != nil {
		slog.Error("Failed to parse guildID", "guild", params.GuildID, "error", err)
		return "Could not generate a leaderboard for this discord server"
	}

	slog.Info("Generating leaderboard", "guild", guildID, "startTime", startTime, "endTime", endTime)

	tournamentTime := time.Now()
	if params.Date != "" {
		tournamentTime = startTime
	}

	tournament, err := l.query.GetActiveTournament(l.ctx, database.GetActiveTournamentParams{
		GuildID:   guildID,
		StartTime: tournamentTime,
	})
	if err == sql.ErrNoRows {
		return "Could not find a tournament for this discord server"
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return "Error getting a tournament for this discord server"
	}

	var cached interface{}
	cacheKey := ""
	if includeEmoji {
		cacheKey = GetLeaderboardCacheKey(guildID)
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

	var header string
	if params.Date != "" {
		header = "Leaderboard for " + params.Date
	} else {
		header = l.generateHeader(&tournament)
	}

	leaderParams := generateLeaderStringParams{
		GuildID:      guildID,
		TournamentID: tournament.ID,
		FirstRound:   true,
		StartTime:    startTime,
		EndTime:      endTime,
		IncludeEmoji: includeEmoji,
	}
	leaderString := l.generateLeaderString(leaderParams)

	all := header + "\n\n" + leaderString

	if includeEmoji {
		l.cache.SetKey(cacheKey, all, 3600)
	}

	return all
}

func (l *Leaderboard) GenerateStats(guildIDString string) string {
	guildID, err := strconv.ParseInt(guildIDString, 10, 64)
	if err != nil {
		slog.Error("Failed to parse guildID", "guild", guildIDString, "error", err)
		return "Could not generate stats for this discord server"
	}

	slog.Info("Generating stats", "guild", guildID)

	tournament, err := l.query.GetActiveTournament(l.ctx, database.GetActiveTournamentParams{
		GuildID:   guildID,
		StartTime: time.Now(),
	})
	if err == sql.ErrNoRows {
		return "Could not find a tournament for this discord server"
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return "Error getting a tournament for this discord server"
	}

	cacheKey := GetStatsCacheKey(guildID)
	var cached interface{}
	cached, err = l.cache.GetKey(cacheKey)
	if err != nil {
		slog.Error("Failed to get stats from cache", "guild", guildID, "error", err)
	}

	// TODO figure out what happens if we get an error before this
	if cached != nil {
		slog.Info("Returning cached stats", "guild", guildID)
		return cached.(string)
	}

	header := l.generateHeader(&tournament)

	statsStr := l.generateStats(tournament.ID)

	all := header + "\n\n" + statsStr

	l.cache.SetKey(cacheKey, all, 3600)

	return all
}

func (l *Leaderboard) generateHeader(tournament *database.Tournament) string {
	startDate := tournament.StartTime.Format("Jan 2, 2006")
	endDate := tournament.EndTime.Format("Jan 2, 2006")

	return fmt.Sprintf("Current Tournament: %s - %s", startDate, endDate)

}

type generateLeaderStringParams struct {
	GuildID      int64
	TournamentID int32
	FirstRound   bool
	StartTime    time.Time
	EndTime      time.Time
	IncludeEmoji bool
}

func (l *Leaderboard) generateLeaderString(params generateLeaderStringParams) string {
	notPlayedX := false
	notPlayedXRaw, err := l.cache.GetKey("notPlayedX")
	if err != nil || notPlayedXRaw == nil {
		slog.Info("Failed to get notPlayedX from cache, defaulting to false", "error", err)
	} else {
		notPlayedX = notPlayedXRaw.(string) == "true"
	}

	addPlusTwenty := false
	addPlusRaw, err := l.cache.GetKey("addPlusTwenty")
	if err != nil || addPlusRaw == nil {
		slog.Info("Failed to get addPlusTwenty from cache, defaulting to false", "error", err)
	} else {
		addPlusTwenty = addPlusRaw.(string) == "true"
	}

	slog.Info("Generating leaderboard string", "guild", params.GuildID, "tournament", params)
	strokeLeaders, err := l.query.GetLeaders(l.ctx, database.GetLeadersParams{
		TournamentID: params.TournamentID,
		RoundDate:    sql.NullTime{Time: params.StartTime, Valid: true},
		RoundDate_2:  sql.NullTime{Time: params.EndTime, Valid: true},
	})

	if err != nil {
		slog.Error("Failed to get leaders", "guild", params.GuildID, "error", err)
		return "Failed to get leaderboard"
	}

	if len(strokeLeaders) == 0 {
		slog.Info("No one has played yet", "guild", params.GuildID)
		return "No one has played yet!"
	}

	previousPlacements := map[int64]int{}
	if params.IncludeEmoji {
		previousPlacements = l.getPreviousPlacements(params.GuildID, params.TournamentID)
	}

	leaderStrings := []string{}
	notYetPlayed := []string{}
	skipCounter := 0
	for i, leader := range strokeLeaders {
		_, err := l.query.HasPlayedToday(l.ctx, database.HasPlayedTodayParams{
			PlayerID:     leader.PlayerID,
			TournamentID: params.TournamentID,
		})

		placementString := ""
		if params.IncludeEmoji {
			placementString = l.getPlacementEmoji(i + 1)
		}

		prev := -1
		if previousPlacement, ok := previousPlacements[leader.PlayerID]; ok {
			prev = previousPlacement
		}

		hasPlayed := true
		if err == sql.ErrNoRows {
			hasPlayed = false
		} else if err != nil {
			slog.Error("Failed to check if player has played today", "guild", params.GuildID, "player", leader.PlayerID, "error", err)
			continue
		}

		movement := "❌"
		addedScore := 0
		if hasPlayed {
			movement = l.getPreviousPlacementEmoji(prev, i+1)
			if addPlusTwenty {
				addedScore = 20
			}
		}

		previousWinString := ""
		if params.IncludeEmoji {
			previousWinString = l.getCrowns(params.GuildID, leader.PlayerID)
		}

		if hasPlayed || notPlayedX {
			strokeString := fmt.Sprintf("%d: <@%d> - %d Total Strokes", i+1-skipCounter, leader.PlayerID, leader.TotalStrokes+int64(addedScore))
			finalString := strings.Join([]string{
				strokeString,
				placementString,
				movement,
				previousWinString,
			}, " ")
			leaderStrings = append(leaderStrings, finalString)
		} else {
			notYetPlayed = append(notYetPlayed, fmt.Sprintf("<@%d> - %d Total Strokes %s", leader.PlayerID, leader.TotalStrokes, placementString))
			skipCounter++
		}
	}

	leaderString := "No one has played today!\n"
	if !notPlayedX {
		if len(leaderStrings) > 0 {
			leaderString = "Leaders\n" + strings.Join(leaderStrings, "\n")
		}
	} else {
		leaderString = strings.Join(leaderStrings, "\n")
	}

	notYetPlayedString := ""
	if len(notYetPlayed) > 0 {
		rand.Shuffle(len(notYetPlayed), func(i, j int) { notYetPlayed[i], notYetPlayed[j] = notYetPlayed[j], notYetPlayed[i] })

		notYetPlayedString = "\n\nNot Played Yet\n" + strings.Join(notYetPlayed, "\n") + "\n"
	}

	return leaderString + notYetPlayedString
}

func (l *Leaderboard) getCrowns(guildID int64, playerID int64) string {
	previousWins, err := l.query.GetTournamentPlacementsByPosition(l.ctx, database.GetTournamentPlacementsByPositionParams{
		GuildID:             guildID,
		PlayerID:            playerID,
		TournamentPlacement: 1,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Failed to get previous wins", "guild", guildID, "player", playerID, "error", err)
		return ""
	} else {
		if previousWins.Count > 0 {
			return fmt.Sprintf("%d 👑", previousWins.Count)
		}
	}

	return ""
}

func (l *Leaderboard) getPlacementEmoji(placement int) string {
	if placement, ok := placements[placement]; ok {
		return placement
	}

	return ""
}

func (l *Leaderboard) getPreviousPlacementEmoji(prev int, current int) string {
	if prev != -1 {
		if prev > current {
			return "⬆️"
		} else if prev < current {
			return "⬇️"
		}
	}

	return ""
}

func (l *Leaderboard) getPreviousPlacements(guildID int64, tournamentID int32) map[int64]int {
	previousPlacements, err := l.query.GetPlacementsForPeriod(l.ctx, database.GetPlacementsForPeriodParams{
		TournamentID: tournamentID,
		RoundDate:    sql.NullTime{Time: time.Now().Add(-24 * time.Hour), Valid: true},
	})
	if err == sql.ErrNoRows {
		slog.Warn("No previous placements", "guild", guildID)
	} else if err != nil {
		slog.Error("Failed to get previous placements", "guild", guildID, "error", err)
	}

	previous := make(map[int64]int)
	for i, previousPlacement := range previousPlacements {
		previous[previousPlacement.PlayerID] = i + 1
	}

	return previous
}

func (l *Leaderboard) generateStats(tournamentID int32) string {
	holeInOneString := l.getHoleInOneLeader(tournamentID)
	worstRoundString := l.getWorstRound(tournamentID)
	mostCommon := l.getFirstMostCommonHole(tournamentID)
	lastCommon := l.getLastMostCommonHole(tournamentID)
	hardestHole := l.getHardestHole(tournamentID)
	bestPerformers := l.getBestPerformers()
	worstPerformers := l.getWorstPerformers()

	statsHeader := "Stats powered by AWS Next Gen Stats"
	statsStr := strings.Join([]string{
		statsHeader,
		holeInOneString,
		worstRoundString,
		mostCommon,
		lastCommon,
		hardestHole,
		bestPerformers,
		worstPerformers,
	}, "\n")

	return statsStr
}

func (l *Leaderboard) getHoleInOneLeader(tournamentID int32) string {
	holeInOneLeader, err := l.query.GetHoleInOneLeader(l.ctx, tournamentID)
	if err == sql.ErrNoRows {
		slog.Warn("No hole in one leader", "tournament", tournamentID)
	} else if err != nil {
		slog.Error("Failed to get hole in one leader", "tournament", tournamentID, "error", err)
	} else {
		plural := ""
		if holeInOneLeader.PlayerID.Int64 != 1 {
			plural = "s"
		}
		return fmt.Sprintf("Most hole in ones: <@%d> with %d hole in one%s", holeInOneLeader.PlayerID.Int64, holeInOneLeader.Count, plural)
	}

	return ""
}

func (l *Leaderboard) getWorstRound(tournamentID int32) string {
	worstRound, _ := l.query.GetWorstRound(l.ctx, tournamentID)
	return fmt.Sprintf("Worst round of the tournament: <@%d>, %d strokes 🤡", worstRound.PlayerID, worstRound.TotalStrokes)
}

func (l *Leaderboard) getFirstMostCommonHole(tournamentID int32) string {
	firstMost, _ := l.query.GetMostCommonHoleForNumber(l.ctx, database.GetMostCommonHoleForNumberParams{
		TournamentID: tournamentID,
		HoleNumber:   0,
	})
	return fmt.Sprintf("Most common opening hole: %s", colors[firstMost.Color])

}

func (l *Leaderboard) getLastMostCommonHole(tournamentID int32) string {
	lastMost, _ := l.query.GetMostCommonHoleForNumber(l.ctx, database.GetMostCommonHoleForNumberParams{
		TournamentID: tournamentID,
		HoleNumber:   4,
	})
	return fmt.Sprintf("Most common finishing hole: %s", colors[lastMost.Color])

}

func (l *Leaderboard) getHardestHole(tournamentID int32) string {
	hole, err := l.query.GetHardestHole(l.ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get hardest hole", "tournament", tournamentID, "error", err)
	}
	return fmt.Sprintf("Hardest hole: %s with an average of %0.2f strokes", colors[hole.Color], hole.Strokes)
}

func (l *Leaderboard) getPerformers(reverse bool) ([]string, string, error) {
	performers, err := l.query.GetStandardDeviation(l.ctx)
	if err == sql.ErrNoRows {
		slog.Info("No std dev found")
		return []string{}, "", err
	} else if err != nil {
		slog.Error("Failed to std dev", "error", err)
		return []string{}, "", err
	}

	if reverse {
		for i, j := 0, len(performers)-1; i < j; i, j = i+1, j-1 {
			performers[i], performers[j] = performers[j], performers[i]
		}
	}

	stdDev := ""
	performerIDs := []string{}
	for i, performer := range performers {
		if i == 0 || performer.StandardDeviation == performers[i-1].StandardDeviation {
			performerIDs = append(performerIDs, fmt.Sprintf("<@%d>", performer.PlayerID))
			stdDev = performer.StandardDeviation
		} else {
			break
		}
	}

	return performerIDs, stdDev, nil
}

func (l *Leaderboard) getWorstPerformers() string {
	worstPerformerString := ""

	worstPerformerIDs, stdDev, err := l.getPerformers(true)
	if err != nil {
		return worstPerformerString
	}

	worstPerformerMentions := strings.Join(worstPerformerIDs, ", ")

	return fmt.Sprintf("[All Time] Least consistent players: %v with a standard deviation of %s strokes", worstPerformerMentions, stdDev)
}

func (l *Leaderboard) getBestPerformers() string {
	bestPerformersString := ""

	bestPerformerIDs, stdDev, err := l.getPerformers(false)
	if err != nil {
		return bestPerformersString
	}

	bestPerformersMentions := strings.Join(bestPerformerIDs, ", ")

	return fmt.Sprintf("[All Time] Most consistent players: %v with a standard deviation of %s strokes", bestPerformersMentions, stdDev)
}
