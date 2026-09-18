package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	limiter, err := New(3, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	key := "user-1"

	for i := 0; i < 3; i++ {
		allowed, remaining, retryAfter := limiter.Allow(key)

		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}

		expectedRemaining := 3 - (i + 1)

		if remaining != expectedRemaining {
			t.Fatalf(
				"expected remaining %d, got %d",
				expectedRemaining,
				remaining,
			)
		}

		if retryAfter != 0 {
			t.Fatalf("expected retryAfter to be 0, got %v", retryAfter)
		}
	}
}

func TestLimiter_BlocksRequestsOverLimit(t *testing.T) {
	limiter, err := New(3, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	key := "user-1"

	for i := 0; i < 3; i++ {
		allowed, _, _ := limiter.Allow(key)

		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	allowed, remaining, retryAfter := limiter.Allow(key)

	if allowed {
		t.Fatal("request over limit should be blocked")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter <= 0 {
		t.Fatal("retryAfter should be greater than zero")
	}

	if retryAfter > time.Minute {
		t.Fatalf("retryAfter should not exceed window, got %v", retryAfter)
	}
}

func TestLimiter_SeparatesClients(t *testing.T) {
	limiter, err := New(1, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	allowed, remaining, retryAfter := limiter.Allow("user-1")

	if !allowed {
		t.Fatal("first request for user-1 should be allowed")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter != 0 {
		t.Fatalf("expected retryAfter to be 0, got %v", retryAfter)
	}

	allowed, remaining, retryAfter = limiter.Allow("user-1")

	if allowed {
		t.Fatal("second request for user-1 should be blocked")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter <= 0 {
		t.Fatal("retryAfter should be greater than zero")
	}

	allowed, remaining, retryAfter = limiter.Allow("user-2")

	if !allowed {
		t.Fatal("first request for user-2 should be allowed")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter != 0 {
		t.Fatalf("expected retryAfter to be 0, got %v", retryAfter)
	}
}

func TestLimiter_ResetsAfterWindow(t *testing.T) {
	limiter, err := New(1, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	key := "user-1"

	allowed, remaining, retryAfter := limiter.Allow(key)

	if !allowed {
		t.Fatal("first request should be allowed")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter != 0 {
		t.Fatalf("expected retryAfter to be 0, got %v", retryAfter)
	}

	allowed, remaining, retryAfter = limiter.Allow(key)

	if allowed {
		t.Fatal("second request should be blocked")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0, got %d", remaining)
	}

	if retryAfter <= 0 {
		t.Fatal("retryAfter should be greater than zero")
	}

	time.Sleep(60 * time.Millisecond)

	allowed, remaining, retryAfter = limiter.Allow(key)

	if !allowed {
		t.Fatal("request should be allowed after window reset")
	}

	if remaining != 0 {
		t.Fatalf("expected remaining 0 after reset, got %d", remaining)
	}

	if retryAfter != 0 {
		t.Fatalf("expected retryAfter to be 0 after reset, got %v", retryAfter)
	}
}

func TestLimiter_ConcurrentAccess(t *testing.T) {
	limiter, err := New(100, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer limiter.Close()

	const requests = 1000

	var wg sync.WaitGroup
	var allowed int
	var mu sync.Mutex

	wg.Add(requests)

	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()

			requestAllowed, _, _ := limiter.Allow("user-1")

			if requestAllowed {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowed != 100 {
		t.Fatalf("expected exactly 100 allowed requests, got %d", allowed)
	}
}

func TestLimiter_RejectsInvalidLimit(t *testing.T) {
	_, err := New(0, time.Minute)

	if err == nil {
		t.Fatal("expected error for invalid limit")
	}
}

func TestLimiter_RejectsInvalidWindow(t *testing.T) {
	_, err := New(10, 0)

	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}
