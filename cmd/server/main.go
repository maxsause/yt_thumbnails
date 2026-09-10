package main

import (
	"log"
	"main/internal/database"
	"main/internal/server"
)

func main() {
	db, err := database.NewDatabase("thumbnail.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := "localhost:50051"
	log.Printf("server starting at %s", port)
	if err = server.Start(port, db); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
