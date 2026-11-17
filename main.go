package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	ticker := time.NewTicker(time.Duration(2000 * time.Millisecond))
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
	task, err := q.CreateTask(ctx, db.CreateTaskParams{
		TaskType:       "http",
		HandlerConfig:  json.RawMessage([]byte(`{"message": "hello kairos"}`)),
		ScheduleType:   "once",
		Timezone:       "Asia/Shanghai",
		MaxRetry:       3,
		RetryBackoff:   30,
		TimeoutSeconds: 30,
		NextRunAt:      time.Now().Add(time.Duration(2 * time.Minute)),
	})
	if err != nil {
		log.Println(fmt.Errorf("failed to create task: %w", err))
		return
	}
	log.Printf("task: %v\n", task)
}

func CheckTask(q *db.Queries) {
	// ctx is for CreateTask
	ctx, Cannel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer Cannel()
	query_tasks, err := q.ListDueTasks(ctx, 10)
	if err != nil {
		log.Println(fmt.Errorf("failed to query tasks, %w", err))
	}

	if len(query_tasks) > 0 {
		log.Println("*************************************************")
		for _, task := range query_tasks {
			log.Printf("task_id: %d, next_run_at: %v\n", task.ID, task.NextRunAt)
		}
		log.Println("*************************************************")
	} else {
		log.Println("no task need run")
	}
}
