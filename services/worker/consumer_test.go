package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	mr "github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestWorkerFlow(t *testing.T) {
	// Start miniredis
	m, err := mr.Run()
	if err != nil {
		t.Fatalf("miniredis failed: %v", err)
	}
	defer m.Close()

	os.Setenv("REDIS_ADDR", m.Addr())

	// Mock provider-manager
	pm := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"text":"mock patch content"}`))
	}))
	defer pm.Close()
	os.Setenv("PROVIDER_MANAGER_URL", pm.URL+"/v1/generate")

	// Push a job
	r := redis.NewClient(&redis.Options{Addr: m.Addr()})
	_, _ = r.XAdd(context.Background(), &redis.XAddArgs{Stream: "agent-jobs", Values: map[string]interface{}{"payload": "{\"repo\":\"simonbrundin/agents\",\"issue\":42}"}}).Result()

	// Run worker in background
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		runWorker(ctx)
	}()

	// wait a bit for worker to process
	time.Sleep(2 * time.Second)
	cancel()
}
