CREATE TYPE task_type AS ENUM ('http');
CREATE TYPE task_status AS ENUM ('active', 'paused', 'deleted', 'completed');
CREATE TYPE schedule_type AS ENUM ('once', 'interval', 'cron');


-- task_type TEXT NOT NULL DEFAULT 'http' CHECK(task_type IN ('http')) ,
-- schedule_type TEXT NOT NULL
-- DEFAULT 'once' CHECK(schdule_type IN ('once', 'interval', 'cron')),
-- status TEXT NOT NULL
-- DEFAULT 'active' CHECK(status IN ('active', 'paused', 'deleted', 'completed')),
CREATE TABLE tasks (
       id bigserial PRIMARY KEY,
       task_type task_type NOT NULL DEFAULT 'http', 
       handler_config JSONB NOT NULL,
       schedule_type schedule_type NOT NULL DEFAULT 'once',
       cron_expr TEXT,
       interval_seconds INT,
       status task_status NOT NULL DEFAULT 'active',
       created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
       timezone TEXT NOT NULL DEFAULT 'UTC',
       retry_count INT NOT NULL DEFAULT 0,
       max_retry INT NOT NULL DEFAULT 3,
       retry_backoff INT NOT NULL DEFAULT 30,
       timeout_seconds INT NOT NULL DEFAULT 60,
       updated_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
       next_run_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX next_run ON tasks (next_run_at);
