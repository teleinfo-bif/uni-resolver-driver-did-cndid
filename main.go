package main

import (
	"log"
	"net/http"
	"uni-resolver-driver-did-cndid/document"
)

func main() {
	http.HandleFunc("/1.0/identifiers/", document.ResolveDID)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
