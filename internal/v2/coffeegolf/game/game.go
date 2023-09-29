package game

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/parser"
)

type Game struct {
	ctx         context.Context
	query       *database.Queries
	cache       *cache.Cache
	parser      *parser.Parser
	leaderboard *leaderboard.Leaderboard
}

func New(ctx context.Context, query *database.Queries, cache *cache.Cache) *Game {
	parser := parser.New(ctx, query)
	leaderboard := leaderboard.New(ctx, query, cache)
	return &Game{
		ctx:         ctx,
		query:       query,
		cache:       cache,
		parser:      parser,
		leaderboard: leaderboard,
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

func (g *Game) GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "coffeegolf",
			Description: "Gets the leaderboard for Coffee Golf",
		},
	}
}

func (g *Game) GetHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coffeegolf": g.LeaderboardCmd,
	}
}
