-- +goose Up
-- Optimizes cursor-based pagination (ORDER BY created_at DESC)
CREATE INDEX idx_message_history ON messages(created_at DESC);

-- Optimizes JOINs for /channels endpoint (fetch channels by user_id)
CREATE INDEX idx_channel_joined_user_id ON channel_joined(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_message_history;
DROP INDEX IF EXISTS idx_channel_joined_user_id;