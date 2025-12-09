package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed in %v", time.Since(start))
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func validateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(task.Title) > 200 {
		return fmt.Errorf("title must be less than 200 characters")
	}
	if task.Priority != "" && task.Priority != "low" && task.Priority != "medium" && task.Priority != "high" {
		return fmt.Errorf("priority must be low, medium, or high")
	}
	if task.Status != "" && task.Status != "pending" && task.Status != "in-progress" && task.Status != "completed" {
		return fmt.Errorf("status must be pending, in-progress, or completed")
	}
	return nil
}

