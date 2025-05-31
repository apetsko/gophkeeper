-- +goose Up
CREATE TABLE user_sessions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               session_token TEXT NOT NULL UNIQUE,
                               user_agent TEXT,
                               ip_address TEXT,
                               created_at TIMESTAMPTZ DEFAULT now(),
                               expires_at TIMESTAMPTZ NOT NULL,
                               revoked BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);

-- CREATE TRIGGER trg_user_sessions_updated_at
--     BEFORE UPDATE ON user_sessions
--     FOR EACH ROW
-- EXECUTE PROCEDURE update_updated_at_column();



-- +goose Down
-- DROP TRIGGER IF EXISTS trg_user_sessions_updated_at ON user_sessions;
DROP TABLE IF EXISTS user_sessions;