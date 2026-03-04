package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple stub server for provider-manager
	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("provider-manager listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
