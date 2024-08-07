package main

import (
	"database/sql"
	"log"

	"github.com/ayushrakesh/go-bank/api"
	db "github.com/ayushrakesh/go-bank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	testDB, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting database")
	}

	store := db.NewStore(testDB)

	server := api.NewServer(store)

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server!", err)
	}
}
