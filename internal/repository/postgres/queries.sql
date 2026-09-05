-- name: GetTask :one
SELECT id, description, status, created_at, updated_at
FROM task_tracker.tasks
WHERE id = $1;

-- name: GetTasks :many
SELECT id, description, status, created_at, updated_at
FROM task_tracker.tasks
ORDER BY id;

-- name: DeleteTask :exec
DELETE
FROM task_tracker.tasks
WHERE id = $1;

-- name: CreateTask :one
INSERT INTO task_tracker.tasks (description, status)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateTask :exec
UPDATE task_tracker.tasks
SET description = $2,
    status      = $3,
    updated_at  = NOW()
WHERE id = $1;