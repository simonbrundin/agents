package main

import (
	"testing"
)

func TestExtractOpenRouterText(t *testing.T) {
	raw := []byte(`{"id":"1","choices":[{"message":{"content":"Hello world"}}]}`)
	got, err := extractOpenRouterText(raw)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != "Hello world" {
		t.Fatalf("unexpected text: %s", got)
	}
}
