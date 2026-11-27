package main

import (
	"context"
	"database/sql"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/Sh-Annn/kairos/db/sqlc"
	"github.com/Sh-Annn/kairos/internal/executor"
	"github.com/Sh-Annn/kairos/internal/scheduler"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// connect database
	conn, err := sql.Open("postgres", "host=localhost dbname=kairos connect_timeout=5 sslmode=disable user=root password=secret")
	if err != nil {
		log.Fatal(err)
	}

	// test conn
	err = conn.PingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connect kairos database")
	defer conn.Close()

	e := executor.NewHTTPExecutor()

	s := scheduler.NewScheduler(db.New(conn), e, 2*time.Second)

	if err := s.Run(ctx); err != nil {
		log.Println(err)
	}

	log.Println("exited program")
}
