-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT                NOT NULL,
    created_at    TIMESTAMPTZ      DEFAULT now(),
    updated_at    TIMESTAMPTZ      DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS users;