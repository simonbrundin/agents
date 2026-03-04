package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/webhook/github", githubWebhook)

	addr := ":8081"
	log.Println("orchestrator listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type IssueEvent struct {
	Action string `json:"action"`
	Issue  struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	} `json:"issue"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

func githubWebhook(w http.ResponseWriter, r *http.Request) {
	// Basic verification skipped for MVP (TODO: verify signature)
	body, _ := io.ReadAll(r.Body)
	var ev IssueEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		log.Println("invalid webhook payload", err)
		w.WriteHeader(400)
		return
	}

	log.Printf("received issue event: %s #%d", ev.Repository.FullName, ev.Issue.Number)

	// For now, enqueue a simple job file to Redis via redis-cli (placeholder)
	// In production use Redis client to push to stream
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("REDIS_URL not set, skipping enqueue (dev mode)")
		w.WriteHeader(200)
		return
	}

	// Simple HTTP call to a queue push endpoint (if present)
	enqueueURL := os.Getenv("ENQUEUE_URL")
	if enqueueURL != "" {
		payload := map[string]interface{}{"repo": ev.Repository.FullName, "issue": ev.Issue}
		b, _ := json.Marshal(payload)
		http.Post(enqueueURL, "application/json", bytes.NewReader(b))
	}

	w.WriteHeader(200)
}
