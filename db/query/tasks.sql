-- name: CreateTask :one
INSERT INTO tasks (
       task_type, handler_config, schedule_type, cron_expr, interval_seconds,
       timezone, max_retry, retry_backoff, timeout_seconds, next_run_at
)
VALUES (
       $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateTaskNextRun :one
UPDATE tasks
SET
    next_run_at = $2
WHERE
    id = $1
RETURNING *;

-- name: ListDueTasks :many
SELECT * FROM tasks
WHERE status = 'active' AND next_run_at <= now()
ORDER BY next_run_at
LIMIT $1;
