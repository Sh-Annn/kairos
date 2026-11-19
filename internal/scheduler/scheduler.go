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

func (s *Scheduler) Run(ctx context.Context) error {
	s.CreateNewTask() // for check later operate, always migratedown before go run

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.Check(ctx)
		}
	}
}

func (s *Scheduler) Check(ctx context.Context) {
	// ctx, cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	// defer cannel()
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

func (s *Scheduler) execTask(ctx context.Context, task *db.Task) error {
	// bad Config
	var handlerConfig domain.HandlerConfig
	err := json.Unmarshal(task.HandlerConfig, &handlerConfig)
	if err != nil {
		return s.markInvalid(ctx, task)
	}

	// Call HTTP
	execCtx, execCannel := context.WithTimeout(ctx, time.Duration(task.IntervalSeconds.Int32)*time.Second)
	body, execErr := s.executor.Exec(execCtx, handlerConfig)
	execCannel()

	if execErr != nil {
		return s.HandleFailure(ctx, task, execErr)
	}
	log.Println(body)
	return s.HandleSuccess(ctx, task)

	// TODO: abstract a function
	// updateTask, err := s.queries.UpdateTaskNextRun(ctx, db.UpdateTaskNextRunParams{
	// 	ID:        task.ID,
	// 	NextRunAt: time.Now().Add(time.Duration(task.TimeoutSeconds) * time.Second),
	// })
	// if err != nil {
	// 	return fmt.Errorf("cannot change next_run_at %w", err)
	// }
	// log.Printf("success change next_run_at, task_id: %d, next_run_at: %v", updateTask.ID, updateTask.NextRunAt)
	// return nil
}

func (s *Scheduler) updateTask(ctx context.Context, task *db.Task) error {
	var err error
	switch task.ScheduleType {
	case db.ScheduleTypeOnce:
		_, err = s.queries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
			ID:     task.ID,
			Status: "deleted",
		})
		if err != nil {
			return fmt.Errorf("cannot change task status: %w", err)
		}
		// TODO: diffrent SchedultType
		// case db.ScheduleTypeInterval:
		// case db.ScheduleTypeCron:
	}

	return nil
}

func (s *Scheduler) CreateNewTask() {
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
	task, err := s.queries.CreateTask(ctx, db.CreateTaskParams{
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
