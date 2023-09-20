package coffeegolf

import (
	"fmt"
	"strings"
	"time"
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

	strokeLeaders := getStrokeLeaders(guildID, tournament.ID)
	if len(strokeLeaders) == 0 {
		return "No one has played yet!"
	}

	leaderStrings := []string{}
	for i, leader := range strokeLeaders {
		leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%s> - %d Total Strokes", i+1, leader.PlayerID, leader.TotalStrokes))
	}

	leaderString := "Leaders\n" + strings.Join(leaderStrings, "\n")

	hole := getHardestHole(guildID, tournament.ID)
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %0.2f strokes\n", hole.Color, hole.Strokes)

	firstMost := mostCommonFirstHole(guildID, tournament.ID)
	lastMost := mostCommonLastHole(guildID, tournament.ID)
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

	statsStr := "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString

	all := tournamentString + "\n\n" + leaderString + "\n" + statsStr

	return all
}
