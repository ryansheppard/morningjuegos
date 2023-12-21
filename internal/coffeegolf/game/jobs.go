package game

import (
	"database/sql"
	"log/slog"
	"math"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

const defaultStrokes = 20

func (g *Game) ConfigureSubscribers() {
	g.messenger.SubscribeAsync(messenger.RoundCreatedKey, g.ProcessAddMissingRounds)
	g.messenger.SubscribeAsync(messenger.TournamentCreatedKey, g.ProcessAddTournamentWinners)
}

func (g *Game) ProcessAddMissingRounds(msg *nats.Msg) {
	slog.Info("Processing add missing rounds message")
	roundCreated, err := messenger.NewRoundCreatedFromJson(msg.Data)
	if err != nil {
		slog.Error("Failed to parse round created message", "error", err)
	}

	tournament, err := g.service.GetActiveTournament(g.ctx, roundCreated.GuildID)
	if err == sql.ErrNoRows {
		slog.Error("No active tournament", "guild", roundCreated.GuildID, "error", err)
		return
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", roundCreated.GuildID, "error", err)
	}

	g.AddMissingRoundsForGuild(roundCreated.GuildID, tournament.ID)
}

func (g *Game) AddMissingRoundsForGuild(guildID int64, tournamentID int32) {
	slog.Info("Adding missing rounds for guild", "guild", guildID, "tournament", tournamentID)

	tournament, err := g.service.GetTournament(g.ctx, tournamentID)
	if err == sql.ErrNoRows {
		slog.Error("No active tournament", "tournamentID", tournamentID, "error", err)
		return
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
	}

	// clear cache
	cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
	g.cache.DeleteKey(g.ctx, cacheKey)

	start := tournament.StartTime

	players, err := g.service.GetUniquePlayersInTournament(g.ctx, tournament.ID)
	if err != nil {
		slog.Error("Failed to get unique players in tournament", "tournamentID", tournament.ID, "error", err)
		return
	}

	// Take the time since the start of the tournament, round the float to the nearest integer,
	// and then subtract one to remove the current day.
	// Missing rounds should only be added after the entire day has passed.
	numDaysPlayed := math.Floor(time.Since(start).Hours()/24) - 1

	for i := float64(0); i <= numDaysPlayed; i++ {
		day := start.Add(time.Duration(i) * 24 * time.Hour)
		if day.Before(start) || day.After(tournament.EndTime) {
			break
		}
		for _, player := range players {
			hasPlayed, err := g.service.HasPlayed(g.ctx, player, tournament.ID, day)
			if err != nil {
				slog.Error("Failed to check if player has played", "player", player, "tournament", tournament, "day", day, "error", err)
				continue
			}
			if !hasPlayed {
				slog.Info("Adding missing round", "player", player, "tournament", tournament, "day", day)
				entry := &database.Round{
					PlayerID:     player,
					TournamentID: tournament.ID,
					TotalStrokes: defaultStrokes,
					Percentage:   "",
					RoundDate:    sql.NullTime{Time: day, Valid: true},
					FirstRound:   true,
					InsertedBy:   "add_missing_rounds",
					OriginalDate: day.Format("Jan 02"),
				}

				_, err := g.service.InsertRound(g.ctx, entry, []*database.Hole{})
				if err != nil {
					slog.Error("Failed to create round", "round", entry, "error", err)
					continue
				}
			}
		}
	}

	slog.Info("Finished adding missing rounds for guild", "guild", guildID)
}

func (g *Game) ProcessAddTournamentWinners(msg *nats.Msg) {
	slog.Info("Processing add tournament winners message")
	roundCreated, err := messenger.NewTournamentCreatedFromJson(msg.Data)
	if err != nil {
		slog.Error("Failed to parse tournament created message", "error", err)
	}

	g.AddTournamentWinnersForGuild(roundCreated.GuildID)
}

func (g *Game) AddTournamentWinnersForGuild(guildID int64) {
	slog.Info("Adding tournament winners for guild", "guild", guildID)

	tournaments, err := g.service.GetInactiveTournaments(g.ctx, guildID)
	if err != nil {
		slog.Error("Failed to get inactive tournaments", "guild", guildID, "error", err)
		return
	}
	if len(tournaments) > 0 {
		cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
		g.cache.DeleteKey(g.ctx, cacheKey)
	}

	for _, tournament := range tournaments {
		g.AddMissingRoundsForGuild(guildID, tournament.ID)

		uniquePlayers, err := g.service.GetUniquePlayersInTournament(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get unique players in tournament", "tournament", tournament, "error", err)
			continue
		}

		// We don't care about the error in this case as an empty list gets returned
		placements, err := g.service.GetTournamentPlacements(g.ctx, tournament.ID)
		if err != nil {
			slog.Error("Failed to get tournament placements", "tournament", tournament, "error", err)
			continue
		}

		if len(uniquePlayers) != len(placements) {
			leaders, err := g.service.GetLeaders(g.ctx, tournament.ID)
			if err != nil {
				slog.Error("Failed to get final leaders", "tournament", tournament, "error", err)
				continue
			}

			err = g.service.CleanTournamentPlacements(g.ctx, tournament.ID)
			if err != nil {
				slog.Error("Failed to clean tournament placements", "tournament", tournament, "error", err)
				continue
			}

			for i, leader := range leaders {
				g.service.CreateTournamentPlacement(g.ctx, leader.PlayerID, tournament.ID, i+1, leader.TotalStrokes, "add_tournament_winners")
			}
		}
	}

	slog.Info("Finished adding tournament winners for guild", "guild", guildID)
}
