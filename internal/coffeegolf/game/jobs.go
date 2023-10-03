package game

import (
	"database/sql"
	"log/slog"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/messages"
)

const defaultStrokes = 20

func (g *Game) ConfigureSubscribers() {
	g.messenger.SubscribeAsync(messages.RoundCreatedKey, g.ProcessAddMissingRounds)
	g.messenger.SubscribeAsync(messages.TournamentCreatedKey, g.ProcessAddTournamentWinners)
}

func (g *Game) ProcessAddMissingRounds(msg *nats.Msg) {
	slog.Info("Processing add missing rounds message")
	roundCreated, err := messages.NewRoundCreatedFromJson(msg.Data)
	if err != nil {
		slog.Error("Failed to parse round created message", "error", err)
	}

	g.AddMissingRoundsForGuild(roundCreated.GuildID)
}

func (g *Game) AddMissingRounds() {
	slog.Info("Adding missing rounds for all guilds")

	guildIDs, err := g.query.GetAllGuilds(g.ctx)
	if err != nil {
		slog.Error("Failed to get all guilds", "error", err)
	}

	for _, guildID := range guildIDs {
		g.AddMissingRoundsForGuild(guildID)
	}

	slog.Info("Finished adding missing rounds for all guilds")
}

func (g *Game) AddMissingRoundsForGuild(guildID int64) {
	slog.Info("Adding missing rounds for guild", "guild", guildID)

	var tournaments []database.Tournament
	tournament, err := g.query.GetActiveTournament(g.ctx, database.GetActiveTournamentParams{
		GuildID:   guildID,
		StartTime: time.Now(),
	})
	if err == sql.ErrNoRows {
		slog.Error("No active tournament", "guild", guildID, "error", err)
		return
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
	}

	// clear cache
	cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
	g.cache.DeleteKey(cacheKey)

	tournaments = append(tournaments, tournament)

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

	slog.Info("Finished adding missing rounds for guild", "guild", guildID)
}

func (g *Game) ProcessAddTournamentWinners(msg *nats.Msg) {
	slog.Info("Processing add tournament winners message")
	roundCreated, err := messages.NewTournamentCreatedFromJson(msg.Data)
	if err != nil {
		slog.Error("Failed to parse tournament created message", "error", err)
	}

	g.AddTournamentWinnersForGuild(roundCreated.GuildID)
}

func (g *Game) AddTournamentWinners() {
	slog.Info("Adding tournament winners for all guilds")
	guilds, err := g.query.GetAllGuilds(g.ctx)
	if err != nil {
		slog.Error("Failed to get all guilds", err)
	}

	for _, guild := range guilds {
		g.AddMissingRoundsForGuild(guild)
	}

	slog.Info("Finished adding tournament winners for all guilds")
}

func (g *Game) AddTournamentWinnersForGuild(guildID int64) {
	slog.Info("Adding tournament winners for guild", "guild", guildID)

	var inactiveTournaments []database.Tournament
	tournaments, err := g.query.GetInactiveTournaments(g.ctx, database.GetInactiveTournamentsParams{
		GuildID: guildID,
		EndTime: time.Now(),
	})
	if err != nil {
		slog.Error("Failed to get inactive tournaments", "guild", guildID, "error", err)
		return
	}
	if len(tournaments) > 0 {
		cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
		g.cache.DeleteKey(cacheKey)
		inactiveTournaments = append(inactiveTournaments, tournaments...)
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

	slog.Info("Finished adding tournament winners for guild", "guild", guildID)
}
