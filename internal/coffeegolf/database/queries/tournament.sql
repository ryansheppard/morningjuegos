-- Tournament Queries
-- name: GetActiveTournament :one
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND start_time <= NOW() AND end_time >= NOW();

-- name: GetTournamentByForDate :one
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND start_time <= $2 AND end_time >= $2;

-- name: GetInactiveTournaments :many
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND end_time < $2;

-- name: GetTournament :one
SELECT * FROM tournament WHERE id = $1;

-- name: CreateTournament :one
INSERT INTO tournament (guild_id, start_time, end_time, inserted_by) VALUES ($1, $2, $3, $4) RETURNING id, guild_id, start_time, end_time, inserted_by;

-- name: GetUniquePlayersInTournament :many
SELECT DISTINCT player_id FROM round WHERE tournament_id = $1;

-- TournamentPlacement Queries
-- name: CreateTournamentPlacement :one
INSERT INTO tournament_placement (tournament_id, player_id, tournament_placement, strokes, inserted_by) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CleanTournamentPlacements :exec
DELETE FROM tournament_placement WHERE tournament_id = $1;

-- name: GetTournamentPlacements :many
SELECT * FROM tournament_placement WHERE tournament_id = $1;

-- name: GetTournamentPlacementsByPosition :many
SELECT COUNT(*) AS count, tournament_placement, player_id
FROM tournament_placement
LEFT JOIN tournament ON tournament_placement.tournament_id = tournament.id
WHERE tournament.guild_id = $1
AND tournament_placement = $2
GROUP BY tournament_placement, player_id;

-- name: GetLeaders :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
GROUP BY player_id
ORDER BY total_strokes ASC;

-- name: GetPlacementsForPeriod :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
AND round_date < $2
GROUP BY player_id
ORDER BY total_strokes ASC;