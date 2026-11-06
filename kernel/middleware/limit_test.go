package middleware

import (
	"context"
	"fmt"
	"labflow/bootstrap"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiter(redisClient, RateLimitConfig{
		Namespace: "test",
		Limit:     2,
		Window:    time.Second * 10,
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

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

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Mock authentication middleware setting "id" to 42
	r.Use(func(c *gin.Context) {
		c.Set("id", 42)
		c.Next()
	})
	r.Use(RateLimiter(redisClient, RateLimitConfig{
		Namespace: "test",
		Limit:     1,
		Window:    time.Second * 10,
	}))
	r.GET("/api/resource", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

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

func TestRateLimiterAuthenticatedUserDifferentTypes(t *testing.T) {
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

	tests := []struct {
		name        string
		userID      interface{}
		expectedKey string
	}{
		{name: "Int64 ID", userID: int64(100), expectedKey: "rate_limit:test:user:100"},
		{name: "Float64 ID", userID: float64(200.0), expectedKey: "rate_limit:test:user:200"},
		{name: "String ID", userID: "custom-user-uuid", expectedKey: "rate_limit:test:user:custom-user-uuid"},
		{name: "Zero Int ID fallback to IP", userID: 0, expectedKey: "rate_limit:test:ip:127.0.0.1"},
		{name: "Empty String ID fallback to IP", userID: "", expectedKey: "rate_limit:test:ip:127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient.Del(ctx, tt.expectedKey)
			t.Cleanup(func() {
				redisClient.Del(ctx, tt.expectedKey)
			})

			gin.SetMode(gin.TestMode)
			r := gin.New()

			r.Use(func(c *gin.Context) {
				c.Set("id", tt.userID)
				c.Next()
			})
			r.Use(RateLimiter(redisClient, RateLimitConfig{
				Namespace: "test",
				Limit:     1,
				Window:    time.Second * 10,
			}))
			r.GET("/api/resource", func(c *gin.Context) {
				c.String(http.StatusOK, "OK")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/resource", nil)
			req.RemoteAddr = "127.0.0.1:1234"
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status OK, got %v", w.Code)
			}

			exists, err := redisClient.Exists(ctx, tt.expectedKey).Result()
			if err != nil || exists == 0 {
				t.Errorf("Expected Redis key %s to exist", tt.expectedKey)
			}
		})
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

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Simulate two limiters with different namespaces sharing the same IP
	r.Use(RateLimiter(redisClient, RateLimitConfig{
		Namespace: "global",
		Limit:     1,
		Window:    time.Second * 10,
	}))
	r.Use(RateLimiter(redisClient, RateLimitConfig{
		Namespace: "auth",
		Limit:     1,
		Window:    time.Second * 10,
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

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
