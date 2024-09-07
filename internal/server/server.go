package server

import (
	"log"
	"net/http"

	"github.com/SevvyP/tasks_v1/internal/database"
)

type Resolver struct {
	Server   http.Server
	Database database.TaskDatabase
}

func NewResolver() *Resolver {
	mux := http.NewServeMux()
	database, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	resolver := &Resolver{
		Server: http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		Database: database,
	}
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			resolver.GetTasks(w, r)
		case http.MethodPost:
			resolver.CreateTask(w, r)
		case http.MethodPut:
			resolver.UpdateTask(w, r)
		case http.MethodDelete:
			resolver.DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return resolver
}

func (r *Resolver) Resolve() error {
	return r.Server.ListenAndServe()
}
