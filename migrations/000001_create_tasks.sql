-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS task_tracker;

CREATE TYPE task_tracker.task_status AS ENUM ('todo', 'in-progress', 'done', 'cancelled');

CREATE TABLE task_tracker.tasks
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description TEXT NOT NULL,
    status task_tracker.task_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL  DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_tracker.tasks;
DROP TYPE IF EXISTS task_tracker.task_status;
DROP SCHEMA IF EXISTS task_tracker;
-- +goose StatementEnd