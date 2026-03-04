package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/v1/generate", generateHandler)

	addr := ":8080"
	log.Println("provider-manager listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type GenRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req GenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Default model
	if req.Model == "" {
		req.Model = "gpt-5-mini"
	}

	// Try OpenRouter first
	respBody, status, err := callOpenRouter(ctx, req)
	if err != nil {
		log.Printf("openrouter error: %v", err)
		http.Error(w, "provider error", http.StatusBadGateway)
		return
	}

	w.WriteHeader(status)
	w.Write(respBody)
}

func callOpenRouter(ctx context.Context, req GenRequest) ([]byte, int, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	apiURL := os.Getenv("OPENROUTER_API_URL")
	if apiURL == "" {
		apiURL = "https://api.openrouter.ai/v1/chat/completions"
	}
	if apiKey == "" {
		return nil, 0, ErrNoAPIKey{}
	}

	body := map[string]interface{}{
		"model":    req.Model,
		"messages": []map[string]string{{"role": "user", "content": req.Prompt}},
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	if req.Temperature != 0 {
		body["temperature"] = req.Temperature
	}

	b, _ := json.Marshal(body)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(b))
	if err != nil {
		return nil, 0, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return respBody, res.StatusCode, nil
}

// ErrNoAPIKey simple error type
type ErrNoAPIKey struct{}

func (ErrNoAPIKey) Error() string { return "OPENROUTER_API_KEY not set" }
