package game

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/parser"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Game struct {
	ctx         context.Context
	query       *database.Queries
	cache       *cache.Cache
	Parser      *parser.Parser
	leaderboard *leaderboard.Leaderboard
	messenger   *messenger.Messenger
	tracer      trace.Tracer
}

// Todo replace with withoptions
func New(ctx context.Context, query *database.Queries, cache *cache.Cache, db *sql.DB, messenger *messenger.Messenger, tracer trace.Tracer) *Game {
	parser := parser.New(ctx, query, db, cache, messenger, tracer)
	leaderboard := leaderboard.New(ctx, query, cache, tracer)
	return &Game{
		ctx:         ctx,
		query:       query,
		cache:       cache,
		Parser:      parser,
		leaderboard: leaderboard,
		messenger:   messenger,
		tracer:      tracer,
	}
}

func (g *Game) LeaderboardCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, span := g.tracer.Start(g.ctx, "leaderboard-command")
	fmt.Println(span.SpanContext().SpanID())
	defer span.End()

	params := leaderboard.GenerateLeaderboardParams{
		GuildID: i.GuildID,
	}

	span.SetAttributes(attribute.String("guild.id", i.GuildID))
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if option, ok := optionMap["date-option"]; ok {
		params.SetDate(option.StringValue())
		span.SetAttributes(attribute.String("date", option.StringValue()))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: g.leaderboard.GenerateLeaderboard(ctx, params),
		},
	})
}

func (g *Game) StatsCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, span := g.tracer.Start(g.ctx, "stats-command")
	defer span.End()

	span.SetAttributes(attribute.String("guild.id", i.GuildID))

	span.AddEvent("generating-stats")
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: g.leaderboard.GenerateStats(ctx, i.GuildID),
		},
	})
	span.AddEvent("finished-generating-stats")
}

func (g *Game) HelpText() string {
	_, span := g.tracer.Start(g.ctx, "help-command")
	defer span.End()

	return `
Use /coffeegolf to the leaderboard for the current tournament
Use /coffeegolf <date> to get the leaderboard for a specific date. Date format is YYYY-MM-DD
Use /coffeestats to get the stats for the current tournament

When posting a Coffee Golf message, the possible reactions are: 
- 👍: First round
- 👌: Bonus round
- 🤯: Parsed the message but failed to insert in to the database
- 🖕: Detected it was a cofee golf message, but failed to parse the message entirely
- 🤬: Something has gone wrong big time
- (no reaction): Not detected as a coffee golf message
	`
}

func (g *Game) HelpCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: g.HelpText(),
		},
	})
}

func (g *Game) GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "coffeegolf",
			Description: "Gets the leaderboard for Coffee Golf",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date-option",
					Description: "The date to get the stats for",
					Required:    false,
				},
			},
		},
		{
			Name:        "coffeestats",
			Description: "Gets the stats for Coffee Golf",
		},
		{
			Name:        "coffeehelp",
			Description: "Gets the help for Coffee Golf",
		},
	}
}

func (g *Game) GetHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coffeegolf":  g.LeaderboardCmd,
		"coffeestats": g.StatsCmd,
		"coffeehelp":  g.HelpCmd,
	}
}
