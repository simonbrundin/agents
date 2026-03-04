package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGenerateHandler_OpenRouter(t *testing.T) {
	// Setup mock OpenRouter server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"1","choices":[{"message":{"content":"Generated text from mock"}}]}`))
	}))
	defer srv.Close()

	os.Setenv("OPENROUTER_API_URL", srv.URL)
	os.Setenv("OPENROUTER_API_KEY", "testkey")

	reqBody := strings.NewReader(`{"prompt":"hello"}`)
	req := httptest.NewRequest("POST", "/v1/generate", reqBody)
	w := httptest.NewRecorder()
	generateHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}

	// simple check for body containing text field
	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)
	if !strings.Contains(body, "Generated text from mock") {
		t.Fatalf("unexpected body: %s", body)
	}
}
