package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tasks/internal/entities"
	"tasks/internal/storage"
	"tasks/pkg/auth"
	"testing"
)

func TestCreateTask_Success(t *testing.T) {
	storage := storage.NewStorage()
	tasksResource := &TasksResource{Storage: storage}

	testTask := entities.Task{Value: "Test Task", Done: false}

	body, _ := json.Marshal(testTask)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	tasksResource.CreateTask(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}
}

func TestCreateTask_BadRequest(t *testing.T) {
	storage := storage.NewStorage()
	tasksResource := &TasksResource{Storage: storage}

	invalidTaskJson := `{"value": 1, "done": false}`

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte(invalidTaskJson)))
	w := httptest.NewRecorder()

	tasksResource.CreateTask(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, res.StatusCode)
	}
}

func TestGetTasks_Success(t *testing.T) {
	storage := storage.NewStorage()
	tasksResource := &TasksResource{Storage: storage}

	testTask := entities.Task{Value: "Test Task", Done: false}

	storage.CreateTask(testTask)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	tasksResource.GetTasks(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	var result map[string][]entities.Task
	err := json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if len(result["tasks"]) == 0 {
		t.Error("expected tasks in response, but got none")
	}
}

func TestDeleteTask_Success(t *testing.T) {
	storage := storage.NewStorage()
	testTask := entities.Task{Value: "Task to be deleted", Done: false}

	taskID, _ := storage.CreateTask(testTask)

	_ = httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
	rec := httptest.NewRecorder()

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}
}


func TestGetTasksWithAuth(t *testing.T) {
	storage := storage.NewStorage()
	tasks := &TasksResource{Storage: storage}
	auth := auth.Auth{Storage: storage}

	testUser := entities.User{Username: "testuser", Password: "testpass"}
	storage.CreateUser(testUser)

	testTasks := []entities.Task{
		{Value: "Task 1", Done: false},
		{Value: "Task 2", Done: true},
	}

	for _, t := range testTasks {
		storage.CreateTask(t)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", auth.CheckAuth(tasks.GetTasks))

	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.SetBasicAuth("testuser", "testpass")

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string][]entities.Task

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
}
