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

-- name: GetRecentBaseball :one
SELECT *
FROM baseball
WHERE id = (SELECT MAX(id) FROM baseball);

-- name: GetRecentFootball :one
SELECT *
FROM football
WHERE id = (SELECT MAX(id) FROM football);

-- name: GetRecentSoccer :one
SELECT *
FROM soccer
WHERE id = (SELECT MAX(id) FROM soccer);
