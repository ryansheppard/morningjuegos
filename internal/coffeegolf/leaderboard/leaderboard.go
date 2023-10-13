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
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/service"
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
	service *service.Service
	cache   *cache.Cache
}

func New(service *service.Service, cache *cache.Cache) *Leaderboard {
	return &Leaderboard{
		service: service,
		cache:   cache,
	}
}

func (l *Leaderboard) GenerateLeaderboard(ctx context.Context, guildIDAsString string) string {
	// includeEmoji gets set false when parsing a date specific leaderboard, so it gets used in a lot of places
	includeEmoji := true

	guildID, err := strconv.ParseInt(guildIDAsString, 10, 64)
	if err != nil {
		slog.Error("Failed to parse guildID", "guild", guildIDAsString, "error", err)
		return "Could not generate a leaderboard for this discord server"
	}

	slog.Info("Generating leaderboard", "guild", guildID)

	tournament, err := l.service.GetActiveTournament(ctx, guildID)
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
		cached, err = l.cache.GetKey(ctx, cacheKey)
		if err != nil {
			slog.Error("Failed to get leaderboard from cache", "guild", guildID, "error", err)
		}
	}

	// TODO figure out what happens if we get an error before this
	if cached != nil {
		slog.Info("Returning cached leaderboard", "guild", guildID)
		return cached.(string)
	}

	header := l.generateHeader(tournament)

	leaderParams := generateLeaderStringParams{
		GuildID:      guildID,
		TournamentID: tournament.ID,
		FirstRound:   true,
		IncludeEmoji: includeEmoji,
	}
	leaderString := l.generateLeaderString(ctx, leaderParams)

	all := header + "\n\n" + leaderString

	if includeEmoji {
		l.cache.SetKey(ctx, cacheKey, all, 3600)
	}

	return all
}

func (l *Leaderboard) GenerateStats(ctx context.Context, guildIDString string) string {
	guildID, err := strconv.ParseInt(guildIDString, 10, 64)
	if err != nil {
		slog.Error("Failed to parse guildID", "guild", guildIDString, "error", err)
		return "Could not generate stats for this discord server"
	}

	slog.Info("Generating stats", "guild", guildID)

	tournament, err := l.service.GetActiveTournament(ctx, guildID)
	if err == sql.ErrNoRows {
		return "Could not find a tournament for this discord server"
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return "Error getting a tournament for this discord server"
	}

	cacheKey := GetStatsCacheKey(guildID)
	var cached interface{}
	cached, err = l.cache.GetKey(ctx, cacheKey)
	if err != nil {
		slog.Error("Failed to get stats from cache", "guild", guildID, "error", err)
	}

	// TODO figure out what happens if we get an error before this
	if cached != nil {
		slog.Info("Returning cached stats", "guild", guildID)
		return cached.(string)
	}

	header := l.generateHeader(tournament)

	statsStr := l.generateStats(ctx, tournament.ID)

	all := header + "\n\n" + statsStr

	l.cache.SetKey(ctx, cacheKey, all, 3600)

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
	IncludeEmoji bool
}

func (l *Leaderboard) generateLeaderString(ctx context.Context, params generateLeaderStringParams) string {
	slog.Info("Generating leaderboard string", "guild", params.GuildID, "tournament", params)
	strokeLeaders, err := l.service.GetLeaders(ctx, params.TournamentID)

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
		previousPlacements = l.getPreviousPlacements(ctx, params.GuildID, params.TournamentID)
	}

	previousWins := map[int64]int64{}
	if params.IncludeEmoji {
		previousWins = l.getPreviousWins(ctx, params.GuildID)
	}

	leaderStrings := []string{}
	notYetPlayed := []string{}
	skipCounter := 0
	for i, leader := range strokeLeaders {
		emojiStrings := []string{}

		prev := -1
		if previousPlacement, ok := previousPlacements[leader.PlayerID]; ok {
			prev = previousPlacement
		}

		placementString := ""
		if params.IncludeEmoji && prev > 0 {
			placementString = l.getPlacementEmoji(prev)
		}

		hasPlayed, err := l.service.HasPlayedToday(ctx, leader.PlayerID, params.TournamentID)
		if err != nil {
			slog.Error("Failed to check if player has played today", "guild", params.GuildID, "player", leader.PlayerID, "error", err)
			continue
		}

		movement := ""
		if hasPlayed {
			movement = l.getPreviousPlacementEmoji(prev, i+1)
			emojiStrings = append(emojiStrings, movement)
		}

		previousWinString := ""
		if params.IncludeEmoji {
			previousWins, ok := previousWins[leader.PlayerID]
			if ok {
				if previousWins > 0 {
					previousWinString = fmt.Sprintf("%d 👑", previousWins)
					emojiStrings = append(emojiStrings, previousWinString)
				}
			}
		}

		emojiString := strings.Join(emojiStrings, " ")

		if hasPlayed {
			strokeString := fmt.Sprintf("%d. <@%d> - %d Strokes %s", i+1-skipCounter, leader.PlayerID, leader.TotalStrokes, emojiString)
			leaderStrings = append(leaderStrings, strokeString)
		} else {
			notYetPlayed = append(notYetPlayed, fmt.Sprintf("<@%d> - %d Strokes %s", leader.PlayerID, leader.TotalStrokes, placementString))
			skipCounter++
		}
	}

	leaderString := "No one has played today!\n"
	if len(leaderStrings) > 0 {
		leaderString = "Leaders\n" + strings.Join(leaderStrings, "\n")
	}

	notYetPlayedString := ""
	if len(notYetPlayed) > 0 {
		rand.Shuffle(len(notYetPlayed), func(i, j int) { notYetPlayed[i], notYetPlayed[j] = notYetPlayed[j], notYetPlayed[i] })

		notYetPlayedString = "\n\nNot Played Yet\n" + strings.Join(notYetPlayed, "\n") + "\n"
	}

	return leaderString + notYetPlayedString
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

