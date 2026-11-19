package main

import (
	"context"
	"database/sql"

	"log"
	"time"

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
	ticker := time.NewTicker(2 * time.Second)
	scheduler := scheduler.NewScheduler(conn, *ticker)

	scheduler.Run()
}
