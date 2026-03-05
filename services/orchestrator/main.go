package main

import (
	"bytes"
	"context"
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

	http.HandleFunc("/webhook/github", handleGithubWebhook)

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

func handleGithubWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("failed to read webhook body", err)
		w.WriteHeader(400)
		return
	}

	var ev IssueEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		log.Println("invalid webhook payload", err)
		w.WriteHeader(400)
		return
	}

	log.Printf("received issue event: %s #%d", ev.Repository.FullName, ev.Issue.Number)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("REDIS_URL not set, skipping enqueue (dev mode)")
		w.WriteHeader(200)
		return
	}

	enqueueURL := os.Getenv("ENQUEUE_URL")
	if enqueueURL != "" {
		payload := map[string]interface{}{"repo": ev.Repository.FullName, "issue": ev.Issue}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("failed to marshal payload", err)
			w.WriteHeader(500)
			return
		}
		resp, err := http.Post(enqueueURL, "application/json", bytes.NewReader(b))
		if err != nil {
			log.Println("failed to enqueue job", err)
			w.WriteHeader(500)
			return
		}
		defer resp.Body.Close()
	} else {
		payload := map[string]interface{}{"repo": ev.Repository.FullName, "issue": ev.Issue.Number}
		p, err := json.Marshal(payload)
		if err != nil {
			log.Println("failed to marshal payload", err)
			w.WriteHeader(500)
			return
		}
		if err := pushJob(context.Background(), string(p)); err != nil {
			log.Printf("failed to push job to redis: %v", err)
		}
	}

	w.WriteHeader(200)
}
