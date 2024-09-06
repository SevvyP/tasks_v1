package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	ID    string `json:"id"`
	Items []Item `json:"items"`
}

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var tasks = []Task{
	{
		ID: "1",
		Items: []Item{
			{Name: "Task 1 Item 1", Price: 10},
			{Name: "Task 1 Item 2", Price: 20},
		},
	},
	{
		ID: "2",
		Items: []Item{
			{Name: "Task 2 Item 1", Price: 15},
			{Name: "Task 2 Item 2", Price: 25},
		},
	},
}

func GetTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task created successfully")
}

func UpdateTask(w http.ResponseWriter, req *http.Request) {
	var updatedTask Task
	err := json.NewDecoder(req.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			tasks[i] = updatedTask
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Task updated successfully")
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, req *http.Request) {
	var taskToDelete Task
	err := json.NewDecoder(req.Body).Decode(&taskToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == taskToDelete.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Task deleted successfully")
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
