package game

import (
	"context"
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/parser"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/service"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

type Game struct {
	ctx         context.Context
	service     *service.Service
	cache       *cache.Cache
	Parser      *parser.Parser
	leaderboard *leaderboard.Leaderboard
	messenger   *messenger.Messenger
}

// Todo replace with withoptions
func New(ctx context.Context, service *service.Service, cache *cache.Cache, db *sql.DB, messenger *messenger.Messenger) *Game {
	parser := parser.New(ctx, service, cache, messenger)
	leaderboard := leaderboard.New(ctx, service, cache)
	return &Game{
		ctx:         ctx,
		service:     service,
		cache:       cache,
		Parser:      parser,
		leaderboard: leaderboard,
		messenger:   messenger,
	}
}

func (g *Game) LeaderboardCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: g.leaderboard.GenerateLeaderboard(i.GuildID),
		},
	})
}

func (g *Game) StatsCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: g.leaderboard.GenerateStats(i.GuildID),
		},
	})
}

func (g *Game) HelpText() string {
	return `
Use /coffeegolf to the leaderboard for the current tournament
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
