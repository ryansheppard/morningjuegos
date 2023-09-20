package coffeegolf

import (
	"fmt"
	"strings"
	"time"
)

// TODO: refactor to use a struct
func generateLeaderboard(guildID string) string {
	now := time.Now().Unix()
	leaders := getLeaders(guildID, 5, now)
	if len(leaders) == 0 {
		return "No one has played yet!"
	}

	leaderStrings := []string{}
	for i, leader := range leaders {
		leaderStrings = append(leaderStrings, fmt.Sprintf("%d: <@%s> - %d Total Strokes", i+1, leader.PlayerID, leader.TotalStrokes))
	}

	leaderString := strings.Join(leaderStrings, "\n")

	hole := getHardestHole(guildID, now)
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %0.2f strokes\n", hole.Color, hole.Strokes)

	firstMost := mostCommonFirstHole(guildID, now)
	lastMost := mostCommonLastHole(guildID, now)
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

	statsStr := "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString

	all := leaderString + "\n" + statsStr

	return all
}
