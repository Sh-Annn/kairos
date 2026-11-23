-- name: CreateTaskRun :one

-- INSERT INTO task_runs (
--        task_id, run_type, run_status, run_message, begin_time, end_time
-- ) VALUES (
--        $1, $2, $3, $4, $5, $6
-- ) RETURNING *;
INSERT INTO task_runs (
       task_id, run_status, run_message, begin_time, end_time
) VALUES (
       $1, $2, $3, $4, $5
) RETURNING *;
