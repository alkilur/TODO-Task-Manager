package sqlite

import (
	"database/sql"
	"strconv"

	srv "yet-another-todo-list/internal/http-server"

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
			date VARCHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
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

func (s *Storage) CreateTask(task *srv.Task) (string, error) {
	stmt, err := s.db.Prepare(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return "", err
	}

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}
