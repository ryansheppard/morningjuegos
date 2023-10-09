package discord

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	cg "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

type Discord struct {
	ctx        context.Context
	Session    *discordgo.Session
	AppID      string
	messenger  *messenger.Messenger
	cache      *cache.Cache
	CoffeeGolf *cg.Game
}

func NewDiscord(ctx context.Context, token string, appID string, messenger *messenger.Messenger, cache *cache.Cache, cg *cg.Game) (*Discord, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("Error creating Discord session", "error", err)
		return nil, err
	}

	return &Discord{
		ctx:        ctx,
		Session:    dg,
		AppID:      appID,
		messenger:  messenger,
		cache:      cache,
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
