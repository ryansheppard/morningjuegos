// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createHole = `-- name: CreateHole :one
INSERT INTO hole (round_id, color, strokes, hole_number) VALUES ($1, $2, $3, $4) RETURNING id, round_id, hole_number, color, strokes, inserted_at
`

type CreateHoleParams struct {
	RoundID    int32
	Color      string
	Strokes    int32
	HoleNumber int32
}

// Hole Queries
func (q *Queries) CreateHole(ctx context.Context, arg CreateHoleParams) (Hole, error) {
	row := q.db.QueryRowContext(ctx, createHole,
		arg.RoundID,
		arg.Color,
		arg.Strokes,
		arg.HoleNumber,
	)
	var i Hole
	err := row.Scan(
		&i.ID,
		&i.RoundID,
		&i.HoleNumber,
		&i.Color,
		&i.Strokes,
		&i.InsertedAt,
	)
	return i, err
}

const createRound = `-- name: CreateRound :one
INSERT INTO round (tournament_id, player_id, total_strokes, original_date, percentage, first_round) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage
`

type CreateRoundParams struct {
	TournamentID int32
	PlayerID     int64
	TotalStrokes int32
	OriginalDate string
	Percentage   string
	FirstRound   bool
}

// Round Queries
func (q *Queries) CreateRound(ctx context.Context, arg CreateRoundParams) (Round, error) {
	row := q.db.QueryRowContext(ctx, createRound,
		arg.TournamentID,
		arg.PlayerID,
		arg.TotalStrokes,
		arg.OriginalDate,
		arg.Percentage,
		arg.FirstRound,
	)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.PlayerID,
		&i.TotalStrokes,
		&i.OriginalDate,
		&i.InsertedAt,
		&i.FirstRound,
		&i.Percentage,
	)
	return i, err
}

const createTournament = `-- name: CreateTournament :one
INSERT INTO tournament (guild_id, start_time, end_time) VALUES ($1, $2, $3) RETURNING id, guild_id, start_time, end_time, inserted_at
`

type CreateTournamentParams struct {
	GuildID   int64
	StartTime time.Time
	EndTime   time.Time
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, createTournament, arg.GuildID, arg.StartTime, arg.EndTime)
	var i Tournament
	err := row.Scan(
		&i.ID,
		&i.GuildID,
		&i.StartTime,
		&i.EndTime,
		&i.InsertedAt,
	)
	return i, err
}

const getActiveTournament = `-- name: GetActiveTournament :one

SELECT id, guild_id, start_time, end_time, inserted_at FROM tournament WHERE guild_id = $1 AND start_time < NOW() AND end_time > NOW()
`

// Player Queries
// Tournament Queries
func (q *Queries) GetActiveTournament(ctx context.Context, guildID int64) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, getActiveTournament, guildID)
	var i Tournament
	err := row.Scan(
		&i.ID,
		&i.GuildID,
		&i.StartTime,
		&i.EndTime,
		&i.InsertedAt,
	)
	return i, err
}

const getAllGuilds = `-- name: GetAllGuilds :many

SELECT DISTINCT guild_id FROM tournament
`

// Guild Queries
func (q *Queries) GetAllGuilds(ctx context.Context) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getAllGuilds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var guild_id int64
		if err := rows.Scan(&guild_id); err != nil {
			return nil, err
		}
		items = append(items, guild_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getHardestHole = `-- name: GetHardestHole :one
SELECT AVG(strokes) AS strokes, color
FROM hole 
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
GROUP BY color
ORDER BY strokes DESC
LIMIT 1
`

type GetHardestHoleRow struct {
	Strokes float64
	Color   string
}

func (q *Queries) GetHardestHole(ctx context.Context, tournamentID int32) (GetHardestHoleRow, error) {
	row := q.db.QueryRowContext(ctx, getHardestHole, tournamentID)
	var i GetHardestHoleRow
	err := row.Scan(&i.Strokes, &i.Color)
	return i, err
}

const getHoleInOneLeader = `-- name: GetHoleInOneLeader :one
SELECT COUNT(*) AS count, round.player_id
FROM hole
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
AND round.player_id IS NOT NULL
AND hole.strokes = 1
GROUP BY round.player_id
ORDER BY count DESC
LIMIT 1
`

type GetHoleInOneLeaderRow struct {
	Count    int64
	PlayerID sql.NullInt64
}

func (q *Queries) GetHoleInOneLeader(ctx context.Context, tournamentID int32) (GetHoleInOneLeaderRow, error) {
	row := q.db.QueryRowContext(ctx, getHoleInOneLeader, tournamentID)
	var i GetHoleInOneLeaderRow
	err := row.Scan(&i.Count, &i.PlayerID)
	return i, err
}

const getLeaders = `-- name: GetLeaders :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
AND inserted_at > $2
AND inserted_at < $3
GROUP BY player_id
ORDER BY total_strokes ASC
`

type GetLeadersParams struct {
	TournamentID int32
	InsertedAt   time.Time
	InsertedAt_2 time.Time
}

type GetLeadersRow struct {
	TotalStrokes int64
	PlayerID     int64
}

func (q *Queries) GetLeaders(ctx context.Context, arg GetLeadersParams) ([]GetLeadersRow, error) {
	rows, err := q.db.QueryContext(ctx, getLeaders, arg.TournamentID, arg.InsertedAt, arg.InsertedAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeadersRow
	for rows.Next() {
		var i GetLeadersRow
		if err := rows.Scan(&i.TotalStrokes, &i.PlayerID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMostCommonHoleForNumber = `-- name: GetMostCommonHoleForNumber :one
