-- +goose Up
-- +goose StatementBegin
DROP INDEX long_url_idx;

ALTER TABLE urls
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX long_url_del_idx ON urls (long_url, is_deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX long_url_del_idx;

ALTER TABLE urls DROP COLUMN is_deleted;

CREATE UNIQUE INDEX long_url_idx ON urls (long_url);
-- +goose StatementEnd