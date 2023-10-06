package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

var ctx = context.TODO()

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

func TestGetTournament(t *testing.T) {
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
