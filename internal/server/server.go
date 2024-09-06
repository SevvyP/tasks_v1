package server

import "net/http"

type Resolver struct {
	Server http.Server
}

func NewResolver() *Resolver {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", GetTasks)
	mux.HandleFunc("/tasks/create", CreateTask)
	mux.HandleFunc("/tasks/update", UpdateTask)
	mux.HandleFunc("/tasks/delete", DeleteTask)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Resolver{Server: *server}
}

func (r *Resolver) Resolve() error {
	return r.Server.ListenAndServe()
}
