package database

import (
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

	// Create client either from URL or with explicit options
	var client *redis.Client

	// Try parsing URL first
	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Failed to parse Redis URL, falling back to direct connection", "error", err)

		// Fall back to direct connection options
		redisHost := os.Getenv("REDIS_HOST")
		if redisHost == "" {
			redisHost = "redis"
		}

		redisPort := os.Getenv("REDIS_PORT")
		if redisPort == "" {
			redisPort = "6379"
		}

		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	} else {
		client = redis.NewClient(opts)
	}

	// Add retry options
	client.Options().MaxRetries = 5
	client.Options().MinRetryBackoff = 100 * time.Millisecond
	client.Options().MaxRetryBackoff = 2 * time.Second

	return client
}

// Helper function to mask password in URL for logging
func maskPasswordInURL(url string) string {
	if strings.Contains(url, "@") {
		parts := strings.Split(url, "@")
		authParts := strings.Split(parts[0], ":")
		return authParts[0] + ":****@" + parts[1]
	}
	return url
}
