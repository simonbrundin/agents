package main

import (
	"encoding/json"
)

// Simplified OpenRouter response model for chat completions (MVP)
// This will need adjustments if OpenRouter response shape differs.

type OpenRouterChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type OpenRouterResponse struct {
	ID      string             `json:"id"`
	Choices []OpenRouterChoice `json:"choices"`
}

func extractOpenRouterText(b []byte) (string, error) {
	var resp OpenRouterResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}
