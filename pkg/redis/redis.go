package redis

import (
	"context"
	"log"
	"time"

	"airline-booking/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(cfg *config.RedisConfig) *RedisClient {
	opt := &redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	//log.Println("Redis connected at:", cfg.Address)
	return &RedisClient{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) Close() {
	if err := r.Client.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	} else {
		log.Println("Redis connection closed")
	}
}
