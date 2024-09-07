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
	mux.HandleFunc("/tasks", resolver.GetTasks)
	mux.HandleFunc("/tasks/create", resolver.CreateTask)
	mux.HandleFunc("/tasks/update", resolver.UpdateTask)
	mux.HandleFunc("/tasks/delete", resolver.DeleteTask)

	return resolver
}

func (r *Resolver) Resolve() error {
	return r.Server.ListenAndServe()
}
