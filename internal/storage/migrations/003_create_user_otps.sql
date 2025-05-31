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

-- CREATE TRIGGER trg_user_otps_updated_at
--     BEFORE UPDATE ON user_otps
--     FOR EACH ROW
-- EXECUTE PROCEDURE update_updated_at_column();


-- +goose Down
-- DROP TRIGGER IF EXISTS trg_user_otps_updated_at ON user_otps;
DROP TABLE IF EXISTS user_otps;