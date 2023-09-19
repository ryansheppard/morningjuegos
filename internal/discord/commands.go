package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/ryansheppard/morningjuegos/internal/games/coffeegolf"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "leaderboard",
			Description: "Get the leaderboard for Coffee Golf",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// TODO: this should not be here
			leaders := coffeegolf.GetLeaders(5)
			leaderStrings := []string{}
			for i, leader := range leaders {
				leaderStrings = append(leaderStrings, fmt.Sprintf("%d: %s - %d Total Strokes", i+1, leader.PlayerName, leader.TotalStrokes))
			}

			leaderString := strings.Join(leaderStrings, "\n")

			hole := coffeegolf.GetHardestHole()
			holeString := fmt.Sprintf("The hardest hole was %s and took an average of %d strokes\n", hole.Color, hole.Strokes)

			firstMost := coffeegolf.MostCommonFirstHole()
			lastMost := coffeegolf.MostCommonLastHole()
			mostCommonString := fmt.Sprintf("The most common first hole was %s and the last was %s", firstMost, lastMost)

			all := leaderString + "\n" + "Stats powered by AWS Next Gen Stats" + "\n" + holeString + "\n" + mostCommonString
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: all,
				},
			})
		},
	}
)
