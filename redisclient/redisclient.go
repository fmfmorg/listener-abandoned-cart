package redisclient

import (
	"context"
	"listener-connection/config"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,  // Replace with your Redis address
		Password: config.RedisPassword, // No password
		DB:       0,                    // Use default DB
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}
