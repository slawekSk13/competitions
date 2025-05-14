package models

import (
	"competition-app/config"
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis client
func InitRedis(cfg *config.Config) error {
	log.Printf("Connecting to Redis at %s:%d...", cfg.RedisHost, cfg.RedisPort)
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr: cfg.GetRedisConnString(),
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	
	log.Println("Redis connection established")
	return nil
}

// CloseRedis closes the Redis client connection
func CloseRedis() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}

// SetCache stores a value in the cache with an expiration time
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetCache retrieves a value from the cache
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// DeleteCache removes a value from the cache
func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// GetContext returns the context used by Redis
func GetContext() context.Context {
	return ctx
}
