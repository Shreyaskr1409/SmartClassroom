package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", NewServer()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
