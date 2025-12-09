package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

var store = &TaskStore{
	tasks: make(map[string]*Task),
}

var idCounter int64

func main() {
	http.HandleFunc("/api/tasks", corsMiddleware(loggingMiddleware(tasksHandler)))
	http.HandleFunc("/api/tasks/", corsMiddleware(loggingMiddleware(taskHandler)))
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Task Manager Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "task-manager-backend",
	})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case "GET":
		getAllTasks(w, r)
	case "POST":
		createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.URL.Path[len("/api/tasks/"):]
	if id == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		getTask(w, r, id)
	case "PATCH":
		updateTask(w, r, id)
	case "DELETE":
		deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	tasks := make([]*Task, 0, len(store.tasks))
	for _, task := range store.tasks {
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateTask(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.mu.Lock()
	idCounter++
	task.ID = strconv.FormatInt(idCounter, 10)
	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	store.tasks[task.ID] = &task
	store.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func getTask(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.RLock()
	task, exists := store.tasks[id]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	task, exists := store.tasks[id]
	if !exists {
		store.mu.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var updates Task
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		store.mu.Unlock()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updates.Title != "" {
		task.Title = updates.Title
	}
	if updates.Description != "" {
		task.Description = updates.Description
	}
	if updates.Status != "" {
		task.Status = updates.Status
	}
	if updates.Priority != "" {
		task.Priority = updates.Priority
	}
	task.UpdatedAt = time.Now()
	store.mu.Unlock()

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	delete(store.tasks, id)
	w.WriteHeader(http.StatusNoContent)
}
