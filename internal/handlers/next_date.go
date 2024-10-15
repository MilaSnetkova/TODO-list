package handlers

import (
	"time"
	"net/http"
	"log"
    
	rep "github.com/MilaSnetkova/TODO-list/internal/repeat"
	"github.com/MilaSnetkova/TODO-list/internal/constants"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowParsed, err := time.Parse(constants.DateFormat, now)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	nextDate, err := rep.NextDate(nowParsed, date, repeat)
	if err != nil {
		http.Error(w, "Cannot calculate the next date", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}