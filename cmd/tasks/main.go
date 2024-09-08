package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/SevvyP/tasks_v1/internal/server"
)

func main() {
	configFile := flag.String("c", "/etc/config.json", "path to config file")
	flag.Parse()
	fmt.Println("Starting server with config file:", *configFile)
	bytes, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	var c server.Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
	err = server.NewResolver(&c).Resolve()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
