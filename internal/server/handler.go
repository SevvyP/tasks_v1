package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SevvyP/tasks_v1/internal/database"
)

func (r *Resolver) GetTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.URL.Query().Get("id") != "" {
		r.GetTaskByID(w, req, req.URL.Query().Get("id"))
		return
	}
	tasks, err := r.Database.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (r *Resolver) GetTaskByID(w http.ResponseWriter, req *http.Request, id string) {
	task, err := r.Database.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (r *Resolver) CreateTask(w http.ResponseWriter, req *http.Request) {
	var task database.Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.Database.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task created successfully")
}

func (r *Resolver) UpdateTask(w http.ResponseWriter, req *http.Request) {
	var updatedTask database.Task
	err := json.NewDecoder(req.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.Database.UpdateTask(updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task updated successfully")
}

func (r *Resolver) DeleteTask(w http.ResponseWriter, req *http.Request) {
	var taskToDelete database.Task
	err := json.NewDecoder(req.Body).Decode(&taskToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.Database.DeleteTask(taskToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task deleted successfully")
}
