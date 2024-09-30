package database

import (
	"authService/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRedis(cfg *config.Config) {
	redisHost := cfg.Redis.Host
	redisPort := cfg.Redis.Port

	RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	fmt.Println("Redis connection established successfully")
}
