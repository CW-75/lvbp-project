-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListActiveGames :many
SELECT * FROM games
WHERE status = 'IN_PROGRESS'
ORDER BY start_time DESC;

-- name: CreateGame :one
INSERT INTO games (
  home_team_id, away_team_id, start_time
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateGame :one
UPDATE games
SET 
  status = $2,
  home_score = $3,
  away_score = $4,
  current_inning = $5,
  is_top_inning = $6
WHERE id = $1
RETURNING *;

-- name: CreateAtBat :one
INSERT INTO at_bats (
  game_id, inning, is_top_inning, batter_id, pitcher_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateAtBatResult :one
UPDATE at_bats
SET 
  result = $2,
  runs_scored = $3,
  outs_recorded = $4
WHERE id = $1
RETURNING *;

-- name: CreatePitch :one
INSERT INTO pitches (
  at_bat_id, pitch_number, coordinate_x, coordinate_y, 
  pitch_result, balls_before, strikes_before, outs_before, velocity_mph
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetCurrentAtBat :one
SELECT * FROM at_bats
WHERE game_id = $1
ORDER BY created_at DESC LIMIT 1;

-- name: GetPitchesForAtBat :many
SELECT * FROM pitches
WHERE at_bat_id = $1
ORDER BY pitch_number ASC;

-- name: GetTeam :one
SELECT * FROM teams
WHERE id = $1 LIMIT 1;

-- name: ListTeams :many
SELECT * FROM teams
ORDER BY name;

-- name: GetOutsForInning :one
SELECT COALESCE(SUM(outs_recorded), 0)::int FROM at_bats
WHERE game_id = $1 AND inning = $2 AND is_top_inning = $3;
