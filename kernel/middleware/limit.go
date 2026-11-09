package middleware

import (
	"fmt"
	"labflow/httpx"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RateLimitConfig struct {
	Namespace string
	Limit     int
	Window    time.Duration
}

// rateLimitDisabled lets load tests run from a single IP without being throttled.
// It is refused in production.
func rateLimitDisabled() bool {
	if os.Getenv("DISABLE_RATE_LIMIT") != "true" {
		return false
	}
	if os.Getenv("ENV") == "production" || os.Getenv("GO_ENV") == "production" {
		log.Println("[ratelimit] DISABLE_RATE_LIMIT ignored in production")
		return false
	}
	log.Println("[ratelimit] rate limiting DISABLED (DISABLE_RATE_LIMIT=true)")
	return true
}

func RateLimiter(redisClient *redis.Client, config RateLimitConfig) httpx.Middleware {
	disabled := rateLimitDisabled()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if disabled || redisClient == nil {
				next.ServeHTTP(w, r)
				return
			}

			ns := config.Namespace
			if ns == "" {
				ns = "default"
			}

			// Key by the authenticated user when the JWT middleware has run, else by IP.
			var key string
			if userID, ok := httpx.UserID(r.Context()); ok && userID > 0 {
				key = fmt.Sprintf("rate_limit:%s:user:%d", ns, userID)
			} else {
				key = fmt.Sprintf("rate_limit:%s:ip:%s", ns, httpx.ClientIP(r))
			}

			now := time.Now().UnixMilli()
			windowMs := config.Window.Milliseconds()
			member := uuid.NewString()

			// Execute atomic Lua script
			res, err := redisClient.Eval(r.Context(), slidingWindowLuaScript, []string{key}, now, windowMs, config.Limit, member).Result()
			if err != nil {
				log.Printf("[RateLimiter] redis error (failing open): %v", err)
				next.ServeHTTP(w, r)
				return
			}

			results, ok := res.([]interface{})
			if !ok || len(results) < 2 {
				log.Printf("[RateLimiter] unexpected redis response format (failing open)")
				next.ServeHTTP(w, r)
				return
			}

			allowed, ok1 := results[0].(int64)
			remaining, ok2 := results[1].(int64)
			if !ok1 || !ok2 {
				log.Printf("[RateLimiter] unexpected redis result types (failing open)")
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

			if allowed == 0 {
				// Estimate retry after: oldest score in ZSET + window duration - now
				retryAfterSec := int(config.Window.Seconds())
				oldest, err := redisClient.ZRangeWithScores(r.Context(), key, 0, 0).Result()
				if err == nil && len(oldest) > 0 {
					elapsedMs := now - int64(oldest[0].Score)
					remainingMs := windowMs - elapsedMs
					if remainingMs > 0 {
						retryAfterSec = int(remainingMs/1000) + 1
					}
				}

				w.Header().Set("Retry-After", strconv.Itoa(retryAfterSec))
				httpx.WriteJSON(w, http.StatusTooManyRequests, map[string]string{
					"error": "rate limit exceeded. please try again later.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
