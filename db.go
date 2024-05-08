package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() error {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		if err := os.Setenv("TODO_DBFILE", "./scheduler.db"); err != nil {
			return fmt.Errorf("error when trying to initialize db: %v", err)
		}
		log.Println("<TODO_DBFILE> has been redefined")
	}

	if _, err := os.Stat(dbPath); err != nil {
		if _, err := os.Create(dbPath); err != nil {
			return fmt.Errorf("error when trying to initialize db: %v", err)
		}
		log.Println("database has been created")
	}
	return nil
}

func TableCreate() error {
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		return fmt.Errorf("error when trying to open db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
	);`)
	if err != nil {
		return fmt.Errorf("error when trying to table create: %v", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);`)
	if err != nil {
		return fmt.Errorf("error when trying to table create: %v", err)
	}
	log.Println("table has been created")
	return nil
}
