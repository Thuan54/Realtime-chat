-- +goose Up
CREATE TABLE channel_messages (
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    PRIMARY KEY (channel_id, message_id)
);

-- +goose Down
DROP TABLE IF EXISTS channel_messages;