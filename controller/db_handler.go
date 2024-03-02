package controller

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()
var db_port = os.Getenv("DB_PORT")
var dataSource = "root:@tcp(localhost:" + db_port + ")/db_kuis"

func connect() *sql.DB {

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
