-- +goose Up
CREATE TABLE user_sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       int         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    session_token TEXT        NOT NULL UNIQUE,
    user_agent    TEXT,
    ip_address    TEXT,
    created_at    TIMESTAMPTZ DEFAULT now(),
    expires_at    TIMESTAMPTZ NOT NULL,
    revoked       BOOLEAN     DEFAULT FALSE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions (user_id);

-- +goose Down
DROP TABLE IF EXISTS user_sessions;