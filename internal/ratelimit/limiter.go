package ratelimit

import (
	"errors"
	"sync"
	"time"
)

type Limiter struct {
	mu sync.Mutex

	limit  int
	window time.Duration

	clients   map[string]*client
	stop      chan struct{}
	closeOnce sync.Once
}

type client struct {
	count       int
	windowStart time.Time
}

func New(limit int, window time.Duration) (*Limiter, error) {
	if limit <= 0 {
		return nil, errors.New("rate limit must be greater than zero")
	}

	if window <= 0 {
		return nil, errors.New("rate limit window must be greater than zero")
	}

	l := &Limiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*client),
		stop:    make(chan struct{}),
	}

	go l.cleanup()

	return l, nil
}

func (l *Limiter) cleanup() {
	ticker := time.NewTicker(l.window)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.removeExpiredClients()

		case <-l.stop:
			return
		}
	}
}

func (l *Limiter) removeExpiredClients() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	for key, c := range l.clients {
		if now.Sub(c.windowStart) >= l.window {
			delete(l.clients, key)
		}
	}
}

func (l *Limiter) Allow(key string) (bool, int, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	c, exists := l.clients[key]

	if !exists || now.Sub(c.windowStart) >= l.window {
		l.clients[key] = &client{
			count:       1,
			windowStart: now,
		}

		return true, l.limit - 1, 0
	}

	if c.count >= l.limit {
		retryAfter := l.window - now.Sub(c.windowStart)

		return false, 0, retryAfter
	}

	c.count++

	return true, l.limit - c.count, 0
}

func (l *Limiter) Close() {
	l.closeOnce.Do(func() {
		close(l.stop)
	})
}
