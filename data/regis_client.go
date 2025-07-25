package data

import (
	"context"
	"fmt"
	"log"
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

	_, err := Client.Ping(Context).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar no Redis: %v", err)
	}

	fmt.Println("Conectado ao Redis com sucesso!")
}
