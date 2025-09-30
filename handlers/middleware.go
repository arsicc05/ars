package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Simple fixed-window per-IP rate limiter
type rateLimiter struct {
	limit     int
	window    time.Duration
	mu        sync.Mutex
	buckets   map[string]*bucket
}

type bucket struct {
	count     int
	windowEnd time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:   limit,
		window:  window,
		buckets: make(map[string]*bucket),
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[ip]
	now := time.Now()
	if !ok || now.After(b.windowEnd) {
		rl.buckets[ip] = &bucket{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}
	if b.count < rl.limit {
		b.count++
		return true
	}
	return false
}

func getIP(r *http.Request) string {
	// Try X-Forwarded-For first if behind proxy; fallback to RemoteAddr
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		return xf
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	rl := newRateLimiter(limit, window)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			if !rl.allow(ip) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}


