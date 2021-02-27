
-- +migrate Up
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(128) PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    content TEXT,
    is_completed BOOLEAN DEFAULT false,
    deadline DATE,
    user_id VARCHAR(128) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX index_tasks_on_user_id_and_deadline ON tasks (user_id, deadline);
-- +migrate Down
DROP TABLE IF EXISTS tasks;