package main

import (
	"database/sql"
	"log"

	"github.com/ayushrakesh/go-bank/api"
	db "github.com/ayushrakesh/go-bank/db/sqlc"
	"github.com/ayushrakesh/go-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("couldnt load config", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error connecting database")
	}

	store := db.NewStore(testDB)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server!", err)
	}
}
