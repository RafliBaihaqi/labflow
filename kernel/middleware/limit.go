package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

func RateLimiter(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	disabled := rateLimitDisabled()
	return func(c *gin.Context) {
		if disabled || redisClient == nil {
			c.Next()
			return
		}

		ns := config.Namespace
		if ns == "" {
			ns = "default"
		}

		// Key resolution: check for authenticated user ID (int, int64, float64, or string) set by JWT middleware
		var key string
		if userIDVal, exists := c.Get("id"); exists {
			switch val := userIDVal.(type) {
			case int:
				if val > 0 {
					key = fmt.Sprintf("rate_limit:%s:user:%d", ns, val)
				}
			case int64:
				if val > 0 {
					key = fmt.Sprintf("rate_limit:%s:user:%d", ns, val)
				}
			case float64:
				if val > 0 {
					key = fmt.Sprintf("rate_limit:%s:user:%d", ns, int(val))
				}
			case string:
				if val != "" {
					key = fmt.Sprintf("rate_limit:%s:user:%s", ns, val)
				}
			}
		}

		if key == "" {
			key = fmt.Sprintf("rate_limit:%s:ip:%s", ns, c.ClientIP())
		}

		now := time.Now().UnixMilli()
		windowMs := config.Window.Milliseconds()
		member := uuid.NewString()

		// Execute atomic Lua script
		res, err := redisClient.Eval(c.Request.Context(), slidingWindowLuaScript, []string{key}, now, windowMs, config.Limit, member).Result()
		if err != nil {
			log.Printf("[RateLimiter] redis error (failing open): %v", err)
			c.Next()
			return
		}

		results, ok := res.([]interface{})
		if !ok || len(results) < 2 {
			log.Printf("[RateLimiter] unexpected redis response format (failing open)")
			c.Next()
			return
		}

		allowed, ok1 := results[0].(int64)
		remaining, ok2 := results[1].(int64)
		if !ok1 || !ok2 {
			log.Printf("[RateLimiter] unexpected redis result types (failing open)")
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if allowed == 0 {
			// Estimate retry after: oldest score in ZSET + window duration - now
			retryAfterSec := int(config.Window.Seconds())
			ctx := c.Request.Context()
			oldest, err := redisClient.ZRangeWithScores(ctx, key, 0, 0).Result()
			if err == nil && len(oldest) > 0 {
				elapsedMs := now - int64(oldest[0].Score)
				remainingMs := windowMs - elapsedMs
				if remainingMs > 0 {
					retryAfterSec = int(remainingMs/1000) + 1
				}
			}

			c.Header("Retry-After", strconv.Itoa(retryAfterSec))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded. please try again later.",
			})
			return
		}

		c.Next()
	}
}
