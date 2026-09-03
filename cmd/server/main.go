package main

import (
	"log"
	"main/internal/server"
)

func main() {
	port := "localhost:50051"
	log.Printf("server starting at %s", port)
	if err := server.Start(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
