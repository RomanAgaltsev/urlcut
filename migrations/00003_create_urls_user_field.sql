-- +goose Up
-- +goose StatementBegin
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE urls
    ADD COLUMN uid UUID NOT NULL DEFAULT uuid_generate_v4();

CREATE INDEX uid_idx ON urls (uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uid_idx;

ALTER TABLE urls DROP COLUMN uid;

DROP
EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd