package internal

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS heartbeats (
		id TEXT PRIMARY KEY,
		last_seen TIMESTAMP
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db
}
