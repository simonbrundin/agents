package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGithubWebhook_Enqueue(t *testing.T) {
	// Mock ENQUEUE_URL server to capture payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		body := string(bodyBytes)
		if !strings.Contains(body, "repo") {
			t.Fatalf("unexpected body: %s", body)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	os.Setenv("ENQUEUE_URL", srv.URL)

	payload := `{"action":"opened","issue":{"number":42,"title":"Test","body":"hi"},"repository":{"full_name":"simonbrundin/agents"}}`
	req := httptest.NewRequest("POST", "/webhook/github", strings.NewReader(payload))
	w := httptest.NewRecorder()
	githubWebhook(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}
