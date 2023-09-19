package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Parser interface {
	ParseGame(message *discordgo.MessageCreate) ParserResponse
}

type ParserResponse struct {
	IsGame   bool
	Inserted bool
}

type Game struct {
	Parser   Parser
	Commands []*discordgo.ApplicationCommand
	Handlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
}

func NewGame(
	parser Parser,
	commands []*discordgo.ApplicationCommand,
	handlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate),
) *Game {
	return &Game{
		Parser:   parser,
		Commands: commands,
		Handlers: handlers,
	}
}

func (g *Game) Register(d *Discord) {
	d.AddParser(g.Parser)

	for _, command := range g.Commands {
		d.AddCommand(command)
	}

	for _, handler := range g.Handlers {
		d.AddCommandHandler(handler)
	}
}
