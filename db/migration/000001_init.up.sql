CREATE TYPE task_status AS ENUM ('active', 'paused', 'deleted');
CREATE TYPE schedule_type AS ENUM ('once', 'interval', 'cron');

CREATE TABLE tasks (
       id bigserial PRIMARY KEY,
       task_type TEXT NOT NULL DEFAULT 'http' CHECK(task_type IN ('http')) ,
       handler_config JSONB NOT NULL,
       schedule_type schedule_type NOT NULL DEFAULT 'once',
       cron_expr TEXT,
       interval_seconds INT,
       status task_status NOT NULL DEFAULT 'active',
       created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
       timezone TEXT NOT NULL DEFAULT 'UTC',
       max_retry INT NOT NULL DEFAULT 3,
       retry_backoff INT NOT NULL DEFAULT 30,
       timeout_seconds INT NOT NULL DEFAULT 60,
       updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
       next_run_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX next_run ON tasks (next_run_at);
