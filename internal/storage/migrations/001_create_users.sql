-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(60)         NOT NULL,
    totp_secret   TEXT,
    totp_enabled  BOOLEAN     DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS users;