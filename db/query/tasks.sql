-- name: CreateTask :one
INSERT INTO tasks (
       task_type, handler_config, schedule_type, cron_expr, interval_seconds,
       timezone, max_retry, retry_backoff, timeout_seconds, next_run_at
)
VALUES (
       $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET
    status = COALESCE(sqlc.narg(status), status),
    retry_count = COALESCE(sqlc.narg(retry_count), retry_count),
    updated_at = COALESCE(sqlc.narg(updated_at), updated_at),
    next_run_at = COALESCE(sqlc.narg(next_run_at), next_run_at)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: ListDueTasks :many
SELECT * FROM tasks
WHERE status = 'active' AND next_run_at <= now()
ORDER BY next_run_at
LIMIT $1;
