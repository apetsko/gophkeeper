-- +goose Up
CREATE TABLE user_keys
(
    id                   SERIAL PRIMARY KEY,
    user_id              int   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    encrypted_master_key BYTEA NOT NULL,
    nonce                BYTEA NOT NULL,
    created_at           TIMESTAMPTZ DEFAULT now(),
    updated_at           TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_user_keys_user_id ON user_keys (user_id);

-- +goose Down
DROP TABLE IF EXISTS user_keys;