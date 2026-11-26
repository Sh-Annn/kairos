package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	db "github.com/Sh-Annn/kairos/db/sqlc"
	"github.com/Sh-Annn/kairos/internal/domain"
	"github.com/Sh-Annn/kairos/internal/executor"
)

type Scheduler struct {
	queries  db.Querier
	executor executor.Executor
	interval time.Duration
}

func NewScheduler(database db.Querier, executor executor.Executor, interval time.Duration) *Scheduler {
	return &Scheduler{
		queries:  database,
		executor: executor,
		interval: interval,
	}
}

// Main Loop
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	CreateNewTask(ctx, s.queries)

	var err error
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err = s.Check(ctx)
			if err != nil {
				return err
			}
		}
	}
}

// List Tasks from database, and run them
func (s *Scheduler) Check(ctx context.Context) error {
	queryTasks, err := s.queries.ListDueTasks(ctx, 10)
	if err != nil {
		log.Println(fmt.Errorf("failed to query tasks, %w", err))
	}

	if len(queryTasks) <= 0 {
		log.Println("no tasks needs run for now, continue check...")
		return nil
	}

	// Who manage the number of retry? Praboly the function execTask,
	// all function call returns them error, then execTask actually
	// clean up.
	log.Println("*************************************************")
	for _, task := range queryTasks {
		log.Printf("task_id: %d, next_run_at: %v\n", task.ID, task.NextRunAt)

		err := s.execTask(ctx, &task)
		if err != nil {
			log.Println(err)
		}

	}
	log.Println("*************************************************")
	return nil
}

// execute each task
func (s *Scheduler) execTask(ctx context.Context, task *db.Task) error {
	taskRunsParams := db.CreateTaskRunParams{
		TaskID:    task.ID,
		RunStatus: "success",
		RunMessage: sql.NullString{
			Valid: false,
		},
		BeginTime: time.Now(),
	}

	// Call HTTP
	execCtx, execCannel := context.WithTimeout(ctx, time.Duration(task.TimeoutSeconds)*time.Second)
	body, execErr := s.executor.Exec(execCtx, task.HandlerConfig)
	execCannel()

	if execErr != nil {
		taskRunsParams.RunStatus = "fail"
		taskRunsParams.RunMessage = sql.NullString{
			String: execErr.Error(),
			Valid:  true,
		}
		taskRunsParams.EndTime = time.Now()
		taskRunsParams.Attempt = task.RetryCount + 1

		_, _ = s.queries.CreateTaskRun(ctx, taskRunsParams)
		return s.HandleFailure(ctx, task, execErr)
	}

	log.Println(string(body))
	taskRunsParams.EndTime = time.Now()
	_, _ = s.queries.CreateTaskRun(ctx, taskRunsParams)
	return s.HandleSuccess(ctx, task)
}

func CreateNewTask(ctx context.Context, queries db.Querier) {
	// ctx is for CreateTask
	ctx, Cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer Cannel()

	handlerConfig := domain.HandlerConfig{
		URL:    "https://httpbin.org/get",
		Method: "GET",
	}
	data, err := json.Marshal(handlerConfig)
	if err != nil {
		panic(err)
	}
	task, err := queries.CreateTask(ctx, db.CreateTaskParams{
		TaskType:      "http",
		HandlerConfig: json.RawMessage(data),

		ScheduleType: "once",
		IntervalSeconds: sql.NullInt32{
			Int32: domain.RandomInt(30, 60),
			Valid: true,
		},
		Timezone:       "Asia/Shanghai",
		MaxRetry:       domain.RandomInt(1, 3),
		RetryBackoff:   domain.RandomInt(2, 5),
		TimeoutSeconds: domain.RandomInt(10, 30),
		NextRunAt:      time.Now().Add(time.Duration(10 * time.Second)),
	})
	if err != nil {
		log.Println(fmt.Errorf("failed to create task: %w", err))
		return
	}
	log.Printf("task: %+v\n", task)
}

func (s *Scheduler) CreateTaskRun(ctx context.Context, params db.CreateTaskRunParams) error {
	_, err := s.queries.CreateTaskRun(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
