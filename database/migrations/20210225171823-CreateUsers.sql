
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(128) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL ,
    created_at DATETIME,
    updated_at DATETIME,
    UNIQUE KEY (email)
);
-- +migrate Down
DROP TABLE IF EXISTS users;