package middleware

import (
	"context"
	"fmt"
	"labflow/bootstrap"
	"labflow/httpx"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func TestRateLimiterCompilation(t *testing.T) {
	// Test that we can define configuration and initialize the limiter function
	config := RateLimitConfig{
		Namespace: "test",
		Limit:     5,
		Window:    time.Minute,
	}
	if config.Limit != 5 || config.Window != time.Minute {
		t.Errorf("Config values mismatch")
	}

	// Verify the RateLimiter signature compiles and initializes
	limiter := RateLimiter(nil, config)
	if limiter == nil {
		t.Errorf("Expected limiter handler function to be initialized")
	}
}

func TestRateLimiterIPChecking(t *testing.T) {
	// Load env from parent directory (since tests run in the package dir)
	_ = godotenv.Load("../.env")
	bootstrap.LoadEnvVariables()

	var redisClient *redis.Client
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("redis init panic: %v", r)
			}
		}()
		redisClient, err = bootstrap.InitRedis()
	}()

	if err != nil {
		t.Skipf("Redis is not available; skipping integration test: %v", err)
		return
	}
	defer redisClient.Close()

	// Clean any previous test keys
	ctx := context.Background()
	redisClient.Del(ctx, "rate_limit:test:ip:127.0.0.1")
	t.Cleanup(func() {
		redisClient.Del(ctx, "rate_limit:test:ip:127.0.0.1")
	})

	r := httpx.Chain(okHandler, RateLimiter(redisClient, RateLimitConfig{
		Namespace: "test",
		Limit:     2,
		Window:    time.Second * 10,
	}))

	// Request 1: Should be allowed (remaining = 1)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w1.Code)
	}
	if w1.Header().Get("X-RateLimit-Limit") != "2" {
		t.Errorf("Expected limit 2, got %v", w1.Header().Get("X-RateLimit-Limit"))
	}
	if w1.Header().Get("X-RateLimit-Remaining") != "1" {
		t.Errorf("Expected remaining 1, got %v", w1.Header().Get("X-RateLimit-Remaining"))
	}

	// Request 2: Should be allowed (remaining = 0)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w2.Code)
	}
	if w2.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Errorf("Expected remaining 0, got %v", w2.Header().Get("X-RateLimit-Remaining"))
	}

	// Request 3: Should be blocked (StatusTooManyRequests)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", w3.Code)
	}
	retryAfter := w3.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Errorf("Expected Retry-After header to be set")
	} else {
		val, err := strconv.Atoi(retryAfter)
		if err != nil || val <= 0 {
			t.Errorf("Expected Retry-After to be a positive integer, got %v", retryAfter)
		}
	}
}

func TestRateLimiterAuthenticatedUser(t *testing.T) {
	_ = godotenv.Load("../.env")
	bootstrap.LoadEnvVariables()

	var redisClient *redis.Client
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("redis init panic: %v", r)
			}
		}()
		redisClient, err = bootstrap.InitRedis()
	}()

	if err != nil {
		t.Skipf("Redis is not available; skipping integration test: %v", err)
		return
	}
	defer redisClient.Close()

	ctx := context.Background()
	redisClient.Del(ctx, "rate_limit:test:user:42")
	t.Cleanup(func() {
		redisClient.Del(ctx, "rate_limit:test:user:42")
	})

	// Authenticated as user 42, as the JWT middleware would do.
	r := httpx.Chain(okHandler, asUser(42), RateLimiter(redisClient, RateLimitConfig{
		Namespace: "test",
		Limit:     1,
		Window:    time.Second * 10,
	}))

	// Request 1: Allowed
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/resource", nil)
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w1.Code)
	}

	// Request 2: Blocked
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/resource", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", w2.Code)
	}

	// Verify key was created under user ID 42 in Redis
	exists, err := redisClient.Exists(ctx, "rate_limit:test:user:42").Result()
	if err != nil || exists == 0 {
		t.Errorf("Expected Redis key rate_limit:user:42 to exist")
	}
}

// The user ID in the context is a typed int, so the limiter no longer needs to
// handle int64, float64 or string IDs. A zero ID is not a real user: key by IP.
func TestRateLimiterZeroUserFallsBackToIP(t *testing.T) {
	_ = godotenv.Load("../.env")
	bootstrap.LoadEnvVariables()

	var redisClient *redis.Client
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("redis init panic: %v", r)
			}
		}()
		redisClient, err = bootstrap.InitRedis()
	}()

	if err != nil {
		t.Skipf("Redis is not available; skipping integration test: %v", err)
		return
	}
	defer redisClient.Close()

	ctx := context.Background()
	key := "rate_limit:test:ip:127.0.0.1"
	redisClient.Del(ctx, key)
	t.Cleanup(func() {
		redisClient.Del(ctx, key)
	})

	r := httpx.Chain(okHandler, asUser(0), RateLimiter(redisClient, RateLimitConfig{
		Namespace: "test",
		Limit:     1,
		Window:    time.Second * 10,
	}))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/resource", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		t.Errorf("Expected Redis key %s to exist", key)
	}
}

func TestRateLimiterNamespaceIsolation(t *testing.T) {
	_ = godotenv.Load("../.env")
	bootstrap.LoadEnvVariables()

	var redisClient *redis.Client
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("redis init panic: %v", r)
			}
		}()
		redisClient, err = bootstrap.InitRedis()
	}()

	if err != nil {
		t.Skipf("Redis is not available; skipping integration test: %v", err)
		return
	}
	defer redisClient.Close()

	ctx := context.Background()
	keyA := "rate_limit:global:ip:127.0.0.1"
	keyB := "rate_limit:auth:ip:127.0.0.1"
	redisClient.Del(ctx, keyA, keyB)
	t.Cleanup(func() {
		redisClient.Del(ctx, keyA, keyB)
	})

	// Simulate two limiters with different namespaces sharing the same IP
	r := httpx.Chain(okHandler,
		RateLimiter(redisClient, RateLimitConfig{
			Namespace: "global",
			Limit:     1,
			Window:    time.Second * 10,
		}),
		RateLimiter(redisClient, RateLimitConfig{
			Namespace: "auth",
			Limit:     1,
			Window:    time.Second * 10,
		}),
	)

	// Request 1: both limiters allow — consumes 1 slot in each namespace key
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w1.Code)
	}

	// Request 2: both limiters should block — each key has count == 1 == limit
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:1234"
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", w2.Code)
	}

	// Verify both keys exist independently
	existsA, _ := redisClient.Exists(ctx, keyA).Result()
	existsB, _ := redisClient.Exists(ctx, keyB).Result()
	if existsA == 0 {
		t.Errorf("Expected Redis key %s to exist", keyA)
	}
	if existsB == 0 {
		t.Errorf("Expected Redis key %s to exist", keyB)
	}
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
})

// asUser stands in for the JWT middleware: it marks the request as authenticated.
func asUser(id int) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), id)))
		})
	}
}
