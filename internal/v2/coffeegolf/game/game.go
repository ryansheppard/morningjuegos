package game

import (
	"context"
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/v2/cache"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/parser"
)

type Game struct {
	ctx         context.Context
	query       *database.Queries
	cache       *cache.Cache
	Parser      *parser.Parser
	leaderboard *leaderboard.Leaderboard
}

func New(ctx context.Context, query *database.Queries, cache *cache.Cache, db *sql.DB) *Game {
	parser := parser.New(ctx, query, db)
	leaderboard := leaderboard.New(ctx, query, cache)
	return &Game{
		ctx:         ctx,
		query:       query,
		cache:       cache,
		Parser:      parser,
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

func (g *Game) HelpText() string {
	return `
Possible responses are: 
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
			Name:        "coffeegolfhelp",
			Description: "Gets the help for Coffee Golf",
		},
	}
}

func (g *Game) GetHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coffeegolf":     g.LeaderboardCmd,
		"coffeegolfhelp": g.HelpCmd,
	}
}
