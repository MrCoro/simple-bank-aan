package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:rahasia@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB


func TestMain(m *testing.M){
	var err error

	testDB, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot Connect to DB:", err)
	} else {
		fmt.Println("Connected to database")
	}
 
	testQueries = New(testDB)

	os.Exit(m.Run())
}