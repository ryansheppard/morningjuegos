package coffeegolf

import (
	"github.com/bwmarrin/discordgo"
)

var LeaderboardCommand = &discordgo.ApplicationCommand{
	Name:        "coffeegolf",
	Description: "Gets the leaderboard for Coffee Golf",
}

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: generateLeaderboard(i.GuildID),
		},
	})
}
