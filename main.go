package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("<repeat> is empty")
	}

	nextDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid <date> format: %w", err)
	}

	repeatParts := strings.Split(repeat, " ")
	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", errors.New("invalid <repeat> format")
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days >= 400 {
			return "", errors.New("invalid <repeat> format")
		}
		nextDate = nextDate.AddDate(0, 0, days)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("invalid <repeat> format")
	}
	return nextDate.Format("20060102"), nil
}

func main() {

	// .env load
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	// DB Init
	if err := InitDB(); err != nil {
		log.Fatal(err)
	}
	if err := TableCreate(); err != nil {
		log.Fatal(err)
	}

	// Set port
	port := os.Getenv("TODO_PORT")
	if port == "" {
		if err := os.Setenv("TODO_PORT", "7540"); err != nil {
			log.Fatal("error when trying to set port: ", err)
		}
		port = os.Getenv("TODO_PORT")
		log.Println("<TODO_PORT> has been redefined")
	}

	// Handlers
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)

	// Run
	log.Printf("Server is starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
