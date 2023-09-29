package game

import (
	"log/slog"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
)

const defaultStrokes = 20

func (g *Game) AddMissingRounds() {
	guilds, err := g.query.GetAllGuilds(g.ctx)
	if err != nil {
		slog.Error("Failed to get all guilds", err)
	}

	var tournaments []database.Tournament
	for _, guild := range guilds {
		tournament, err := g.query.GetActiveTournament(g.ctx, guild.GuildID)
		if err != nil {
			slog.Error("Failed to get active tournament", "guild", guild, err)
			continue
		}
		if tournament != (database.Tournament{}) {
			tournaments = append(tournaments, tournament)
		}
	}

	now := time.Now().Unix()

	for _, tournament := range tournaments {
		start := tournament.StartTime.Unix()
		// Skip tournaments that started less than 24 hours ago
		if now-start < 86400 {
			continue
		}

		numDaysPlayed := (now - start) / 86400
		players, err := g.query.GetUniquePlayersInTournament(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get unique players in tournament", "tournament", tournament, err)
			continue
		}

		for i := int64(0); i < numDaysPlayed; i++ {
			day := start + (i * 86400)
			for _, player := range players {
				exists, err := g.query.HasPlayed(g.ctx, database.HasPlayedParams{
					PlayerID:     player,
					TournamentID: tournament.ID,
					DateTrunc:    day,
				})
				if err != nil {
					slog.Error("Failed to check if player has played", "player", player, "tournament", tournament, "day", day, err)
				}
				if exists == (database.Round{}) {
					entry := &database.Round{
						PlayerID:     player,
						TournamentID: tournament.ID,
						TotalStrokes: defaultStrokes,
						InsertedAt:   time.Unix(day, 0),
						Percentage:   "",
					}

					g.query.CreateRound(g.ctx, database.CreateRoundParams{
						TournamentID: entry.TournamentID,
						PlayerID:     entry.PlayerID,
						TotalStrokes: defaultStrokes,
						OriginalDate: "",
						FirstRound:   true,
					},
					)
				}
			}
		}
	}
}

// func (g *Game) AddTournamentWinners() {
// 	guilds, err := g.query.GetAllGuilds(g.ctx)
// 	if err != nil {
// 		slog.Error("Failed to get all guilds", err)
// 	}

// 	var inactiveTournaments []database.Tournament
// 	for _, guild := range guilds {
// 		tournaments := g.query.GetInactiveTournaments(guild)
// 		if len(tournaments) > 0 {
// 			inactiveTournaments = append(inactiveTournaments, tournaments...)
// 		}
// 	}

// 	for _, tournament := range inactiveTournaments {
// 		uniquePlayers := cg.Query.getUniquePlayersInTournament(tournament.ID)
// 		placements := cg.Query.getTournamentPlacements(tournament.ID)

// 		if len(uniquePlayers) != len(placements) {
// 			cg.Query.createTournamentPlacements(tournament.ID, tournament.GuildID)
// 		}
// 	}
// }
