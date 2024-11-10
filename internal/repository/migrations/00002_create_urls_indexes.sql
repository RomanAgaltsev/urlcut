-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS url_id_idx ON urls (url_id);

CREATE UNIQUE INDEX IF NOT EXISTS long_url_idx ON urls (long_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS long_url_idx;

DROP INDEX IF EXISTS url_id_idx;
-- +goose StatementEnd