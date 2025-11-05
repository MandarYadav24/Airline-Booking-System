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

// NewRedisClient initializes and tests a Redis connection.
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

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", cfg.Address, err)
	}

	log.Printf("Connected to Redis at %s", cfg.Address)

	return &RedisClient{
		Client: client,
		Ctx:    context.Background(),
	}
}

// GetClient returns the underlying Redis client.
func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}

// Close safely closes the Redis client connection.
func (r *RedisClient) Close() {
	if r.Client == nil {
		return
	}
	if err := r.Client.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	} else {
		log.Println("Redis connection closed")
	}
}
