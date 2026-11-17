package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "host=localhost dbname=kairos connect_timeout=5")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("connect kairos database")
}
