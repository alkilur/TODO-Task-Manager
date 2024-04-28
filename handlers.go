package main

import (
	"log"
	"net/http"
	"time"
)

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
		log.Printf("error writing response: %v", err)
	}
}
