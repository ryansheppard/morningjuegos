package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	cg "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

type Discord struct {
	Session    *discordgo.Session
	AppID      string
	Messenger  *messenger.Messenger
	CoffeeGolf *cg.Game
}

func NewDiscord(token string, appID string, messenger *messenger.Messenger, cg *cg.Game) (*Discord, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		return nil, err
	}

	return &Discord{
		Session:    dg,
		AppID:      appID,
		Messenger:  messenger,
		CoffeeGolf: cg,
	}, nil
}

func (d *Discord) Configure() error {
	for _, command := range d.CoffeeGolf.GetCommands() {
		d.AddCommand(command)
	}

	d.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.CoffeeGolf.GetHandlers()[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	d.Session.AddHandler(d.messageCreate)

	d.Session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildMessageReactions

	err := d.Session.Open()
	if err != nil {
		slog.Error("Error opening Discord session: ", "error", err)
		return err
	}

	return nil
}

func (d *Discord) AddCommand(command *discordgo.ApplicationCommand) {
	slog.Info("Adding command", "command", command.Name)
	_, err := d.Session.ApplicationCommandCreate(d.AppID, "", command)
	if err != nil {
		slog.Error("Cannot create command", "command", command.Name, "error", err)
	}
}
