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
			if task.Repeat == "d 1" {
				task.Date = time.Now().Format("20060102")
			} else {
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Invalid repeat format"}`, http.StatusBadRequest)
					return
				}
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

	default:
		http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
	}
}
