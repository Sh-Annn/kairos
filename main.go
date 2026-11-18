package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	db "github.com/Sh-Annn/kairos/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	// connect database
	conn, err := sql.Open("postgres", "host=localhost dbname=kairos connect_timeout=5 sslmode=disable user=root password=secret")
	if err != nil {
		log.Fatal(err)
	}
	err = conn.PingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("connect kairos database")

	// conn is *sql.DB, func New will return a *db.Queries for sql query
	q := db.New(conn)

	// cycle of search database and run task is 2s
	ticker := time.NewTicker(2 * time.Second)
	CreateNewTask(q)
	for {
		//CreateNewTask(q)
		CheckTask(q)

		// Sleep: time = max(triker, all database time)
		<-ticker.C
	}
}

func CreateNewTask(q *db.Queries) {
	// ctx is for CreateTask
	ctx, Cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer Cannel()

	handlerConfig := HandlerConfig{
		URL:    "",
		Method: "GET",
	}
	data, err := json.Marshal(handlerConfig)
	if err != nil {
		panic(err)
	}
	task, err := q.CreateTask(ctx, db.CreateTaskParams{
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

type HandlerConfig struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func CheckTask(q *db.Queries) {
	// ctx is for CreateTask
	ctx, cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cannel()
	queryTasks, err := q.ListDueTasks(ctx, 10)
	if err != nil {
		log.Println(fmt.Errorf("failed to query tasks, %w", err))
	}

	if len(queryTasks) > 0 {
		log.Println("*************************************************")
		for _, task := range queryTasks {
			log.Printf("task_id: %d, next_run_at: %v\n", task.ID, task.NextRunAt)
			var handlerConfig HandlerConfig
			err := json.Unmarshal(task.HandlerConfig, &handlerConfig)
			if err != nil {
				log.Println(fmt.Errorf("cannot unmarshal json URL and Method: %w\n", err))
				continue
			}

			// TODO: valid HandlerConfig
			if handlerConfig.URL == "" {
				_, err := q.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
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
			execTask(ctx, q, &task, handlerConfig.URL, handlerConfig.Method)
		}
		log.Println("*************************************************")
	} else {
		log.Println("no task need run")
	}
}

// get https://httpbin.org for now
func execTask(ctx context.Context, q *db.Queries, task *db.Task, url string, method string) {
	if method != http.MethodGet {
		return
	}

	err := execHTTP(ctx, url, time.Duration(task.TimeoutSeconds)*time.Second)
	if err != nil {
		log.Println(fmt.Errorf("cannot exec http: %w", err))
		return
	}

	updateTask, err := q.UpdateTaskNextRun(ctx, db.UpdateTaskNextRunParams{
		ID:        task.ID,
		NextRunAt: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		log.Println(fmt.Errorf("cannot change next_run_at %w", err))
		return
	}
	log.Printf("success change next_run_at, task_id: %d, next_run_at: %v", updateTask.ID, updateTask.NextRunAt)
}

func execHTTP(ctx context.Context, url string, timeout time.Duration) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(fmt.Errorf("cannot create request: %w", err))
		return err
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad http statusCode")
	}
	log.Println(resp.Header)

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(body))
	return nil
}
