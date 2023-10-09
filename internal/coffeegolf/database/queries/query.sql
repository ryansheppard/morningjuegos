-- Guild Queries

-- name: GetAllGuilds :many
SELECT DISTINCT guild_id FROM tournament;

-- Player Queries

-- Tournament Queries
-- name: GetActiveTournament :one
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND start_time <= NOW() AND end_time >= NOW();

-- name: GetTournamentByForDate :one
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND start_time <= $2 AND end_time >= $2;

-- name: GetInactiveTournaments :many
SELECT id, guild_id, start_time, end_time, inserted_by FROM tournament WHERE guild_id = $1 AND end_time < $2;

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