func (l *Leaderboard) getPreviousPlacements(ctx context.Context, guildID int64, tournamentID int32) map[int64]int {
	previousPlacements, err := l.service.GetPlacementsForPeriod(ctx, tournamentID, time.Now().Add(-24*time.Hour))
	if err != nil {
		slog.Error("Failed to get previous placements", "guild", guildID, "error", err)
		return nil
	}

	previous := make(map[int64]int)
	for i, previousPlacement := range previousPlacements {
		previous[previousPlacement.PlayerID] = i + 1
	}

	return previous
}

func (l *Leaderboard) getPreviousWins(ctx context.Context, guildID int64) map[int64]int64 {
	previousWins, err := l.service.GetTournamentPlacementsByPosition(ctx, guildID, 1)
	if err != nil {
		slog.Error("Failed to get previous wins", "guild", guildID, "error", err)
		return nil
	}

	previous := make(map[int64]int64)
	for _, previousWin := range previousWins {
		previous[previousWin.PlayerID] = previousWin.Wins
	}

	return previous
}

// TODO: this should not have a newline for empty results
func (l *Leaderboard) generateStats(ctx context.Context, tournamentID int32) string {
	holeInOneString := l.getHoleInOneLeader(ctx, tournamentID)
	bestRoundString := l.getBestRounds(ctx, tournamentID)
	worstRoundString := l.getWorstRounds(ctx, tournamentID)
	mostCommon := l.getFirstMostCommonHole(ctx, tournamentID)
	lastCommon := l.getLastMostCommonHole(ctx, tournamentID)
	hardestHole := l.getHardestHole(ctx, tournamentID)
	bestPerformers := l.getBestPerformers(ctx)
	worstPerformers := l.getWorstPerformers(ctx)

	statsHeader := "Stats powered by AWS Next Gen Stats"
	statsStr := strings.Join([]string{
		statsHeader,
		holeInOneString,
		bestRoundString,
		worstRoundString,
		mostCommon,
		lastCommon,
		hardestHole,
		bestPerformers,
		worstPerformers,
	}, "\n")

	return statsStr
}

func (l *Leaderboard) getHoleInOneLeader(ctx context.Context, tournamentID int32) string {
	holeInOneLeaders, err := l.service.GetHoleInOneLeaders(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get hole in one leaders", "tournament", tournamentID, "error", err)
		return ""
	}

	if len(holeInOneLeaders) == 0 {
		return ""
	}

	holeInOnes := int64(0)
	leaders := []string{}
	for i, leader := range holeInOneLeaders {
		if i == 0 || leader.Count == holeInOneLeaders[0].Count {
			leaders = append(leaders, fmt.Sprintf("<@%d>", leader.PlayerID))
			holeInOnes = leader.Count
		} else {
			break
		}
	}

	plural := ""
	if holeInOnes > 1 {
		plural = "s"
	}

	holeInOneMentions := strings.Join(leaders, ", ")

	return fmt.Sprintf("Most hole in ones: %s with %d hole in one%s", holeInOneMentions, holeInOnes, plural)
}

