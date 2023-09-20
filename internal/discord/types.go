package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/game"
)

var Parsers []game.Parser

type Discord struct {
	Discord *discordgo.Session
	AppID   string
}

func NewDiscord(token string, appID string) *Discord {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageCreate)

	// TODO: clean these up
	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	return &Discord{
		Discord: dg,
		AppID:   appID,
	}
}

func (d *Discord) RegisterGame(g *game.Game) {
	d.AddParser(g.Parser)

	for _, command := range g.Commands {
		d.AddCommand(command)
	}

	for _, handler := range g.Handlers {
		d.AddCommandHandler(handler)
	}
}

func (d *Discord) AddParser(parser game.Parser) {
	Parsers = append(Parsers, parser)
}

func (d *Discord) AddCommand(command *discordgo.ApplicationCommand) {
	_, err := d.Discord.ApplicationCommandCreate(d.AppID, "", command)
	if err != nil {
		fmt.Errorf("Cannot create '%v' command: %v", command.Name, err)
	}
}

func (d *Discord) AddCommandHandler(handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
	d.Discord.AddHandler(handler)
}
