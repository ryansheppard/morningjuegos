package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

type StandardDeviation struct {
	PlayerID int64
	StdDev   string
}

func (s *Service) GetStandardDeviation(ctx context.Context, reverse bool) ([]StandardDeviation, error) {
	performers, err := s.queries.GetStandardDeviation(ctx)
	if err == sql.ErrNoRows {
		slog.Info("No std dev found")
		return nil, err
	} else if err != nil {
		slog.Error("Failed to std dev", "error", err)
		return nil, err
	}

	if reverse {
		for i, j := 0, len(performers)-1; i < j; i, j = i+1, j-1 {
			performers[i], performers[j] = performers[j], performers[i]
		}
	}

	var stdDevs []StandardDeviation
	for _, performer := range performers {
		stdDevs = append(stdDevs, StandardDeviation{
			PlayerID: performer.PlayerID,
			StdDev:   performer.StandardDeviation,
		})
	}

	return stdDevs, nil
}

type HoleStats struct {
	Color   string
	Strokes float64
}

func (s *Service) GetHardestHole(ctx context.Context, tournamentID int32) (*HoleStats, error) {
	hardestHole, err := s.queries.GetHardestHole(ctx, tournamentID)
	if err != nil {
		slog.Error("Failed to get hardest hole", "error", err)
		return nil, err
	}

	return &HoleStats{
		Color:   hardestHole.Color,
		Strokes: hardestHole.Strokes,
	}, nil
}

func (s *Service) GetMostCommonHoleForNumber(ctx context.Context, tournamentID int32, holeNumber int32) (*HoleStats, error) {
	hole, err := s.queries.GetMostCommonHoleForNumber(ctx, database.GetMostCommonHoleForNumberParams{
		TournamentID: tournamentID,
		HoleNumber:   holeNumber,
	})
	if err != nil {
		slog.Error("Failed to get most common hole for number", "error", err)
		return nil, err
	}

	return &HoleStats{
		Color:   hole.Color,
		Strokes: float64(hole.Strokes),
	}, nil
}

type HoleInOneLeaders struct {
	PlayerID int64
	Count    int64
}

func (s *Service) GetHoleInOneLeaders(ctx context.Context, tournamentID int32) ([]*HoleInOneLeaders, error) {
	holeInOneLeaders, err := s.queries.GetHoleInOneLeaders(ctx, tournamentID)
	if err == sql.ErrNoRows {
		slog.Warn("No hole in one leaders", "tournament", tournamentID)
		return nil, err
	} else if err != nil {
		slog.Error("Failed to get hole in one leaders", "tournament", tournamentID, "error", err)
		return nil, err
	}

	var leaders []*HoleInOneLeaders
	for _, leader := range holeInOneLeaders {
		leaders = append(leaders, &HoleInOneLeaders{
			PlayerID: leader.PlayerID.Int64,
			Count:    leader.Count,
		})
	}

	return leaders, nil
}

func (s *Service) GetBestRounds(ctx context.Context, tournamentID int32) ([]*Leader, error) {
	rounds, err := s.queries.GetBestRounds(ctx, tournamentID)
	if err == sql.ErrNoRows {
		slog.Warn("No best rounds", "tournament", tournamentID)
		return nil, err
	} else if err != nil {
		slog.Error("Failed to get best rounds", "tournament", tournamentID, "error", err)
		return nil, err
	}

	var leaders []*Leader
	for _, round := range rounds {
		leaders = append(leaders, &Leader{
			PlayerID:     round.PlayerID,
			TotalStrokes: int64(round.TotalStrokes),
		})
	}

	return leaders, nil
}

func (s *Service) GetWorstRounds(ctx context.Context, tournamentID int32) ([]*Leader, error) {
	rounds, err := s.queries.GetWorstRounds(ctx, tournamentID)
	if err == sql.ErrNoRows {
		slog.Warn("No best rounds", "tournament", tournamentID)
		return nil, err
	} else if err != nil {
		slog.Error("Failed to get best rounds", "tournament", tournamentID, "error", err)
		return nil, err
	}

	var leaders []*Leader
	for _, round := range rounds {
		leaders = append(leaders, &Leader{
			PlayerID:     round.PlayerID,
			TotalStrokes: int64(round.TotalStrokes),
		})
	}

	return leaders, nil
}
