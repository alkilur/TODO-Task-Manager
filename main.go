package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// DB Create
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "./scheduler.db"

		log.Println("<dbPath> has been redefined")
	}

	if _, err := os.Stat(dbPath); err != nil {
		if _, err := os.Create(dbPath); err != nil {
			log.Fatal(err)
		}
		log.Println("database has been created")
	}

	// Table create
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8) NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128)
	);
	CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);`)
	if err != nil {
		log.Fatal(err)
	}

	// Main handler
	http.Handle("/", http.FileServer(http.Dir("./web")))
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
		log.Println("<port> has been redefined")
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
