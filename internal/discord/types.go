package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
)

type Discord struct {
	Discord    *discordgo.Session
	AppID      string
	CoffeeGolf *coffeegolf.CoffeeGolf
}

func NewDiscord(token string, appID string, cg *coffeegolf.CoffeeGolf) *Discord {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	return &Discord{
		Discord:    dg,
		AppID:      appID,
		CoffeeGolf: cg,
	}
}

func (d *Discord) Configure() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "coffeegolf",
			Description: "Gets the leaderboard for Coffee Golf",
		},
	}

	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coffeegolf": d.CoffeeGolf.LeaderboardCmd,
	}
	for _, command := range commands {
		d.AddCommand(command)
	}

	for _, handler := range handlers {
		d.AddCommandHandler(handler)
	}

	d.Discord.AddHandler(d.messageCreate)

	// TODO: clean these up
	d.Discord.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions

	err := d.Discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return err
	}

	return nil
}

func (d *Discord) AddCommand(command *discordgo.ApplicationCommand) {
	_, err := d.Discord.ApplicationCommandCreate(d.AppID, "", command)
	if err != nil {
		fmt.Printf("Cannot create '%v' command: %v", command.Name, err)
	}
}

func (d *Discord) AddCommandHandler(handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
	d.Discord.AddHandler(handler)
}
