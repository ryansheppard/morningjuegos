package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

func (s *Service) HasPlayedToday(ctx context.Context, playerID int64, tournamentID int32) (bool, error) {
	_, err := s.queries.HasPlayedToday(ctx, database.HasPlayedTodayParams{
		PlayerID:     playerID,
		TournamentID: tournamentID,
	})

	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		slog.Error("Failed to check if player has played today", "player", playerID, "tournament", tournamentID, "error", err)
		return false, err
	}

	return false, nil
}
func (s *Service) HasPlayed(ctx context.Context, playerID int64, tournamentID int32, day time.Time) (bool, error) {
	roundDate := sql.NullTime{
		Time:  day,
		Valid: true,
	}
	_, err := s.queries.HasPlayed(ctx, database.HasPlayedParams{
		PlayerID:     playerID,
		TournamentID: tournamentID,
		RoundDate:    roundDate,
	})
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		slog.Error("Failed to check if player has played", "player", playerID, "tournament", tournamentID, "day", day, "error", err)
		return false, err
	}

	return false, nil
}

func (s *Service) InsertRound(ctx context.Context, round *database.Round, holes []*database.Hole) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return false, err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	insertedRound, err := qtx.CreateRound(ctx, database.CreateRoundParams{
		TournamentID: round.TournamentID,
		PlayerID:     round.PlayerID,
		OriginalDate: round.OriginalDate,
		TotalStrokes: round.TotalStrokes,
		Percentage:   round.Percentage,
		FirstRound:   round.FirstRound,
		InsertedBy:   round.InsertedBy,
		RoundDate:    round.RoundDate,
	})

	if err != nil {
		slog.Error("Failed to insert round", "round", round, "error", err)
		return false, err
	}

	for _, hole := range holes {
		hole.RoundID = insertedRound.ID
		_, err = qtx.CreateHole(ctx, database.CreateHoleParams{
			RoundID:    hole.RoundID,
			Color:      hole.Color,
			Strokes:    hole.Strokes,
			HoleNumber: hole.HoleNumber,
			InsertedBy: hole.InsertedBy,
		})
		if err != nil {
			slog.Error("Failed to insert hole", "hole", hole, "error", err)
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return false, err
	}

	return true, nil
}
