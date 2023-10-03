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

const cleanTournamentPlacements = `-- name: CleanTournamentPlacements :exec
DELETE FROM tournament_placement WHERE tournament_id = $1
`

func (q *Queries) CleanTournamentPlacements(ctx context.Context, tournamentID int32) error {
	_, err := q.db.ExecContext(ctx, cleanTournamentPlacements, tournamentID)
	return err
}

const createHole = `-- name: CreateHole :one
INSERT INTO hole (round_id, color, strokes, hole_number, inserted_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, round_id, hole_number, color, strokes, inserted_at, inserted_by
`

type CreateHoleParams struct {
	RoundID    int32
	Color      string
	Strokes    int32
	HoleNumber int32
	InsertedBy string
}

// Hole Queries
func (q *Queries) CreateHole(ctx context.Context, arg CreateHoleParams) (Hole, error) {
	row := q.db.QueryRowContext(ctx, createHole,
		arg.RoundID,
		arg.Color,
		arg.Strokes,
		arg.HoleNumber,
		arg.InsertedBy,
	)
	var i Hole
	err := row.Scan(
		&i.ID,
		&i.RoundID,
		&i.HoleNumber,
		&i.Color,
		&i.Strokes,
		&i.InsertedAt,
		&i.InsertedBy,
	)
	return i, err
}

const createRound = `-- name: CreateRound :one
INSERT INTO round
(tournament_id, player_id, total_strokes, original_date, percentage, first_round, inserted_by, round_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage, inserted_by, round_date
`

type CreateRoundParams struct {
	TournamentID int32
	PlayerID     int64
	TotalStrokes int32
	OriginalDate string
	Percentage   string
	FirstRound   bool
	InsertedBy   string
	RoundDate    sql.NullTime
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
		arg.InsertedBy,
		arg.RoundDate,
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
		&i.InsertedBy,
		&i.RoundDate,
	)
	return i, err
}

const createTournament = `-- name: CreateTournament :one
INSERT INTO tournament (guild_id, start_time, end_time, inserted_by) VALUES ($1, $2, $3, $4) RETURNING id, guild_id, start_time, end_time, inserted_at, inserted_by
`

type CreateTournamentParams struct {
	GuildID    int64
	StartTime  time.Time
	EndTime    time.Time
	InsertedBy string
}

func (q *Queries) CreateTournament(ctx context.Context, arg CreateTournamentParams) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, createTournament,
		arg.GuildID,
		arg.StartTime,
		arg.EndTime,
		arg.InsertedBy,
	)
	var i Tournament
	err := row.Scan(
		&i.ID,
		&i.GuildID,
		&i.StartTime,
		&i.EndTime,
		&i.InsertedAt,
		&i.InsertedBy,
	)
	return i, err
}

const createTournamentPlacement = `-- name: CreateTournamentPlacement :one
INSERT INTO tournament_placement (tournament_id, player_id, tournament_placement, strokes, inserted_by) VALUES ($1, $2, $3, $4, $5) RETURNING tournament_id, player_id, tournament_placement, strokes, inserted_at, inserted_by
`

type CreateTournamentPlacementParams struct {
	TournamentID        int32
	PlayerID            int64
	TournamentPlacement int32
	Strokes             int32
	InsertedBy          string
}

// TournamentPlacement Queries
func (q *Queries) CreateTournamentPlacement(ctx context.Context, arg CreateTournamentPlacementParams) (TournamentPlacement, error) {
	row := q.db.QueryRowContext(ctx, createTournamentPlacement,
		arg.TournamentID,
		arg.PlayerID,
		arg.TournamentPlacement,
		arg.Strokes,
		arg.InsertedBy,
	)
	var i TournamentPlacement
	err := row.Scan(
		&i.TournamentID,
		&i.PlayerID,
		&i.TournamentPlacement,
		&i.Strokes,
		&i.InsertedAt,
		&i.InsertedBy,
	)
	return i, err
}

const getActiveTournament = `-- name: GetActiveTournament :one

SELECT id, guild_id, start_time, end_time, inserted_at, inserted_by FROM tournament WHERE guild_id = $1 AND start_time <= $2 AND end_time >= $2
`

type GetActiveTournamentParams struct {
	GuildID   int64
	StartTime time.Time
}

