package scheduler

import (
	"context"
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

	//	CreateNewTask(ctx, s.queries)

	for {
		s.Check(ctx)
		// select {
		// case <-ctx.Done():
		// 	return nil
		// case <-ticker.C:
		// 	continue
		// }
		<-ticker.C
	}
}

// List Tasks from database, and run them
func (s *Scheduler) Check(ctx context.Context) {
	queryTasks, err := s.queries.ListDueTasks(ctx, 10)
	if err != nil {
		log.Println(fmt.Errorf("failed to query tasks, %w", err))
	}

	if len(queryTasks) <= 0 {
		log.Println("no task need run")
		return
	}

	log.Println("*************************************************")
	for _, task := range queryTasks {
		log.Printf("task_id: %d, next_run_at: %v\n", task.ID, task.NextRunAt)

		err := s.execTask(ctx, &task)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("*************************************************")
}

// execute each task
func (s *Scheduler) execTask(ctx context.Context, task *db.Task) error {
	// bad Config
	var handlerConfig domain.HandlerConfig
	err := json.Unmarshal(task.HandlerConfig, &handlerConfig)
	if err != nil {
		return s.markInvalid(ctx, task)
	}

	// Call HTTP
	// execCtx, execCannel := context.WithTimeout(ctx, time.Duration(task.IntervalSeconds.Int32)*time.Second)
	body, execErr := s.executor.Exec(ctx, handlerConfig)
	// execCannel()

	if execErr != nil {
		return s.HandleFailure(ctx, task, domain.HTTPError(execErr.Error()))
	}
	log.Println(string(body))
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

		ScheduleType:   "once",
		Timezone:       "Asia/Shanghai",
		MaxRetry:       3,
		RetryBackoff:   30,
		TimeoutSeconds: 30,
		NextRunAt:      time.Now().Add(time.Duration(2 * time.Second)),
	})
	if err != nil {
		log.Println(fmt.Errorf("failed to create task: %w", err))
		return
	}
	log.Printf("task: %+v\n", task)
}
