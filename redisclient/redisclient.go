package redisclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"listener-connection/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddress,
		Password:     config.RedisPassword,
		DB:           0,
		ReadTimeout:  30 * time.Second, // Read timeout
		WriteTimeout: 30 * time.Second, // Write timeout
		MaxRetries:   3,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: config.IsProduction,
		},
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
}
