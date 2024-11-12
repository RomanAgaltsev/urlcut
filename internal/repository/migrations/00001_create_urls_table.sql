-- +goose Up
CREATE TABLE
    urls (
        id SERIAL PRIMARY KEY,
        long_url VARCHAR(8000) NOT NULL,
        base_url VARCHAR(100) NOT NULL,
        url_id VARCHAR(20) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

-- +goose Down
DROP TABLE urls;