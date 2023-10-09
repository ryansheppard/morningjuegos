package service

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

var ctx = context.TODO()

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestGetActiveTournament(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "guild_id", "start_time", "end_time", "inserted_by"}).
		AddRow(1, 1, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour), "test")
	mock.ExpectQuery("SELECT.+").WithArgs(1).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	tournament, err := service.GetActiveTournament(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting tournament: %s", err)
	}

	if tournament.ID != 1 {
		t.Errorf("expected tournament ID to be 1, got %d", tournament.ID)
	}
}

func TestGetActiveTournamentFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetActiveTournament(ctx, 1)

	if err == nil {
		t.Errorf("expected error while getting tournament, got nil")
	}
}

func TestGetInactiveTournaments(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "guild_id", "start_time", "end_time", "inserted_by"}).
		AddRow(1, 1, time.Now().Add(-96*time.Hour), time.Now().Add(-72*time.Hour), "test")
	mock.ExpectQuery("SELECT.+").WithArgs(1, AnyTime{}).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	tournament, err := service.GetInactiveTournaments(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting tournament: %s", err)
	}

	if len(tournament) != 1 {
		t.Errorf("expected tournament length to be 1, got %d", len(tournament))
	}

	if tournament[0].ID != 1 {
		t.Errorf("expected tournament ID to be 1, got %d", tournament[0].ID)
	}
}

func TestGetInactiveTournamentFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetInactiveTournaments(ctx, 1)

	if err == nil {
		t.Errorf("expected error while getting tournament, got nil")
	}
}

func TestCreateTournament(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	now := time.Now()
	start := time.Date(2023, 10, 1, 0, 0, 0, 0, now.Location())
	end := time.Date(2023, 10, 10, 23, 59, 59, 0, now.Location())
	mock.ExpectQuery("INSERT INTO tournament.+").
		WithArgs(1, start, end, "test").
		WillReturnRows(sqlmock.NewRows([]string{"id", "guild_id", "start_time", "end_time", "inserted_by"}).AddRow(1, 1, start, end, "test"))

	queries := database.New(d)
	service := New(d, queries)

	tournament, err := service.CreateTournament(ctx, 1, start, end, "test")
	if err != nil {
		t.Errorf("error was not expected while inserting tournament: %s", err)
	}

	if tournament.ID != 1 {
		t.Errorf("expected tournament ID to be 1, got %d", tournament.ID)
	}
}

func TestGetUniquePlayersInTournament(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	rows := sqlmock.NewRows([]string{"player_id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)
	mock.ExpectQuery("SELECT.+").WithArgs(1).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	players, err := service.GetUniquePlayersInTournament(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting players: %s", err)
	}

	if len(players) != 3 {
		t.Errorf("expected players length to be 3, got %d", len(players))
	}
}

func TestGetUniquePlayersFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetUniquePlayersInTournament(ctx, 1)

	if err == nil {
		t.Errorf("expected error while getting players, got nil")
	}
}

func TestGetTournamentPlacements(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	rows := sqlmock.NewRows([]string{"tournament_id", "player_id", "tournament_placement", "strokes", "inserted_at", "inserted_by"}).
		AddRow(1, 1, 1, 10, time.Now(), "test").
		AddRow(1, 2, 2, 11, time.Now(), "test").
		AddRow(1, 3, 3, 12, time.Now(), "test")
	mock.ExpectQuery("SELECT.+").WithArgs(1).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	placements, err := service.GetTournamentPlacements(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting placements: %s", err)
	}

	if len(placements) != 3 {
		t.Errorf("expected placements length to be 3, got %d", len(placements))
	}
}

func TestGetTournamentPlacementsFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetTournamentPlacements(ctx, 1)

	if err == nil {
		t.Errorf("expected error while getting placements, got nil")
	}
}

func TestGetFinalLeaders(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	rows := sqlmock.NewRows([]string{"player_id", "total_strokes"}).
		AddRow(1, 10).
		AddRow(2, 11).
		AddRow(3, 12)
	mock.ExpectQuery("SELECT.+").WithArgs(1).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	leaders, err := service.GetFinalLeaders(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting leaders: %s", err)
	}

	if len(leaders) != 3 {
		t.Errorf("expected leaders length to be 3, got %d", len(leaders))
	}
}

func TestGetFinalLeadersFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetFinalLeaders(ctx, 1)

	if err == nil {
		t.Errorf("expected error while getting leaders, got nil")
	}
}

func TestCleanTournamentPlacements(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer d.Close()
	mock.ExpectExec("DELETE FROM tournament_placement.+").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	queries := database.New(d)
	service := New(d, queries)

	err = service.CleanTournamentPlacements(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while cleaning placements: %s", err)
	}
}

func TestCleanTournamentPlacementsFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	queries := database.New(d)
	service := New(d, queries)

	err = service.CleanTournamentPlacements(ctx, 1)
	if err == nil {
		t.Errorf("expected error while cleaning placements, got nil")
	}
}

func TestCreateTournamentPlacement(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)

	}
	defer d.Close()

	mock.ExpectQuery("INSERT INTO tournament_placement.+").
		WithArgs(1, 1, 1, 10, "test").
		WillReturnRows(sqlmock.NewRows([]string{"tournament_id", "player_id", "tournament_placement", "strokes", "inserted_at", "inserted_by"}).AddRow(1, 1, 1, 1, time.Now(), "test"))

	queries := database.New(d)
	service := New(d, queries)

	err = service.CreateTournamentPlacement(ctx, 1, 1, 1, 10, "test")
	if err != nil {
		t.Errorf("error was not expected while creating placement: %s", err)
	}
}

func TestCreateTournamentPlacementFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)

	}
	defer d.Close()

	queries := database.New(d)
	service := New(d, queries)

	err = service.CreateTournamentPlacement(ctx, 1, 1, 1, 10, "test")
	if err == nil {
		t.Errorf("expected error while creating placement, got nil")
	}
}

func TestGetLeaders(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"player_id", "total_strokes"}).
		AddRow(1, 10).
		AddRow(2, 11).
		AddRow(3, 12)
	mock.ExpectQuery("SELECT.+").WithArgs(1).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	leaders, err := service.GetLeaders(ctx, 1)
	if err != nil {
		t.Errorf("error was not expected while getting leaders: %s", err)
	}

	if len(leaders) != 3 {
		t.Errorf("expected leaders length to be 3, got %d", len(leaders))
	}
}

func TestGetLeadersFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetLeaders(ctx, 1)
	if err == nil {
		t.Errorf("expected error while getting leaders, got nil")
	}
}

func TestGetTournamentPlacementsByPosition(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"count", "tournament_placement", "player_id"}).
		AddRow(2, 2, 1).
		AddRow(2, 2, 2)
	mock.ExpectQuery("SELECT.+").WithArgs(1, 2).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	placements, err := service.GetTournamentPlacementsByPosition(ctx, 1, 2)
	if err != nil {
		t.Errorf("error was not expected while getting placements: %s", err)
	}

	if len(placements) != 2 {
		t.Errorf("expected placements to be 2, got %d", len(placements))
	}
}

func TestGetTournamentPlacementsByPositionFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetTournamentPlacementsByPosition(ctx, 1, 1)
	if err == nil {
		t.Errorf("expected error while getting placements, got nil")
	}
}

func TestGetPlacementsForPeriod(t *testing.T) {
	t.Parallel()

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"total_strokes", "player_id"}).
		AddRow(12, 1).
		AddRow(14, 2).
		AddRow(13, 3)
	mock.ExpectQuery("SELECT.+").WithArgs(1, AnyTime{}).WillReturnRows(rows)

	queries := database.New(d)
	service := New(d, queries)

	placements, err := service.GetPlacementsForPeriod(ctx, 1, time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Errorf("error was not expected while getting placements: %s", err)
	}

	if len(placements) != 3 {
		t.Errorf("expected placements to be 2, got %d", len(placements))
	}
}

func TestGetPlacementsFails(t *testing.T) {
	t.Parallel()

	d, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	queries := database.New(d)
	service := New(d, queries)

	_, err = service.GetPlacementsForPeriod(ctx, 1, time.Now().Add(-24*time.Hour))
	if err == nil {
		t.Errorf("expected error while getting placements, got nil")
	}
}
