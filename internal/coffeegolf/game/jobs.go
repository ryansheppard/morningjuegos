package game

import (
	"database/sql"
	"log/slog"
	"math"
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

	for _, tournament := range tournaments {
		start := tournament.StartTime

		players, err := g.query.GetUniquePlayersInTournament(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get unique players in tournament", "tournament", tournament, "error", err)
			continue
		}

		// Take the time since the start of the tournament, round the float to the nearest integer,
		// and then subtract one to remove the current day.
		// Missing rounds should only be added after the entire day has passed.
		numDaysPlayed := math.Floor(time.Since(start).Hours()/24) - 1

		for i := float64(0); i < numDaysPlayed; i++ {
			day := start.Add(time.Duration(i) * 24 * time.Hour)
			nullTime := sql.NullTime{
				Time:  day,
				Valid: true,
			}
			for _, player := range players {
				_, err := g.query.HasPlayed(g.ctx, database.HasPlayedParams{
					PlayerID:     player,
					TournamentID: tournament.ID,
					RoundDate:    nullTime,
				})

				if err == sql.ErrNoRows {
					slog.Info("Adding missing round", "player", player, "tournament", tournament, "day", day)
					entry := &database.Round{
						PlayerID:     player,
						TournamentID: tournament.ID,
						TotalStrokes: defaultStrokes,
						Percentage:   "",
						RoundDate:    nullTime,
					}

					_, err := g.query.CreateRound(g.ctx, database.CreateRoundParams{
						TournamentID: entry.TournamentID,
						PlayerID:     entry.PlayerID,
						TotalStrokes: defaultStrokes,
						OriginalDate: day.Format("Jan 02"),
						FirstRound:   true,
						InsertedBy:   "add_missing_rounds",
						RoundDate:    entry.RoundDate,
					},
					)
					if err != nil {
						slog.Error("Failed to create round", "round", entry, "error", err)
						continue
					}
				} else if err != nil {
					slog.Error("Failed to check if player has played", "player", player, "tournament", tournament, "day", day, "error", err)
					continue
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
