package database

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func NewRDB() *redis.Client {
	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opts)
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
