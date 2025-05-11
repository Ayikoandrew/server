package database

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func NewRDB() *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisPort := os.Getenv("REDIS_PORT")

		if redisPassword != "" {
			url = fmt.Sprintf("redis://:%s@redis:%s/0", redisPassword, redisPort)
		} else {
			url = fmt.Sprintf("redis://redis:%s/0", redisPort)
		}
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opts)
	client.Options().MaxRetries = 5
	client.Options().MinRetryBackoff = 100 * time.Millisecond
	client.Options().MaxRetryBackoff = 2 * time.Millisecond
	return client
}

func GetRedisClient() *redis.Client {
	once.Do(func() {
		rdb = NewRDB()
	})
	return rdb
}

func InitRedis() {
	GetRedisClient()
}
