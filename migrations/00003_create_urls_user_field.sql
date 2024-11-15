-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
    ADD COLUMN uid UUID;

CREATE UNIQUE INDEX uid_idx ON urls (uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uid_idx;

ALTER TABLE urls DROP COLUMN uid;
-- +goose StatementEnd