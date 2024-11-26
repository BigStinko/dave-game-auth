-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserMatches :many
SELECT
	m.*,
	px.username as player_x_username,
	py.username as player_y_username,
	w.username as winner_username,
	CASE
		WHEN m.winner_id = $1 THEN true
		ELSE false
	END as won
FROM matches m
JOIN users px ON m.player_x_id = px.id
JOIN users py ON m.player_y_id = py.id
JOIN users w ON m.winner_id = w.id
WHERE m.player_x_id = $1 OR m.player_y_id = $1
ORDER BY m.played_at DESC
LIMIT 10;

-- name: GetUserStats :one
SELECT
	COUNT(*) as total_matches,
	COUNT(CASE WHEN winner_id = $1 THEN 1 END) as wins,
	COUNT(CASE WHEN (player_x_id = $1 OR player_y_id = $1) AND winner_id != $1 THEN 1 END) as losses
FROM matches
WHERE player_x_id = $1 OR player_y_id = $1;
