package main

import (
	"context"
	"database/sql"
	"time"

	"log"

	"github.com/Sh-Annn/kairos/db/sqlc"
	"github.com/Sh-Annn/kairos/internal/executor"
	"github.com/Sh-Annn/kairos/internal/scheduler"

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

	e := executor.NewHTTPExecutor()

	s := scheduler.NewScheduler(db.New(conn), e, 2*time.Second)

	if err := s.Run(context.Background()); err != nil {
		log.Println(err)
	}
}
