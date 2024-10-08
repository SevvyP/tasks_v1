package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SevvyP/tasks_v1/pkg/model"
)

// GetTasks retrieves a list of tasks from the database and sends them as a JSON response.
// If the "id" query parameter is provided, it retrieves a specific task by ID instead.
// The retrieved tasks are encoded as JSON and sent in the response body.
// If there is an error retrieving the tasks from the database, an HTTP 500 Internal Server Error is returned.
// If a task with the specified ID is not found, an HTTP 404 Not Found is returned.
func (r *Resolver) GetTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := req.URL.Query().Get("id")
	user := req.URL.Query().Get("user_id")

	// Check if both "id" and "user" query parameters are present
	if id != "" && user != "" {
		http.Error(w, "Cannot query by both id and user", http.StatusBadRequest)
		return
	}

	// Check for "id" query parameter
	if id != "" {
		r.GetTaskByID(w, req, id)
		return
	}

	// Check for "user" query parameter
	if user != "" {
		r.GetTasksByUserID(w, req, user)
		return
	}

	tasks, err := r.Database.GetTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID retrieves a task by its ID from the database and sends it as a JSON response.
// If the task is found, it is encoded as JSON and sent in the response body.
// If the task is not found, an HTTP 404 Not Found is returned.
// If there is an error retrieving the task from the database, an HTTP 500 Internal Server Error is returned.
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

// GetTasksByUser retrieves tasks by user from the database and sends them as a JSON response.
// If the tasks are found, they are encoded as JSON and sent in the response body.
// If there is an error retrieving the tasks from the database, an HTTP 500 Internal Server Error is returned.
func (r *Resolver) GetTasksByUserID(w http.ResponseWriter, req *http.Request, user string) {
	tasks, err := r.Database.GetTasksByUserID(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask creates a new task in the database based on the JSON request body.
// If the task is created successfully, an HTTP 201 Created response is returned.
// If there is an error creating the task, an HTTP 500 Internal Server Error is returned.
func (r *Resolver) CreateTask(w http.ResponseWriter, req *http.Request) {
	var task model.Task
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

// UpdateTask updates an existing task in the database based on the JSON request body.
// If the task is updated successfully, an HTTP 200 OK response is returned.
// If there is an error updating the task, an HTTP 500 Internal Server Error is returned.
func (r *Resolver) UpdateTask(w http.ResponseWriter, req *http.Request) {
	var updatedTask model.Task
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

// DeleteTask deletes a task from the database based on the JSON request body.
// If the task is deleted successfully, an HTTP 200 OK response is returned.
// If there is an error deleting the task, an HTTP 500 Internal Server Error is returned.
func (r *Resolver) DeleteTask(w http.ResponseWriter, req *http.Request) {
	var taskToDelete model.Task
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
