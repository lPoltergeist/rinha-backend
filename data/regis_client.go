package data

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Client  *redis.Client
	Context = context.Background()
)

func InitRedis() {

	Addr := os.Getenv("REDIS_HOST")
	if Addr == "" {
		Addr = "redis"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: Addr + ":6379",
	})
}
