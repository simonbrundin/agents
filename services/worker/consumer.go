package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func runWorker(ctx context.Context) {
	r := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	defer r.Close()

	stream := "agent-jobs"
	for {
		// XREAD BLOCK 0 STREAMS agent-jobs $
		res, err := r.XRead(ctx, &redis.XReadArgs{Streams: []string{stream, "0-0"}, Count: 1, Block: 5 * time.Second}).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("redis read error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, str := range res {
			for _, msg := range str.Messages {
				payload := msg.Values["payload"].(string)
				log.Printf("worker got payload: %s", payload)
				// Simulate calling provider-manager to get patch
				patch, err := requestPatch(ctx, payload)
				if err != nil {
					log.Printf("provider-manager error: %v", err)
					continue
				}
				log.Printf("simulated apply patch: %s", patch)
				// Simulate creating PR
				log.Printf("simulated create PR for payload: %s", payload)
			}
		}
	}
}

func requestPatch(ctx context.Context, payload string) (string, error) {
	pmURL := os.Getenv("PROVIDER_MANAGER_URL")
	if pmURL == "" {
		pmURL = "http://localhost:8080/v1/generate"
	}
	reqBody := map[string]interface{}{"prompt": fmt.Sprintf("Fix the issue described: %s", payload)}
	b, _ := json.Marshal(reqBody)
	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Post(pmURL, "application/json", bytesNewReader(b))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var out map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	if t, ok := out["text"].(string); ok {
		return t, nil
	}
	return "", fmt.Errorf("no text in provider response")
}

func bytesNewReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
