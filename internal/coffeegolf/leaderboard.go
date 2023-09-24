package coffeegolf

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/utils"
)

// TODO: refactor to use a struct
func generateLeaderboard(guildID string) string {
	tournament := getActiveTournament(guildID, false)
	if tournament == nil {
		return "No active tournament"
	}

	startDate := time.Unix(tournament.Start, 0).Format("Jan 2, 2006")
	endDate := time.Unix(tournament.End, 0).Format("Jan 2, 2006")

	tournamentString := fmt.Sprintf("Current Tournament: %s - %s", startDate, endDate)

	strokeLeaders := getStrokeLeaders(guildID, tournament.ID, 100)
	if len(strokeLeaders) == 0 {
		return "No one has played yet!"
	}

	leaderStrings := []string{}
	notYetPlayed := []string{}
	startOfDay := utils.GetStartofDay(time.Now().Unix())
	skipCounter := 0
	for i, leader := range strokeLeaders {
		hasPlayedToday := checkIfPlayerHasRound(leader.PlayerID, tournament.ID, startOfDay)
		if hasPlayedToday {
			leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%s> - %d Total Strokes", i+1-skipCounter, leader.PlayerID, leader.TotalStrokes))
		} else {
			notYetPlayed = append(notYetPlayed, fmt.Sprintf("<@%s> - %d Total Strokes", leader.PlayerID, leader.TotalStrokes))
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

	hole := getHardestHole(guildID, tournament.ID)
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %0.2f strokes", hole.Color, hole.Strokes)

	firstMost := mostCommonFirstHole(guildID, tournament.ID)
	lastMost := mostCommonLastHole(guildID, tournament.ID)
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

	worstRound := getWorstRound(guildID, tournament.ID)
	worstRoundString := fmt.Sprintf("The worst round was %d strokes by <@%s>", worstRound.TotalStrokes, worstRound.PlayerID)

	statsStr := "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString + "\n" + worstRoundString

	all := tournamentString + "\n\n" + leaderString + "\n" + notYetPlayedString + "\n" + statsStr

	return all
}
