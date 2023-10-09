package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

func (s *Service) GetActiveTournament(ctx context.Context, guildID int64) (*database.Tournament, error) {
	tournament, err := s.queries.GetActiveTournament(ctx, guildID)
	if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return nil, err

	}

	return &database.Tournament{
		ID:        tournament.ID,
		GuildID:   tournament.GuildID,
		StartTime: tournament.StartTime,
		EndTime:   tournament.EndTime,
	}, nil
}

func (s *Service) GetInactiveTournaments(ctx context.Context, guildID int64) ([]*database.Tournament, error) {
	tournaments, err := s.queries.GetInactiveTournaments(ctx, database.GetInactiveTournamentsParams{
		GuildID: guildID,
		EndTime: time.Now(),
	})
	if err != nil {
		slog.Error("Failed to get inactive tournaments", "guild", guildID, "error", err)
		return nil, err
	}

	var inactiveTournaments []*database.Tournament
	for _, tournament := range tournaments {
		inactiveTournaments = append(inactiveTournaments, &database.Tournament{
			ID:        tournament.ID,
			GuildID:   tournament.GuildID,
			StartTime: tournament.StartTime,
			EndTime:   tournament.EndTime,
		})
	}

	return inactiveTournaments, nil
}

func (s *Service) CreateTournament(ctx context.Context, guildID int64, start time.Time, end time.Time, insertedBy string) (*database.Tournament, error) {
	tournament, err := s.queries.CreateTournament(ctx, database.CreateTournamentParams{
		GuildID:    guildID,
		StartTime:  start,
		EndTime:    end,
		InsertedBy: insertedBy,
	})
	if err != nil {
		slog.Error("Failed to create tournament", "guild", guildID, "error", err)
		return nil, err
	}

	return &database.Tournament{
		ID:         tournament.ID,
		GuildID:    tournament.GuildID,
		StartTime:  tournament.StartTime,
		EndTime:    tournament.EndTime,
		InsertedBy: tournament.InsertedBy,
	}, nil
}

func (s *Service) GetOrCreateTournament(ctx context.Context, guildID int64, insertedBy string) (*database.Tournament, bool, error) {
	tournament, err := s.GetActiveTournament(ctx, guildID)
	if err == sql.ErrNoRows {
		slog.Info("No active tournament found, creating one", "guild", guildID)
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		endDate := end.AddDate(0, 0, defaultTouramentLength-1)

		tournament, err := s.CreateTournament(ctx, guildID, start, endDate, insertedBy)
		if err != nil {
			slog.Error("Failed to create tournament", "guild", guildID, "error", err)
			return nil, false, err
		}

		slog.Info("Created tournament", "tournament", tournament.ID, "guildID", guildID)
		return &database.Tournament{
			ID:         tournament.ID,
			GuildID:    tournament.GuildID,
			StartTime:  tournament.StartTime,
			EndTime:    tournament.EndTime,
			InsertedBy: tournament.InsertedBy,
		}, true, nil
	} else if err != nil {
		slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
		return nil, false, err
	}

	return &database.Tournament{
		ID:         tournament.ID,
		GuildID:    tournament.GuildID,
		StartTime:  tournament.StartTime,
		EndTime:    tournament.EndTime,
		InsertedBy: tournament.InsertedBy,
	}, false, nil
}

func (s *Service) GetUniquePlayersInTournament(ctx context.Context, tournamentID int32) ([]int64, error) {
	players, err := s.queries.GetUniquePlayersInTournament(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get unique players in tournament", "tournamentID", tournamentID, "error", err)
		return nil, err
	}

	return players, nil
}

func (s *Service) GetTournamentPlacements(ctx context.Context, tournamentID int32) ([]database.TournamentPlacement, error) {
	placements, err := s.queries.GetTournamentPlacements(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get tournament placements", "tournamentID", tournamentID, "error", err)
		return nil, err
	}

	return placements, nil
}

