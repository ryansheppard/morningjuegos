package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	cg "github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/game"
)

type Discord struct {
	Discord    *discordgo.Session
	AppID      string
	CoffeeGolf *cg.Game
}

func NewDiscord(token string, appID string, cg *cg.Game) (*Discord, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		return nil, err
	}

	return &Discord{
		Discord:    dg,
		AppID:      appID,
		CoffeeGolf: cg,
	}, nil
}

func (d *Discord) Configure() error {
	for _, command := range d.CoffeeGolf.GetCommands() {
		d.AddCommand(command)
	}

	d.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.CoffeeGolf.GetHandlers()[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	d.Discord.AddHandler(d.messageCreate)

	// TODO: clean these up
	d.Discord.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions

	err := d.Discord.Open()
	if err != nil {
		slog.Error("Error opening Discord session: ", "error", err)
		return err
	}

	return nil
}

func (d *Discord) AddCommand(command *discordgo.ApplicationCommand) {
	slog.Info("Adding command", "command", command.Name)
	_, err := d.Discord.ApplicationCommandCreate(d.AppID, "", command)
	if err != nil {
		slog.Error("Cannot create command", "command", command.Name, "error", err)
	}
}
