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
		headers    map[string]string
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
		{
			name:       "X-Forwarded-For takes precedence",
			remoteAddr: "10.0.0.1:1234",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.7, 10.0.0.1",
			},
			expected: "203.0.113.7",
		},
		{
			name:       "X-Real-IP is used when forwarded is absent",
			remoteAddr: "10.0.0.1:1234",
			headers: map[string]string{
				"X-Real-IP": "198.51.100.8",
			},
			expected: "198.51.100.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

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

func TestMiddleware_WithCustomKeyFunc(t *testing.T) {
	limiter, err := New(1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	keyFunc := func(r *http.Request) string {
		return r.Header.Get("X-Auth-User")
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := MiddlewareWithKeyFunc(limiter, keyFunc)(next)

	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.Header.Set("X-Auth-User", "user-1")
	req1.RemoteAddr = "10.0.0.1:1234"

	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("X-Auth-User", "user-1")
	req2.RemoteAddr = "10.0.0.2:4321"

	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected rate limit for same custom key, got %d", rec2.Code)
	}

	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req3.Header.Set("X-Auth-User", "user-2")
	req3.RemoteAddr = "10.0.0.2:4321"

	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Fatalf("expected different custom key to be allowed, got %d", rec3.Code)
	}
}