// Player Queries
// Tournament Queries
func (q *Queries) GetActiveTournament(ctx context.Context, arg GetActiveTournamentParams) (Tournament, error) {
	row := q.db.QueryRowContext(ctx, getActiveTournament, arg.GuildID, arg.StartTime)
	var i Tournament
	err := row.Scan(
		&i.ID,
		&i.GuildID,
		&i.StartTime,
		&i.EndTime,
		&i.InsertedAt,
		&i.InsertedBy,
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

const getFinalLeaders = `-- name: GetFinalLeaders :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
GROUP BY player_id
ORDER BY total_strokes ASC
`

type GetFinalLeadersRow struct {
	TotalStrokes int64
	PlayerID     int64
}

func (q *Queries) GetFinalLeaders(ctx context.Context, tournamentID int32) ([]GetFinalLeadersRow, error) {
	rows, err := q.db.QueryContext(ctx, getFinalLeaders, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFinalLeadersRow
	for rows.Next() {
		var i GetFinalLeadersRow
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

const getInactiveTournaments = `-- name: GetInactiveTournaments :many
SELECT id, guild_id, start_time, end_time, inserted_at, inserted_by FROM tournament WHERE guild_id = $1 AND end_time < $2
`

type GetInactiveTournamentsParams struct {
	GuildID int64
	EndTime time.Time
}

func (q *Queries) GetInactiveTournaments(ctx context.Context, arg GetInactiveTournamentsParams) ([]Tournament, error) {
	rows, err := q.db.QueryContext(ctx, getInactiveTournaments, arg.GuildID, arg.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tournament
	for rows.Next() {
		var i Tournament
		if err := rows.Scan(
			&i.ID,
			&i.GuildID,
			&i.StartTime,
			&i.EndTime,
			&i.InsertedAt,
			&i.InsertedBy,
		); err != nil {
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

const getLeaders = `-- name: GetLeaders :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
AND round_date >= $2
AND round_date <= $3
GROUP BY player_id
ORDER BY total_strokes ASC
`

type GetLeadersParams struct {
	TournamentID int32
	RoundDate    sql.NullTime
	RoundDate_2  sql.NullTime
}

type GetLeadersRow struct {
	TotalStrokes int64
	PlayerID     int64
}

func (q *Queries) GetLeaders(ctx context.Context, arg GetLeadersParams) ([]GetLeadersRow, error) {
	rows, err := q.db.QueryContext(ctx, getLeaders, arg.TournamentID, arg.RoundDate, arg.RoundDate_2)
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
AND round_date < $2
GROUP BY player_id
ORDER BY total_strokes ASC
`

type GetPlacementsForPeriodParams struct {
	TournamentID int32
	RoundDate    sql.NullTime
}

type GetPlacementsForPeriodRow struct {
	TotalStrokes int64
	PlayerID     int64
}

func (q *Queries) GetPlacementsForPeriod(ctx context.Context, arg GetPlacementsForPeriodParams) ([]GetPlacementsForPeriodRow, error) {
	rows, err := q.db.QueryContext(ctx, getPlacementsForPeriod, arg.TournamentID, arg.RoundDate)
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

const getStandardDeviation = `-- name: GetStandardDeviation :many
SELECT player_id, round(stddev(total_strokes), 3) as standard_deviation
FROM round
WHERE inserted_by = 'parser'
AND first_round = 't'
GROUP BY player_id
ORDER BY standard_deviation
`

type GetStandardDeviationRow struct {
	PlayerID          int64
	StandardDeviation string
}

// Stats
func (q *Queries) GetStandardDeviation(ctx context.Context) ([]GetStandardDeviationRow, error) {
	rows, err := q.db.QueryContext(ctx, getStandardDeviation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStandardDeviationRow
	for rows.Next() {
		var i GetStandardDeviationRow
		if err := rows.Scan(&i.PlayerID, &i.StandardDeviation); err != nil {
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

const getTournamentPlacements = `-- name: GetTournamentPlacements :many
SELECT tournament_id, player_id, tournament_placement, strokes, inserted_at, inserted_by FROM tournament_placement WHERE tournament_id = $1
`

func (q *Queries) GetTournamentPlacements(ctx context.Context, tournamentID int32) ([]TournamentPlacement, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentPlacements, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TournamentPlacement
	for rows.Next() {
		var i TournamentPlacement
		if err := rows.Scan(
			&i.TournamentID,
			&i.PlayerID,
			&i.TournamentPlacement,
			&i.Strokes,
			&i.InsertedAt,
			&i.InsertedBy,
		); err != nil {
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

const getTournamentPlacementsByPosition = `-- name: GetTournamentPlacementsByPosition :one
SELECT COUNT(*) AS count, tournament_placement
FROM tournament_placement
LEFT JOIN tournament ON tournament_placement.tournament_id = tournament.id
WHERE tournament.guild_id = $1
AND player_id = $2
AND tournament_placement = $3
GROUP BY tournament_placement, player_id
`

type GetTournamentPlacementsByPositionParams struct {
	GuildID             int64
	PlayerID            int64
	TournamentPlacement int32
}

type GetTournamentPlacementsByPositionRow struct {
	Count               int64
	TournamentPlacement int32
}

func (q *Queries) GetTournamentPlacementsByPosition(ctx context.Context, arg GetTournamentPlacementsByPositionParams) (GetTournamentPlacementsByPositionRow, error) {
	row := q.db.QueryRowContext(ctx, getTournamentPlacementsByPosition, arg.GuildID, arg.PlayerID, arg.TournamentPlacement)
	var i GetTournamentPlacementsByPositionRow
	err := row.Scan(&i.Count, &i.TournamentPlacement)
	return i, err
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
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage, inserted_by, round_date
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
		&i.InsertedBy,
		&i.RoundDate,
	)
	return i, err
}

const hasPlayed = `-- name: HasPlayed :one
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage, inserted_by, round_date
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND round_date = $3
`

type HasPlayedParams struct {
	PlayerID     int64
	TournamentID int32
	RoundDate    sql.NullTime
}

func (q *Queries) HasPlayed(ctx context.Context, arg HasPlayedParams) (Round, error) {
	row := q.db.QueryRowContext(ctx, hasPlayed, arg.PlayerID, arg.TournamentID, arg.RoundDate)
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
		&i.InsertedBy,
		&i.RoundDate,
	)
	return i, err
}

const hasPlayedToday = `-- name: HasPlayedToday :one
SELECT id, tournament_id, player_id, total_strokes, original_date, inserted_at, first_round, percentage, inserted_by, round_date 
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND round_date = CURRENT_DATE
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
		&i.InsertedBy,
		&i.RoundDate,
	)
	return i, err
}
