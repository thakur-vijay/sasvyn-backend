package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	newHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected JSON content type, got %q", contentType)
	}
	if body := recorder.Body.String(); body != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected response body: %q", body)
	}
}
