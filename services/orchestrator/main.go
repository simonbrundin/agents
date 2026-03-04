package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("orchestrator listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
