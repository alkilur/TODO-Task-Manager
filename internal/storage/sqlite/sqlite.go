package sqlite

import (
	"database/sql"
	"strconv"
	"time"

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

func (s *Storage) GetTasks(searchQuery string) ([]srv.Task, error) {
	var rows *sql.Rows

	date, err := time.Parse("02.01.2006", searchQuery)
	if err == nil {
		rows, err = s.db.Query(`
			SELECT * FROM scheduler
			WHERE date = ? LIMIT 50`,
			date.Format(srv.TimeLayout))
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = s.db.Query(`
			SELECT * FROM scheduler
			WHERE title LIKE $searchQuery OR comment LIKE $searchQuery
			ORDER BY date LIMIT 50`,
			sql.Named("searchQuery", "%"+searchQuery+"%"))
		if err != nil {
			return nil, err
		}
	}

	tasks := make([]srv.Task, 0)
	for rows.Next() {
		task := srv.Task{}
		if err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(task *srv.Task) error {
	stmt, err := s.db.Prepare(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if affectedRows == 0 || err != nil {
		return srv.ErrInvalidID
	}

	return nil
}

func (s *Storage) CompleteTask(id string) error {

	task := srv.Task{}
	row := s.db.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return err
	}

	if task.Repeat == "" {
		res, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			return err
		}
		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			return srv.ErrInvalidID
		}
	} else {
		nextDate, err := srv.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		res, err := s.db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
		if err != nil {
			return err
		}
		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			return srv.ErrInvalidID
		}
	}
	return nil
}
