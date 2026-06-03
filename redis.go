package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, cfg RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis successfully!")
	return rdb
}
