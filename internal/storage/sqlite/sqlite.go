package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(DBPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			comment TEXT,
			date VARCHAR(8) NOT NULL,
			repeat VARCHAR(128)
		);
		CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
	`)
	if err != nil {
		return nil, err
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}
