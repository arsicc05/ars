package main

import (
	"context"
	"net"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"strconv"
	"syscall"
	"time"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type TokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(refillRate float64, capacity int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	if elapsed > 0 {
		tb.tokens += elapsed * tb.refillRate
		if tb.tokens > float64(tb.capacity) {
			tb.tokens = float64(tb.capacity)
		}
		tb.lastRefill = now
	}

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

type RateLimiter struct {
	buckets    sync.Map // key: client identifier -> *TokenBucket
	refillRate float64
	capacity   int
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{refillRate: rps, capacity: burst}
}

func (rl *RateLimiter) getBucket(key string) *TokenBucket {
	if v, ok := rl.buckets.Load(key); ok {
		return v.(*TokenBucket)
	}
	bucket := NewTokenBucket(rl.refillRate, rl.capacity)
	actual, _ := rl.buckets.LoadOrStore(key, bucket)
	return actual.(*TokenBucket)
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := getClientIP(r)
		bucket := rl.getBucket(client)
		if !bucket.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func getEnvFloat(name string, def float64) float64 {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f <= 0 {
		return def
	}
	return f
}

func getEnvInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil || i <= 0 {
		return def
	}
	return i
}

func main() {
	rps := getEnvFloat("RATE_LIMIT_RPS", 5)
	burst := getEnvInt("RATE_LIMIT_BURST", 10)
	limiter := NewRateLimiter(rps, burst)

	configRepo := repositories.NewConfigInMemRepository()
	groupRepo := repositories.NewConfigGroupInMemRepository()
	
	configService := services.NewConfigService(configRepo)
	groupService := services.NewConfigGroupService(groupRepo)
	
	config := model.NewConfig("db_config", 2)
	config.AddParameter("username", "pera")
	config.AddParameter("port", "5432")
	config.AddParameter("host", "localhost")
	_ = configService.Add(config)
	
	group := model.NewConfigGroup("web_configs", 1)
	
	webConfig := model.NewGroupConfig("web_server")
	webConfig.AddParameter("port", "8080")
	webConfig.AddParameter("host", "0.0.0.0")
	webConfig.AddLabel("environment", "development")
	webConfig.AddLabel("team", "backend")
	
	group.AddConfig(webConfig)
	_ = groupService.Add(group)
	
	configHandler := handlers.NewConfigHandler(configService)
	groupHandler := handlers.NewConfigGroupHandler(groupService)
	
	router := mux.NewRouter()
	router.Use(limiter.Middleware)
	
	router.HandleFunc("/configs", configHandler.GetAll).Methods("GET")
	router.HandleFunc("/configs", configHandler.Create).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")
	
	router.HandleFunc("/groups", groupHandler.GetAll).Methods("GET")
	router.HandleFunc("/groups", groupHandler.Create).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.Get).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.Delete).Methods("DELETE")
	
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.AddConfig).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/configs/{configName}", groupHandler.GetConfig).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}/configs/{configName}", groupHandler.RemoveConfig).Methods("DELETE")
	// labels-based operations
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.GetConfigsByLabels).Methods("GET").Queries("labels", "{labels}")
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.DeleteConfigsByLabels).Methods("DELETE").Queries("labels", "{labels}")

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8000: %v\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
