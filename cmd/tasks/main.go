package main

import (
	"log"

	"github.com/SevvyP/tasks_v1/internal/server"
)

func main() {
	err := server.NewResolver().Resolve()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
