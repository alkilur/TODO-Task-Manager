package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// func responseWithError(w http.ResponseWriter, err error) {
// 	log.Printf("response with error: %v", err)
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
// 		http.Error(w, "Error when trying to send error response", http.StatusInternalServerError)
// 		log.Printf("error when trying to send error response: %v", encodeErr)
// 	}
// }

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		http.Error(w, "Invalid <now> format", http.StatusBadRequest)
		log.Printf("invalid <now> format: %v", err)
		return
	}
	nextDate, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		log.Printf("NextDate func error: %v", err)
		return
	}
	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Printf("error when trying to send response: %v", err)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error":"id field is empty"}`, http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to get task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to open db: %v", err)
			return
		}
		defer db.Close()

		task := Task{}
		row := db.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
		if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error":"Error when trying to get task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to row.Scan(): %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
			log.Printf("error when trying to send response: %v", err)
		}

	case http.MethodPost:
		task := Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Title is empty"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			if _, err := time.Parse("20060102", task.Date); err != nil {
				http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
				return
			}
		}
		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Invalid repeat format"}`, http.StatusBadRequest)
				return
			}
			if task.Date < time.Now().Format("20060102") {
				task.Date = nextDate
			}
		}

		db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to create task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to create task: %v", err)
			return
		}
		defer db.Close()

		res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
			task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to create task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to create task: %v", err)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Error when trying to get task ID"}`, http.StatusInternalServerError)
			log.Printf("error when trying to get task ID: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"id": strconv.FormatInt(id, 10)}); err != nil {
			http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
			log.Printf("error when trying to send response: %v", err)
		}

	case http.MethodPut:
		task := Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			http.Error(w, `{"error":"ID is empty"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Title is empty"}`, http.StatusBadRequest)
			return
		}

		if task.Date != "" {
			if _, err := time.Parse("20060102", task.Date); err != nil {
				http.Error(w, `{"error":"Invalid date format"}`, http.StatusBadRequest)
				return
			}
		}
		if task.Date == "" || task.Date < time.Now().Format("20060102") {
			task.Date = time.Now().Format("20060102")
		}

		if task.Repeat != "" {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Invalid repeat format"}`, http.StatusBadRequest)
				return
			}
			if task.Date < time.Now().Format("20060102") {
				task.Date = nextDate
			}
		}

		db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to update task"}`, http.StatusBadRequest)
			log.Printf("error when trying to open db: %v", err)
			return
		}
		defer db.Close()

		res, err := db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
			task.Date, task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to update task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to update task: %v", err)
			return
		}

		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			http.Error(w, `{"error":"Invalid task id"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
			log.Printf("error when trying to send response: %v", err)
		}

	case http.MethodDelete:
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error":"id field is empty"}`, http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to delete task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to open db: %v", err)
			return
		}
		defer db.Close()

		res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to delete task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to delete task from db: %v", err)
			return
		}
		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			http.Error(w, `{"error":"Invalid task id"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
			log.Printf("error when trying to send response: %v", err)
		}

	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Invalid HTTP method"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, `{"error":"Error when trying to get tasks"}`, http.StatusInternalServerError)
		log.Printf("error when trying to open db: %v", err)
		return
	}
	defer db.Close()

	var rows *sql.Rows

	searchValue := r.FormValue("search")
	date, err := time.Parse("02.01.2006", searchValue)
	if err == nil {
		rows, err = db.Query(`SELECT * FROM scheduler WHERE date = ? LIMIT 50`, date.Format("20060102"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to get tasks"}`, http.StatusInternalServerError)
			log.Printf("error when trying to select tasks: %v", err)
			return
		}
	} else {
		rows, err = db.Query(`SELECT * FROM scheduler WHERE title LIKE $searchValue OR comment LIKE $searchValue ORDER BY date LIMIT 50`,
			sql.Named("searchValue", "%"+searchValue+"%"))
		if err != nil {
			http.Error(w, `{"error":"Error when trying to get tasks"}`, http.StatusInternalServerError)
			log.Printf("error when trying to select tasks: %v", err)
			return
		}
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		task := Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error":"Error when trying to get tasks"}`, http.StatusInternalServerError)
			log.Printf("error when trying to rows.Scan(): %v", err)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, `{"error":"Error when trying to get tasks"}`, http.StatusInternalServerError)
		log.Printf("error in rows.Err(): %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(map[string][]Task{"tasks": tasks}); err != nil {
		http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
		log.Printf("error when trying to send response: %v", err)
	}
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Invalid HTTP method"}`, http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, `{"error":"id field is empty"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		http.Error(w, `{"error":"Error when trying to done task"}`, http.StatusInternalServerError)
		log.Printf("error when trying to open db: %v", err)
		return
	}
	defer db.Close()

	task := Task{}
	row := db.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		http.Error(w, `{"error":"Error when trying to done task"}`, http.StatusInternalServerError)
		log.Printf("error when trying to row.Scan(): %v", err)
		return
	}

	if task.Repeat == "" {
		res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to done task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to delete task from db: %v", err)
			return
		}
		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			http.Error(w, `{"error":"Invalid task id"}`, http.StatusBadRequest)
			return
		}
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to date calculate"}`, http.StatusInternalServerError)
			log.Printf("error when trying to calculate nextDate: %v", err)
			return
		}
		res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
		if err != nil {
			http.Error(w, `{"error":"Error when trying to update task"}`, http.StatusInternalServerError)
			log.Printf("error when trying to update task: %v", err)
			return
		}
		affectedRows, err := res.RowsAffected()
		if affectedRows == 0 || err != nil {
			http.Error(w, `{"error":"Invalid task id"}`, http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		http.Error(w, `{"error":"Error when trying to send response"}`, http.StatusInternalServerError)
		log.Printf("error when trying to send response: %v", err)
	}
}
