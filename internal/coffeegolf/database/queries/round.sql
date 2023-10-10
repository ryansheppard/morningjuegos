-- Round Queries
-- name: CreateRound :one
INSERT INTO round
(tournament_id, player_id, total_strokes, original_date, percentage, first_round, inserted_by, round_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: HasPlayedToday :one
SELECT * 
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND round_date = CURRENT_DATE
AND first_round = TRUE;

-- name: HasPlayed :one
SELECT *
FROM round
WHERE player_id = $1
AND tournament_id = $2
AND round_date = $3
AND first_round = TRUE;
