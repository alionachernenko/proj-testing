package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tasks/internal/entities"
	"tasks/internal/storage"
	"tasks/pkg/auth"
)

func main() {
	mux := http.NewServeMux()

	storage := storage.NewStorage()

	tasks := TasksResource{
		Storage: storage,
	}

	auth := auth.Auth{
		Storage: storage,
	}

	users := UsersResource{
		Storage: storage,
	}
	
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /tasks", auth.CheckAuth(tasks.GetTasks))
	mux.HandleFunc("POST /tasks", auth.CheckAuth(tasks.CreateTask))
	mux.HandleFunc("PUT /tasks/{id}", auth.CheckAuth(tasks.UpdateTask))
	mux.HandleFunc("DELETE /tasks/{id}", auth.CheckAuth(tasks.DeleteTask))

	http.ListenAndServe(":8080", mux)
}

type TasksResource struct {
	Storage *storage.Storage
}

func (tr *TasksResource) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := tr.Storage.GetTasks()
	res := map[string][]entities.Task{"tasks": tasks}

	err := json.NewEncoder(w).Encode(res)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (tr *TasksResource) CreateTask(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Println("Failed to decode: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskID, ok := tr.Storage.CreateTask(reqBody)

	if !ok {
		fmt.Print("Error creating task", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBody.ID = taskID

}

func (tr *TasksResource) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var reqBody entities.Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := tr.Storage.UpdadeTask(id, reqBody)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tr *TasksResource) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ok := tr.Storage.DeleteTask(id)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type UsersResource struct {
	Storage *storage.Storage
}

func (ur *UsersResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println("Failed to decode: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	ok := ur.Storage.CreateUser(user)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
