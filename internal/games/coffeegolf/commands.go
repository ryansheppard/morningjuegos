package coffeegolf

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "coffeegolf",
		Description: "Gets the leaderboard for Coffee Golf",
	},
}

var Handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: generateLeaderboard(i.GuildID),
			},
		})
	},
}
