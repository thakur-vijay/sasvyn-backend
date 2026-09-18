package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestMiddleware_AllowsRequest(t *testing.T) {
	limiter, err := New(5, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(limiter)(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.10:54321"

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Header().Get("X-RateLimit-Limit"); got != "5" {
		t.Fatalf("expected X-RateLimit-Limit=5, got %q", got)
	}

	if got := rec.Header().Get("X-RateLimit-Remaining"); got != "4" {
		t.Fatalf("expected X-RateLimit-Remaining=4, got %q", got)
	}

	if got := rec.Header().Get("Retry-After"); got != "" {
		t.Fatalf("Retry-After should not be present, got %q", got)
	}
}

func TestMiddleware_BlocksRequest(t *testing.T) {
	limiter, err := New(1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	nextCalls := 0

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(limiter)(next)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.10:54321"

	rec1 := httptest.NewRecorder()

	handler.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.10:54322"

	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusTooManyRequests,
			rec2.Code,
		)
	}

	if nextCalls != 1 {
		t.Fatalf(
			"expected next handler to be called once, got %d",
			nextCalls,
		)
	}

	if got := rec2.Header().Get("X-RateLimit-Limit"); got != "1" {
		t.Fatalf("expected X-RateLimit-Limit=1, got %q", got)
	}

	if got := rec2.Header().Get("X-RateLimit-Remaining"); got != "0" {
		t.Fatalf("expected X-RateLimit-Remaining=0, got %q", got)
	}

	retryAfter := rec2.Header().Get("Retry-After")

	if retryAfter == "" {
		t.Fatal("Retry-After header should be present")
	}

	seconds, err := strconv.Atoi(retryAfter)
	if err != nil {
		t.Fatalf("Retry-After should contain seconds, got %q", retryAfter)
	}

	if seconds < 1 {
		t.Fatalf("Retry-After should be at least 1 second, got %d", seconds)
	}
}

func TestMiddleware_SeparatesClients(t *testing.T) {
	limiter, err := New(1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Middleware(limiter)(next)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "user-1"

	rec1 := httptest.NewRecorder()

	handler.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "user-2"

	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec2, req2)

	if rec1.Code != http.StatusOK {
		t.Fatalf("user-1 expected status %d, got %d", http.StatusOK, rec1.Code)
	}

	if rec2.Code != http.StatusOK {
		t.Fatalf("user-2 expected status %d, got %d", http.StatusOK, rec2.Code)
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		expected   string
	}{
		{
			name:       "IPv4 with port",
			remoteAddr: "192.168.1.10:54321",
			expected:   "192.168.1.10",
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[::1]:54321",
			expected:   "::1",
		},
		{
			name:       "IP without port",
			remoteAddr: "192.168.1.10",
			expected:   "192.168.1.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			got := clientIP(req)

			if got != tt.expected {
				t.Fatalf(
					"expected IP %q, got %q",
					tt.expected,
					got,
				)
			}
		})
	}
}
