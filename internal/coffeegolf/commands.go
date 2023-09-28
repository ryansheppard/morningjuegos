// Package coffeegolf handles parsing and scoring the Coffee Golf game
package coffeegolf

import (
	"github.com/bwmarrin/discordgo"
)

// var commands = []*discordgo.ApplicationCommand{
// 	{
// 		Name:        "coffeegolf",
// 		Description: "Gets the leaderboard for Coffee Golf",
// 	},
// }

// var handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 	"coffeegolf": cg.leaderboardCmd,
// 	//  func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 	// 		Data: &discordgo.InteractionResponseData{
// 	// 			Content: generateLeaderboard(i.GuildID),
// 	// 		},
// 	// 	})
// 	// },
// }

func (cg *CoffeeGolf) LeaderboardCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cg.generateLeaderboard(i.GuildID),
		},
	})

}
