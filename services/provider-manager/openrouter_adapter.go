package main

import (
	"bytes"
	"context"
	"encoding/json"
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
		return "", ErrNoAPIKey{}
	}
	body := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// Try to parse OpenRouter JSON and return the first choice content.
	if txt, err := extractOpenRouterText(respBody); err == nil && txt != "" {
		return txt, nil
	}
	// Fall back to returning raw body if parsing yields no content or errors
	return string(respBody), nil
}
