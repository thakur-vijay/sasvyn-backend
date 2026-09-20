package ratelimit

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	authRefreshPath = "/auth/refresh"
	socialLoginPath = "/auth/socialLogin"
)

type KeyFunc func(*http.Request) string

func Middleware(limiter *Limiter) func(http.Handler) http.Handler {
	return MiddlewareWithKeyFunc(limiter, clientIP)
}

func MiddlewareWithKeyFunc(limiter *Limiter, keyFunc KeyFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if key == "" {
				key = clientIP(r)
			}

			allowed, remaining, retryAfter := limiter.Allow(key)

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if !allowed {
				retrySeconds := int(retryAfter.Seconds())

				if retrySeconds < 1 {
					retrySeconds = 1
				}

				w.Header().Set("Retry-After", strconv.Itoa(retrySeconds))

				http.Error(
					w,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		first := strings.TrimSpace(strings.Split(xff, ",")[0])
		if first != "" {
			return strings.Trim(first, "[]")
		}
	}

	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return strings.Trim(xrip, "[]")
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

func PolicyMiddleware(
	defaultLimiter *Limiter,
	authRefreshLimiter *Limiter,
	socialLoginLimiter *Limiter,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := defaultLimiter

			switch {
			case r.Method == http.MethodPost && hasRoutePath(r.URL.Path, authRefreshPath):
				limiter = authRefreshLimiter

			case r.Method == http.MethodPost && hasRoutePath(r.URL.Path, socialLoginPath):
				limiter = socialLoginLimiter
			}

			key := clientIP(r)

			allowed, remaining, retryAfter := limiter.Allow(key)

			w.Header().Set(
				"X-RateLimit-Limit",
				strconv.Itoa(limiter.limit),
			)
			w.Header().Set(
				"X-RateLimit-Remaining",
				strconv.Itoa(remaining),
			)

			if !allowed {
				retrySeconds := int(retryAfter.Seconds())

				if retrySeconds < 1 {
					retrySeconds = 1
				}

				w.Header().Set(
					"Retry-After",
					strconv.Itoa(retrySeconds),
				)

				http.Error(
					w,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func PolicyMiddlewareWithKeyFunc(
	defaultLimiter *Limiter,
	authRefreshLimiter *Limiter,
	socialLoginLimiter *Limiter,
	keyFunc KeyFunc,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := defaultLimiter

			switch {
			case r.Method == http.MethodPost && hasRoutePath(r.URL.Path, authRefreshPath):
				limiter = authRefreshLimiter

			case r.Method == http.MethodPost && hasRoutePath(r.URL.Path, socialLoginPath):
				limiter = socialLoginLimiter
			}

			key := keyFunc(r)
			if key == "" {
				key = clientIP(r)
			}

			allowed, remaining, retryAfter := limiter.Allow(key)

			w.Header().Set(
				"X-RateLimit-Limit",
				strconv.Itoa(limiter.limit),
			)
			w.Header().Set(
				"X-RateLimit-Remaining",
				strconv.Itoa(remaining),
			)

			if !allowed {
				retrySeconds := int(retryAfter.Seconds())

				if retrySeconds < 1 {
					retrySeconds = 1
				}

				w.Header().Set(
					"Retry-After",
					strconv.Itoa(retrySeconds),
				)

				http.Error(
					w,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasRoutePath(path, route string) bool {
	return path == route || strings.HasSuffix(path, route)
}
