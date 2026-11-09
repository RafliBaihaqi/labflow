package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func ConnectToDB() (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Invalid DATABASE_URL:", err)
	}

	// Bound the pool. Without a cap, every in-flight request under load gets its
	// own connection until Postgres refuses with "too many clients"
	// (default max_connections = 100). Idle connections close after 5 minutes.
	cfg.MaxConns = int32(envInt("DB_MAX_OPEN_CONNS", 25))
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to connect to the DB:", err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to connect to the DB:", err)
	}
	return pool, nil
}

func envInt(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return fallback
}

func InitRedis() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL environment variable is not set")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse REDIS_URL: %v", err))
	}

	client := redis.NewClient(opts)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
	return client, nil
}