SELECT COUNT(color) AS strokes, color
FROM hole
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
AND hole_number = $2
GROUP BY color
ORDER BY strokes DESC
LIMIT 1
`

type GetMostCommonHoleForNumberParams struct {
	TournamentID int32
	HoleNumber   int32
}

type GetMostCommonHoleForNumberRow struct {
	Strokes int64
	Color   string
}

func (q *Queries) GetMostCommonHoleForNumber(ctx context.Context, arg GetMostCommonHoleForNumberParams) (GetMostCommonHoleForNumberRow, error) {
	row := q.db.QueryRowContext(ctx, getMostCommonHoleForNumber, arg.TournamentID, arg.HoleNumber)
	var i GetMostCommonHoleForNumberRow
	err := row.Scan(&i.Strokes, &i.Color)
	return i, err
}

const getPlacementsForPeriod = `-- name: GetPlacementsForPeriod :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
AND inserted_at < $2
GROUP BY player_id
ORDER BY total_strokes ASC
`

type GetPlacementsForPeriodParams struct {
	TournamentID int32
	InsertedAt   time.Time
}

type GetPlacementsForPeriodRow struct {
	TotalStrokes int64
	PlayerID     int64
}

func (q *Queries) GetPlacementsForPeriod(ctx context.Context, arg GetPlacementsForPeriodParams) ([]GetPlacementsForPeriodRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlacementsForPeriod, arg.TournamentID, arg.InsertedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPlacementsForPeriodRow
	for rows.Next() {
		var i GetPlacementsForPeriodRow
		if err := rows.Scan(&i.TotalStrokes, &i.PlayerID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUniquePlayersInTournament = `-- name: GetUniquePlayersInTournament :many
SELECT DISTINCT player_id FROM round WHERE tournament_id = $1
`

func (q *Queries) GetUniquePlayersInTournament(ctx context.Context, tournamentID int32) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getUniquePlayersInTournament, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var player_id int64
		if err := rows.Scan(&player_id); err != nil {
			return nil, err
		}
		items = append(items, player_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorstRound = `-- name: GetWorstRound :one
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
ORDER BY total_strokes DESC
LIMIT 1
`

func (q *Queries) GetWorstRound(ctx context.Context, tournamentID int32) (Round, error) {
	row := q.db.QueryRowContext(ctx, getWorstRound, tournamentID)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.PlayerID,
		&i.TotalStrokes,
		&i.OriginalDate,
		&i.InsertedAt,
		&i.FirstRound,
		&i.Percentage,
	)
	return i, err
}

const hasPlayed = `-- name: HasPlayed :one
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND date_trunc('day', inserted_at) = date_trunc('day', $3)
`

type HasPlayedParams struct {
	PlayerID     int64
	TournamentID int32
	DateTrunc    int64
}

func (q *Queries) HasPlayed(ctx context.Context, arg HasPlayedParams) (Round, error) {
	row := q.db.QueryRowContext(ctx, hasPlayed, arg.PlayerID, arg.TournamentID, arg.DateTrunc)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.PlayerID,
		&i.TotalStrokes,
		&i.OriginalDate,
		&i.InsertedAt,
		&i.FirstRound,
		&i.Percentage,
	)
	return i, err
}

const hasPlayedToday = `-- name: HasPlayedToday :one
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage 
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND date_trunc('day', inserted_at) = date_trunc('day', NOW())
`

type HasPlayedTodayParams struct {
	PlayerID     int64
	TournamentID int32
}

func (q *Queries) HasPlayedToday(ctx context.Context, arg HasPlayedTodayParams) (Round, error) {
	row := q.db.QueryRowContext(ctx, hasPlayedToday, arg.PlayerID, arg.TournamentID)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.PlayerID,
		&i.TotalStrokes,
		&i.OriginalDate,
		&i.InsertedAt,
		&i.FirstRound,
		&i.Percentage,
	)
	return i, err
}
