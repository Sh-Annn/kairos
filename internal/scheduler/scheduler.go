package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/Sh-Annn/kairos/db/sqlc"
	"github.com/Sh-Annn/kairos/internal/executor"
	"github.com/Sh-Annn/kairos/util"
)

type Scheduler struct {
	Queries  *db.Queries
	Ticker   time.Ticker
	Executor executor.Executor
}

func NewScheduler(database *sql.DB, ticker time.Ticker) *Scheduler {
	return &Scheduler{
		Queries: db.New(database),
		Ticker:  ticker,
	}
}

func (s *Scheduler) Run() {
	s.CreateNewTask() // for check later operate, always migratedown before go run
	for {
		s.Check()

		<-s.Ticker.C
	}
}

func (s *Scheduler) Check() {
	ctx, cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cannel()
	queryTasks, err := s.Queries.ListDueTasks(ctx, 10)
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
		var handlerConfig util.HandlerConfig
		err := json.Unmarshal(task.HandlerConfig, &handlerConfig)
		if err != nil {
			log.Println(fmt.Errorf("cannot unmarshal json URL and Method: %w\n", err))
			continue
		}

		// TODO: valid HandlerConfig
		if handlerConfig.URL == "" {
			_, err := s.Queries.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
				ID:     task.ID,
				Status: "deleted",
			})
			if err != nil {
				log.Println(fmt.Errorf("cannot change task status: %w", err))
				continue
			}
			log.Println("task's url is not avilable, changed task staus")
			continue
		}
		log.Println(handlerConfig.Method, handlerConfig.URL)

		ctx, Cannel := context.WithTimeout(context.Background(), time.Duration(task.TimeoutSeconds)*time.Second)
		defer Cannel()

		s.execTask(ctx, &task, handlerConfig)
	}
	log.Println("*************************************************")
}

func (s *Scheduler) execTask(ctx context.Context, task *db.Task, handlerConfig util.HandlerConfig) {
	if handlerConfig.Method != http.MethodGet {
		return
	}

	s.Executor = executor.NewHttpExecutor(handlerConfig, time.Duration(task.TimeoutSeconds))
	err := s.Executor.Exec(ctx)
	if err != nil {
		log.Println(fmt.Errorf("cannot exec http: %w", err))
		return
	}

	updateTask, err := s.Queries.UpdateTaskNextRun(ctx, db.UpdateTaskNextRunParams{
		ID:        task.ID,
		NextRunAt: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		log.Println(fmt.Errorf("cannot change next_run_at %w", err))
		return
	}
	log.Printf("success change next_run_at, task_id: %d, next_run_at: %v", updateTask.ID, updateTask.NextRunAt)
}

func (s *Scheduler) CreateNewTask() {
	// ctx is for CreateTask
	ctx, Cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer Cannel()

	handlerConfig := util.HandlerConfig{
		URL:    "",
		Method: "GET",
	}
	data, err := json.Marshal(handlerConfig)
	if err != nil {
		panic(err)
	}
	task, err := s.Queries.CreateTask(ctx, db.CreateTaskParams{
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
