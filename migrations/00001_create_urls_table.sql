-- +goose Up
CREATE TABLE
    urls (
        id SERIAL PRIMARY KEY,
        long_url TEXT NOT NULL,
        base_url VARCHAR(100) NOT NULL,
        url_id VARCHAR(20) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

-- +goose Down
DROP TABLE IF EXISTS urls;