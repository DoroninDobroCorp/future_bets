-- +migrate Up
CREATE TABLE IF NOT EXISTS tg_testbot.bot_files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_date ON tg_testbot.bot_files (date);
CREATE INDEX IF NOT EXISTS idx_filename ON tg_testbot.bot_files (filename);

CREATE TABLE IF NOT EXISTS tg_testbot.bot_subscribers (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    chat_id BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_username ON tg_testbot.bot_subscribers (username);
CREATE INDEX IF NOT EXISTS idx_chat_id ON tg_testbot.bot_subscribers (chat_id);

-- +migrate Down
DROP TABLE IF EXISTS bot_files;
DROP TABLE IF EXISTS bot_subscribers;