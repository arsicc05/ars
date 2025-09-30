package handlers

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

type cachedResponse struct {
	status int
	head   http.Header
	body   []byte
	exp    time.Time
}

type idempotencyStore struct {
	mu   sync.Mutex
	data map[string]cachedResponse
}

func newIdempotencyStore() *idempotencyStore {
	return &idempotencyStore{data: make(map[string]cachedResponse)}
}

func (s *idempotencyStore) get(key string) (cachedResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp, ok := s.data[key]
	if !ok || time.Now().After(resp.exp) {
		return cachedResponse{}, false
	}
	return resp, true
}

func (s *idempotencyStore) set(key string, resp cachedResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = resp
	// naive cleanup could be added here periodically
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.status = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	_, _ = rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func IdempotencyMiddleware(ttl time.Duration) func(http.Handler) http.Handler {
	store := newIdempotencyStore()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			if cached, ok := store.get(key); ok {
				for k, vv := range cached.head {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				w.WriteHeader(cached.status)
				_, _ = w.Write(cached.body)
				return
			}
			recorder := &responseRecorder{ResponseWriter: w, status: 200}
			next.ServeHTTP(recorder, r)
			// copy headers
			headCopy := make(http.Header)
			for k, vv := range recorder.Header() {
				copySlice := make([]string, len(vv))
				copy(copySlice, vv)
				headCopy[k] = copySlice
			}
			store.set(key, cachedResponse{
				status: recorder.status,
				head:   headCopy,
				body:   recorder.body.Bytes(),
				exp:    time.Now().Add(ttl),
			})
		})
	}
}


