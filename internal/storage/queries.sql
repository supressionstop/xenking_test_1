-- name: SaveLine :one
INSERT INTO lines (
    sport, rate, created_at
) VALUES (
             $1, $2, $3
         )
RETURNING *;