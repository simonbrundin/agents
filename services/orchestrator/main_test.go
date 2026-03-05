package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGithubWebhook_MissingEnv(t *testing.T) {
	t.Setenv("REDIS_URL", "")
	t.Setenv("ENQUEUE_URL", "")

	payload := `{"action":"opened","issue":{"number":1,"title":"Test","body":"body"},"repository":{"full_name":"owner/repo"}}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/github", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	handleGithubWebhook(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHandleGithubWebhook_InvalidPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook/github", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()

	handleGithubWebhook(rec, req)

	if rec.Code != 400 {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandleGithubWebhook_EnqueueURLSet(t *testing.T) {
	t.Setenv("REDIS_URL", "redis://localhost")
	t.Setenv("ENQUEUE_URL", "http://localhost:9999/enqueue")

	payload := `{"action":"opened","issue":{"number":1,"title":"Test","body":"body"},"repository":{"full_name":"owner/repo"}}`
	req := httptest.NewRequest(http.MethodPost, "/webhook/github", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	handleGithubWebhook(rec, req)

	if rec.Code != 500 {
		t.Errorf("expected status 500 (failed to post), got %d", rec.Code)
	}
}
