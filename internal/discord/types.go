package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var Parsers []Parser

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
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentMessageContent | discordgo.IntentGuildMessageReactions

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	return &Discord{
		Discord: dg,
		AppID:   appID,
	}
}

func (d *Discord) AddParser(parser Parser) {
	Parsers = append(Parsers, parser)
}

func (d *Discord) AddCommand(command *discordgo.ApplicationCommand, handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
	_, err := d.Discord.ApplicationCommandCreate(d.AppID, "", command)
	if err != nil {
		log.Panicf("Cannot create '%v' command: %v", command.Name, err)
	}

	d.Discord.AddHandler(handler)
}
