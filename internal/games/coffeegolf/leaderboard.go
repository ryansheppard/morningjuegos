package coffeegolf

import (
	"fmt"
	"strings"
)

// TODO: refactor to use a struct
func generateLeaderboard(guildID string) string {
	leaders := GetLeaders(guildID, 5)
	leaderStrings := []string{}
	for i, leader := range leaders {
		leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%s> - %d Total Strokes", i+1, leader.PlayerID, leader.TotalStrokes))
	}

	leaderString := strings.Join(leaderStrings, "\n")

	hole := GetHardestHole(guildID)
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %d strokes\n", hole.Color, hole.Strokes)

	firstMost := MostCommonFirstHole(guildID)
	lastMost := MostCommonLastHole(guildID)
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

	all := leaderString + "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString

	return all
}
