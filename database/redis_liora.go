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
	if url == "" {

		host := getEnv("REDIS_HOST", "redis")
		port := getEnv("REDIS_PORT", "6379")
		password := os.Getenv("REDIS_PASSWORD")
		url = fmt.Sprintf("redis://:%s@%s:%s/0", password, host, port)
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Invalid Redis URL", "error", err)
		panic(fmt.Sprintf("Invalid Redis URL: %v", err))
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		maskedURL := strings.ReplaceAll(url, os.Getenv("REDIS_PASSWORD"), "****")
		slog.Error("Failed to connect to Redis",
			"error", err,
			"url", maskedURL)
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	slog.Info("Successfully connected to Redis")
	return client
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
