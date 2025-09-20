package main

import (
	"database/sql"
	"log"

	"github.com/ayushrakesh/gopay/api"
	db "github.com/ayushrakesh/gopay/db/sqlc"
	"github.com/ayushrakesh/gopay/util"
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

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server!", err)
	}
}