// TOOD: dedupe these two methods
func (l *Leaderboard) getBestRounds(ctx context.Context, tournamentID int32) string {
	rounds, err := l.service.GetBestRounds(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get best rounds", "tournament", tournamentID, "error", err)
		return ""
	}

	if len(rounds) == 0 {
		return ""
	}

	strokes := int64(0)
	mentions := []string{}
	for i, round := range rounds {
		if i == 0 || round.TotalStrokes == rounds[0].TotalStrokes {
			mentions = append(mentions, fmt.Sprintf("<@%d>", round.PlayerID))
			strokes = int64(round.TotalStrokes)
		} else {
			break
		}
	}

	plural := ""
	if len(mentions) > 1 {
		plural = "s"
	}

	bestMentions := strings.Join(mentions, ", ")

	return fmt.Sprintf("Best round%s of the tournament: %s, %d strokes 🙇", plural, bestMentions, strokes)
}

func (l *Leaderboard) getWorstRounds(ctx context.Context, tournamentID int32) string {
	rounds, err := l.service.GetWorstRounds(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get best rounds", "tournament", tournamentID, "error", err)
		return ""
	}

	if len(rounds) == 0 {
		return ""
	}

	strokes := int64(0)
	mentions := []string{}
	for i, round := range rounds {
		if i == 0 || round.TotalStrokes == rounds[0].TotalStrokes {
			mentions = append(mentions, fmt.Sprintf("<@%d>", round.PlayerID))
			strokes = int64(round.TotalStrokes)
		} else {
			break
		}
	}

	plural := ""
	if len(mentions) > 1 {
		plural = "s"
	}

	worstMentions := strings.Join(mentions, ", ")

	return fmt.Sprintf("Worst round%s of the tournament: %s, %d strokes 🤡%s", plural, worstMentions, strokes, plural)
}

func (l *Leaderboard) getFirstMostCommonHole(ctx context.Context, tournamentID int32) string {
	firstMost, err := l.service.GetMostCommonHoleForNumber(ctx, tournamentID, 0)
	if err != nil {
		slog.Error("Failed to get most common hole for number", "tournament", tournamentID, "error", err)
		return ""
	}

	return fmt.Sprintf("Most common opening hole: %s", colors[firstMost.Color])

}

func (l *Leaderboard) getLastMostCommonHole(ctx context.Context, tournamentID int32) string {
	lastMost, err := l.service.GetMostCommonHoleForNumber(ctx, tournamentID, 4)
	if err != nil {
		slog.Error("Failed to get most common hole for number", "tournament", tournamentID, "error", err)
		return ""
	}

	return fmt.Sprintf("Most common finishing hole: %s", colors[lastMost.Color])
}

func (l *Leaderboard) getHardestHole(ctx context.Context, tournamentID int32) string {
	hole, err := l.service.GetHardestHole(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get hardest hole", "tournament", tournamentID, "error", err)
	}
	return fmt.Sprintf("Hardest hole: %s with an average of %0.2f strokes", colors[hole.Color], hole.Strokes)
}

func (l *Leaderboard) getPerformers(ctx context.Context, reverse bool) ([]string, string, error) {
	performers, err := l.service.GetStandardDeviation(ctx, reverse)
	if err != nil {
		slog.Error("Failed to get performers", "error", err)
		return nil, "", err
	}

	stdDev := ""
	performerIDs := []string{}
	for i, performer := range performers {
		if i == 0 || performer.StdDev == performers[i-1].StdDev {
			performerIDs = append(performerIDs, fmt.Sprintf("<@%d>", performer.PlayerID))
			stdDev = performer.StdDev
		} else {
			break
		}
	}

	return performerIDs, stdDev, nil
}

func (l *Leaderboard) getWorstPerformers(ctx context.Context) string {
	worstPerformerString := ""

	worstPerformerIDs, stdDev, err := l.getPerformers(ctx, true)
	if err != nil {
		return worstPerformerString
	}

	worstPerformerMentions := strings.Join(worstPerformerIDs, ", ")

	return fmt.Sprintf("[All Time] Least consistent players: %v with a standard deviation of %s strokes", worstPerformerMentions, stdDev)
}

func (l *Leaderboard) getBestPerformers(ctx context.Context) string {
	bestPerformersString := ""

	bestPerformerIDs, stdDev, err := l.getPerformers(ctx, false)
	if err != nil {
		return bestPerformersString
	}

	bestPerformersMentions := strings.Join(bestPerformerIDs, ", ")

	return fmt.Sprintf("[All Time] Most consistent players: %v with a standard deviation of %s strokes", bestPerformersMentions, stdDev)
}
