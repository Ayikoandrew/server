package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		rdb = NewRDB()
	})
	return rdb
}

func InitRedis() {
	GetRedisClient()
}

func NewRDB() *redis.Client {
	url := os.Getenv("REDIS_URL")

	// If REDIS_URL is not provided, construct it from components
	if url == "" {
		redisHost := os.Getenv("REDIS_HOST")
		if redisHost == "" {
			redisHost = "redis" // Default to service name in docker-compose
		}

		redisPort := os.Getenv("REDIS_PORT")
		if redisPort == "" {
			redisPort = "6379" // Default Redis port
		}

		redisPassword := os.Getenv("REDIS_PASSWORD")

		// Format URL with password
		if redisPassword != "" {
			url = fmt.Sprintf("redis://:%s@%s:%s/0", redisPassword, redisHost, redisPort)
		} else {
			url = fmt.Sprintf("redis://%s:%s/0", redisHost, redisPort)
		}
	}

	// Log connection attempt (without showing password)
	slog.Info("Connecting to Redis", "url", maskPasswordInURL(url))

	// Parse URL first
	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Failed to parse Redis URL", "error", err)
		// Don't fall back to direct connection - if URL parsing fails, it's better to fail fast
		panic(fmt.Sprintf("Invalid Redis URL: %v", err))
	}

	// Create client with options
	client := redis.NewClient(opts)

	// Add retry options
	client.Options().MaxRetries = 5
	client.Options().MinRetryBackoff = 100 * time.Millisecond
	client.Options().MaxRetryBackoff = 2 * time.Second

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	return client
}

// maskPasswordInURL masks the password in a Redis URL for safe logging
// Example: "redis://:secret@redis:6379/0" becomes "redis://:****@redis:6379/0"
func maskPasswordInURL(url string) string {
	// Handle empty URL
	if url == "" {
		return url
	}

	// Split into parts before and after @
	atParts := strings.Split(url, "@")
	if len(atParts) < 2 {
		return url // No @ found, return as-is
	}

	// Split the auth part (before @)
	authPart := atParts[0]
	colonParts := strings.Split(authPart, ":")

	// If we have a password part (format is redis://:password@host)
	if len(colonParts) >= 3 {
		colonParts[2] = "****" // Mask password
		authPart = strings.Join(colonParts, ":")
	}

	// Reconstruct the URL
	return authPart + "@" + strings.Join(atParts[1:], "@")
}
