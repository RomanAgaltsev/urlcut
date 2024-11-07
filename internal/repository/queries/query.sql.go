// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"
)

const createURL = `-- name: CreateURL :one
INSERT INTO urls (long_url, base_url, url_id)
VALUES ($1, $2, $3) RETURNING id, long_url, base_url, url_id, created_at
`

type CreateURLParams struct {
	LongUrl string
	BaseUrl string
	UrlID   string
}

func (q *Queries) CreateURL(ctx context.Context, arg CreateURLParams) (Url, error) {
	row := q.db.QueryRowContext(ctx, createURL, arg.LongUrl, arg.BaseUrl, arg.UrlID)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.LongUrl,
		&i.BaseUrl,
		&i.UrlID,
		&i.CreatedAt,
	)
	return i, err
}

const getURL = `-- name: GetURL :one
SELECT id, long_url, base_url, url_id, created_at
FROM urls
WHERE url_id = $1 LIMIT 1
`

func (q *Queries) GetURL(ctx context.Context, urlID string) (Url, error) {
	row := q.db.QueryRowContext(ctx, getURL, urlID)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.LongUrl,
		&i.BaseUrl,
		&i.UrlID,
		&i.CreatedAt,
	)
	return i, err
}
