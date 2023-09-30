-- Guild Queries

-- name: GetAllGuilds :many
SELECT DISTINCT guild_id FROM tournament;

-- Player Queries

-- Tournament Queries
-- name: GetActiveTournament :one
SELECT * FROM tournament WHERE guild_id = $1 AND start_time <= $2 AND end_time >= $2;

-- name: GetInactiveTournaments :many
SELECT * FROM tournament WHERE guild_id = $1 AND end_time < $2;

-- name: CreateTournament :one
INSERT INTO tournament (guild_id, start_time, end_time, inserted_by) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUniquePlayersInTournament :many
SELECT DISTINCT player_id FROM round WHERE tournament_id = $1;

-- TournamentPlacement Queries
-- name: CreateTournamentPlacement :one
INSERT INTO tournament_placement (tournament_id, player_id, tournament_placement, strokes, inserted_by) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CleanTournamentPlacements :exec
DELETE FROM tournament_placement WHERE tournament_id = $1;

-- name: GetTournamentPlacements :many
SELECT * FROM tournament_placement WHERE tournament_id = $1;

-- name: GetTournamentPlacementsByPosition :one
SELECT COUNT(*) AS count, tournament_placement
FROM tournament_placement
LEFT JOIN tournament ON tournament_placement.tournament_id = tournament.id
WHERE tournament.guild_id = $1
AND player_id = $2
AND tournament_placement = $3
GROUP BY tournament_placement, player_id;

-- Round Queries
-- name: CreateRound :one
INSERT INTO round
(tournament_id, player_id, total_strokes, original_date, percentage, first_round, inserted_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: HasPlayedToday :one
SELECT * 
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND date_trunc('day', inserted_at) = date_trunc('day', NOW());

-- name: HasPlayed :one
SELECT *
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND date_trunc('day', inserted_at) = date_trunc('day', $3);

-- name: GetLeaders :many
SELECT SUM(total_strokes) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
AND inserted_at > $2
AND inserted_at < $3
GROUP BY player_id
ORDER BY total_strokes ASC;

-- name: GetFinalLeaders :many
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
AND inserted_at < $2
GROUP BY player_id
ORDER BY total_strokes ASC;

-- name: GetWorstRound :one
SELECT *
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
ORDER BY total_strokes DESC
LIMIT 1;

-- Hole Queries
-- name: CreateHole :one
INSERT INTO hole (round_id, color, strokes, hole_number, inserted_by) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetHardestHole :one
SELECT AVG(strokes) AS strokes, color
FROM hole 
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
GROUP BY color
ORDER BY strokes DESC
LIMIT 1;

-- name: GetMostCommonHoleForNumber :one
SELECT COUNT(color) AS strokes, color
FROM hole
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
AND hole_number = $2
GROUP BY color
ORDER BY strokes DESC
LIMIT 1;

-- name: GetHoleInOneLeader :one
SELECT COUNT(*) AS count, round.player_id
FROM hole
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
AND round.player_id IS NOT NULL
AND hole.strokes = 1
GROUP BY round.player_id
ORDER BY count DESC
LIMIT 1;