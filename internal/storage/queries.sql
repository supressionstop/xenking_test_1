-- name: SaveBaseball :one
INSERT INTO baseball (
    sport, rate, created_at
) VALUES (
             $1, $2, $3
         )
RETURNING *;
-- name: SaveFootball :one
INSERT INTO football (
    sport, rate, created_at
) VALUES (
             $1, $2, $3
         )
RETURNING *;
-- name: SaveSoccer :one
INSERT INTO soccer (
    sport, rate, created_at
) VALUES (
             $1, $2, $3
         )
RETURNING *;

