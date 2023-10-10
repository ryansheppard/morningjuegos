-- Stats queries
-- name: GetWorstRounds :many
SELECT CAST(MAX(total_strokes) AS INTEGER) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
GROUP BY player_id
ORDER BY total_strokes DESC;

-- name: GetBestRounds :many
SELECT CAST(MIN(total_strokes) AS INTEGER) AS total_strokes, player_id
FROM round
WHERE tournament_id = $1
AND first_round = TRUE
GROUP BY player_id
ORDER BY total_strokes ASC;

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

-- name: GetHoleInOneLeaders :many
SELECT COUNT(*) AS count, round.player_id
FROM hole
LEFT JOIN round ON hole.round_id = round.id
WHERE round.tournament_id = $1
AND round.first_round = TRUE
AND round.player_id IS NOT NULL
AND hole.strokes = 1
GROUP BY round.player_id
ORDER BY count DESC;

-- Stats
-- name: GetStandardDeviation :many
SELECT player_id, round(stddev(total_strokes), 3) as standard_deviation
FROM round
WHERE inserted_by = 'parser'
AND first_round = 't'
GROUP BY player_id
ORDER BY standard_deviation;