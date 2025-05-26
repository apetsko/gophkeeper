-- +goose Up
CREATE TABLE user_data (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           minio_object_id TEXT,
                           data_name TEXT NOT NULL,
                           data_type TEXT NOT NULL,
                           meta JSONB DEFAULT '{}'::JSONB,
                           created_at TIMESTAMPTZ DEFAULT now(),
                           updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_user_data_user_id ON user_data(user_id);
CREATE INDEX idx_user_data_data_type ON user_data(data_type);

CREATE TRIGGER trg_user_data_updated_at
    BEFORE UPDATE ON user_data
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS trg_user_data_updated_at ON user_data;
DROP TABLE IF EXISTS user_data;