package game

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
)

const defaultStrokes = 20

func (g *Game) AddMissingRounds() {
	slog.Info("Adding missing rounds")
	guildIDs, err := g.query.GetAllGuilds(g.ctx)
	if err != nil {
		slog.Error("Failed to get all guilds", "error", err)
	}

	var tournaments []database.Tournament
	for _, guildID := range guildIDs {
		tournament, err := g.query.GetActiveTournament(g.ctx, database.GetActiveTournamentParams{
			GuildID:   guildID,
			StartTime: time.Now(),
		})
		if err == sql.ErrNoRows {
			slog.Error("No active tournament", "guild", guildID, "error", err)
			continue
		} else if err != nil {
			slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		}

		// clear cache
		cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
		g.cache.DeleteKey(cacheKey)

		tournaments = append(tournaments, tournament)
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
			slog.Error("Failed to get unique players in tournament", "tournament", tournament, "error", err)
			continue
		}

		for i := int64(0); i < numDaysPlayed; i++ {
			day := start + (i * 86400)
			for _, player := range players {
				_, err := g.query.HasPlayed(g.ctx, database.HasPlayedParams{
					PlayerID:     player,
					TournamentID: tournament.ID,
					DateTrunc:    day,
				})
				if err == sql.ErrNoRows {
					slog.Info("Adding missing round", "player", player, "tournament", tournament, "day", day)
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
						InsertedBy:   "add_missing_rounds",
					},
					)
				} else if err != nil {
					slog.Error("Failed to check if player has played", "player", player, "tournament", tournament, "day", day, "error", err)
				}
			}
		}
	}

	slog.Info("Finished adding missing rounds")
}

func (g *Game) AddTournamentWinners() {
	slog.Info("Adding tournament winners")
	guilds, err := g.query.GetAllGuilds(g.ctx)
	if err != nil {
		slog.Error("Failed to get all guilds", err)
	}

	var inactiveTournaments []database.Tournament
	for _, guild := range guilds {
		tournaments, err := g.query.GetInactiveTournaments(g.ctx, database.GetInactiveTournamentsParams{
			GuildID: guild,
			EndTime: time.Now(),
		})
		if err != nil {
			slog.Error("Failed to get inactive tournaments", "guild", guild, "error", err)
			continue
		}
		if len(tournaments) > 0 {
			cacheKey := leaderboard.GetLeaderboardCacheKey(guild)
			g.cache.DeleteKey(cacheKey)
			inactiveTournaments = append(inactiveTournaments, tournaments...)
		}
	}

	for _, tournament := range inactiveTournaments {
		uniquePlayers, err := g.query.GetUniquePlayersInTournament(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get unique players in tournament", "tournament", tournament, "error", err)
			continue
		}

		placements, err := g.query.GetTournamentPlacements(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get tournament placements", "tournament", tournament, "error", err)
			continue
		}

		if len(uniquePlayers) != len(placements) {
			placements, err := g.query.GetFinalLeaders(g.ctx, tournament.ID)
			if err != nil {
				slog.Error("Failed to get final leaders", "tournament", tournament, "error", err)
				continue
			}
			// transaction
			g.query.CleanTournamentPlacements(g.ctx, tournament.ID)
			for i, placement := range placements {
				g.query.CreateTournamentPlacement(g.ctx, database.CreateTournamentPlacementParams{
					TournamentID:        tournament.ID,
					PlayerID:            placement.PlayerID,
					TournamentPlacement: int32(i + 1),
					Strokes:             int32(placement.TotalStrokes),
					InsertedBy:          "add_tournament_winners",
				})
			}
		}
	}
	slog.Info("Finished adding tournament winners")
}
