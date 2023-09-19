package coffeegolf

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var LeaderboardCommand = &discordgo.ApplicationCommand{
	Name:        "coffeegolf",
	Description: "Gets the leaderboard for Coffee Golf",
}

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	leaders := GetLeaders(5)
	leaderStrings := []string{}
	for i, leader := range leaders {
		leaderStrings = append(leaderStrings, fmt.Sprintf("%d: %s - %d Total Strokes", i+1, leader.PlayerName, leader.TotalStrokes))
	}

	leaderString := strings.Join(leaderStrings, "\n")

	hole := GetHardestHole()
	holeString := fmt.Sprintf("The hardest hole was %s and took an average of %d strokes\n", hole.Color, hole.Strokes)

	firstMost := MostCommonFirstHole()
	lastMost := MostCommonLastHole()
	mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

	all := leaderString + "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: all,
		},
	})
}
