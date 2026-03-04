package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func callOpenRouterRaw(ctx context.Context, prompt string, model string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	apiURL := os.Getenv("OPENROUTER_API_URL")
	if apiURL == "" {
		apiURL = "https://api.openrouter.ai/v1/chat/completions"
	}
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY not set")
	}
	body := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)
	// TODO: parse response and return text only
	return string(respBody), nil
}
