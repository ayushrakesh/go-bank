package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting database")
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
