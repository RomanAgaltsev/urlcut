-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls
(
    id
    serial
    PRIMARY
    KEY,
    long_url
    TEXT
    NOT
    NULL,
    base_url
    VARCHAR
(
    100
) NOT NULL,
    url_id VARCHAR
(
    20
) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW
(
)
    );

CREATE UNIQUE INDEX IF NOT EXISTS url_id_idx ON urls (url_id, long_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS url_id_idx;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd