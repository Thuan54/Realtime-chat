-- +goose Up
-- Seed Users
INSERT INTO users (id, email, password_hashed) VALUES
(1, 'alice@example.com', '$2a$12$LJ3m9sW9vJ1K2H3G4F5D6eA7b8c9d0e1f2g3h4i5j6k7l8m9n0o'),
(2, 'bob@example.com',   '$2a$12$LJ3m9sW9vJ1K2H3G4F5D6eA7b8c9d0e1f2g3h4i5j6k7l8m9n0o'),
(3, 'charlie@example.com', '$2a$12$LJ3m9sW9vJ1K2H3G4F5D6eA7b8c9d0e1f2g3h4i5j6k7l8m9n0o');

-- Seed Channels
INSERT INTO channels (id, name) VALUES
(1, 'general'),
(2, 'random'),
(3, 'dev-updates');

-- Seed Memberships
INSERT INTO channel_joined (channel_id, user_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 1),
(2, 3),
(3, 2);

-- Seed Messages
INSERT INTO messages (id, sender_id, content, created_at) VALUES
(1, 1, 'Welcome to the chat!', NOW() - INTERVAL '29 days'),
(2, 2, 'Thanks, Alice! Excited to be here.', NOW() - INTERVAL '28 days'),
(3, 3, 'Hello everyone!', NOW() - INTERVAL '20 days'),
(4, 1, 'Check out the new docs.', NOW() - INTERVAL '19 days'),
(5, 2, 'Will do. Also, server is looking stable.', NOW() - INTERVAL '18 days'),
(6, 3, 'Just pushed a fix for the WS reconnect logic.', NOW() - INTERVAL '12 days');

-- Seed Channels-Messages Junction
INSERT INTO channel_messages (channel_id, message_id) VALUES
(1, 1), (1, 2), (1, 3),
(2, 4), (2, 6),
(3, 5);

-- Reset identity sequences to prevent collisions on future inserts
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('channels_id_seq', (SELECT MAX(id) FROM channels));
SELECT setval('messages_id_seq', (SELECT MAX(id) FROM messages));

-- +goose Down
DELETE FROM channel_messages;
DELETE FROM messages;
DELETE FROM channel_joined;
DELETE FROM channels;
DELETE FROM users;