type Leader struct {
	PlayerID     int64
	TotalStrokes int64
	Wins         int64
}

func (s *Service) GetFinalLeaders(ctx context.Context, tournamentID int32) ([]*Leader, error) {
	placements, err := s.queries.GetFinalLeaders(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get final leaders", "tournamentID", tournamentID, "error", err)
		return nil, err
	}

	leaders := []*Leader{}
	for _, placement := range placements {
		leaders = append(leaders, &Leader{
			PlayerID:     placement.PlayerID,
			TotalStrokes: placement.TotalStrokes,
		})
	}

	return leaders, nil
}

func (s *Service) CleanTournamentPlacements(ctx context.Context, tournamentID int32) error {
	err := s.queries.CleanTournamentPlacements(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to clean tournament placements", "tournamentID", tournamentID, "error", err)
		return err
	}
	return nil
}

func (s *Service) CreateTournamentPlacement(ctx context.Context, playerID int64, tournamentID int32, placement int, strokes int64, insertedBy string) error {
	_, err := s.queries.CreateTournamentPlacement(ctx, database.CreateTournamentPlacementParams{
		TournamentID:        tournamentID,
		PlayerID:            playerID,
		TournamentPlacement: int32(placement),
		Strokes:             int32(strokes),
		InsertedBy:          insertedBy,
	})
	if err != nil {
		slog.Error("Failed to create tournament placement", "tournamentID", tournamentID, "playerID", playerID, "placement", placement, "strokes", strokes, "error", err)
		return err
	}

	return nil
}

func (s *Service) GetLeaders(ctx context.Context, tournamentID int32) ([]*Leader, error) {
	leaders, err := s.queries.GetLeaders(ctx, tournamentID)

	if err != nil {
		slog.Error("Failed to get leaders", "tournamentid", tournamentID, "error", err)
		return nil, err
	}

	var leaderList []*Leader
	for _, leader := range leaders {
		leaderList = append(leaderList, &Leader{
			PlayerID:     leader.PlayerID,
			TotalStrokes: leader.TotalStrokes,
		})
	}

	return leaderList, nil
}

func (s *Service) GetTournamentPlacementsByPosition(ctx context.Context, guildID int64, tournamentPlacement int32) ([]*Leader, error) {
	previousWins, err := s.queries.GetTournamentPlacementsByPosition(ctx, database.GetTournamentPlacementsByPositionParams{
		GuildID:             guildID,
		TournamentPlacement: tournamentPlacement,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Failed to get previous placements", "guild", guildID, "error", err)
		return nil, err
	}

	var leaderList []*Leader
	for _, previousWin := range previousWins {
		leaderList = append(leaderList, &Leader{
			PlayerID: previousWin.PlayerID,
			Wins:     previousWin.Count,
		})
	}

	return leaderList, nil
}

func (s *Service) GetPlacementsForPeriod(ctx context.Context, tournamentID int32, roundDate time.Time) ([]*Leader, error) {
	previousPlacements, err := s.queries.GetPlacementsForPeriod(ctx, database.GetPlacementsForPeriodParams{
		TournamentID: tournamentID,
		RoundDate:    sql.NullTime{Time: roundDate, Valid: true},
	})
	if err == sql.ErrNoRows {
		slog.Warn("No previous placements", "tournament", tournamentID, "roundDate", roundDate, "error", "sql.ErrNoRows")
		return nil, err
	} else if err != nil {
		slog.Error("Failed to get previous placements", "tournament", tournamentID, "roundDate", roundDate, "error", err)
		return nil, err
	}

	var previousPlacementList []*Leader
	for _, previousPlacement := range previousPlacements {
		previousPlacementList = append(previousPlacementList, &Leader{
			PlayerID:     previousPlacement.PlayerID,
			TotalStrokes: previousPlacement.TotalStrokes,
		})
	}

	return previousPlacementList, nil
}
