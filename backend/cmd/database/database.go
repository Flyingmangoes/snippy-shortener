package database

import (
	"database/sql"
	"log"
	"log/slog"
	_"github.com/lib/pq"
)

func NewDatabaseConnection(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database: ",err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to establish db connection: ", err)
	}

	msg := "Connected"
	slog.Info("[DEBUG]", "DB Status", msg)
	return db
}