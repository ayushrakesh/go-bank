package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ayushrakesh/go-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	var err error
	config, errr := util.LoadConfig("../..")
	if errr != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Error connecting database")
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
