package database

import (
	"database/sql"
	"log"
	"os"
)

func Connect() *sql.DB {
	connect_string := os.Getenv("DB_CONN_URL")
	db, err := sql.Open("postgres", connect_string)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
