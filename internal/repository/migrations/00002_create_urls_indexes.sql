-- +goose Up
CREATE UNIQUE INDEX long_url_idx ON urls (long_url);

-- +goose Down
DROP INDEX long_url_idx;