-- 1. ID
-- 2. task.ID 外键
-- 3. 执行类型（begin表示初次, retry表示重试）
-- 4. 执行成功、失败的状态（'success', 'fail'）
-- 5. 错误信息（success对应 null, fail 对应一个字符串）
-- 6. 任务开始的时间
-- 7. 任务结束的时间

CREATE TABLE task_runs (
       id bigserial PRIMARY KEY,
       task_id bigint NOT NULL,
       -- 暂时不实现
       --run_type TEXT NOT NULL DEFAULT 'begin' CHECK (run_type IN ('begin', 'retry')),
       run_status TEXT NOT NULL CHECK (run_status IN ('success', 'fail')),
       run_message TEXT DEFAULT NULL,
       begin_time TIMESTAMPTZ NOT NULL DEFAULT (now()),
       end_time TIMESTAMPTZ NOT NULL,
       CONSTRAINT runs_and_id
                  FOREIGN KEY (task_id)
                  REFERENCES tasks(id)
);
