-- +goose Up
CREATE TABLE channel_joined (
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (channel_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS channel_joined;