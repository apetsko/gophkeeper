-- +goose Up
CREATE TABLE user_otps
(
    id         SERIAL PRIMARY KEY,
    user_id    INT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    otp_code   TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    used       BOOLEAN     DEFAULT FALSE,
    CONSTRAINT otp_unique_code UNIQUE (user_id, otp_code)
);

CREATE INDEX idx_user_otps_user_id ON user_otps (user_id);

-- +goose Down
DROP TABLE IF EXISTS user_otps;