-- name: GetURL :one
SELECT *
FROM urls
WHERE url_id = $1 LIMIT 1;

-- name: GetURLByLong :one
SELECT *
FROM urls
WHERE long_url = $1 LIMIT 1;

-- name: StoreURL :one
INSERT INTO urls (long_url, base_url, url_id)
VALUES ($1, $2, $3) RETURNING *;