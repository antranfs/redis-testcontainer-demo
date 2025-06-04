package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService represents a simple Redis service wrapper
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new Redis service instance
func NewRedisService(redisURL string) *RedisService {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)
	return &RedisService{client: client}
}

// Set stores a key-value pair in Redis with optional expiration
func (r *RedisService) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key from Redis
func (r *RedisService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Ping tests the Redis connection
func (r *RedisService) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}

// UserCache represents a simple user caching service
type UserCache struct {
	redis *RedisService
}

// NewUserCache creates a new user cache service
func NewUserCache(redisService *RedisService) *UserCache {
	return &UserCache{redis: redisService}
}

// CacheUser stores user data in cache
func (u *UserCache) CacheUser(ctx context.Context, userID string, userData string) error {
	key := fmt.Sprintf("user:%s", userID)
	return u.redis.Set(ctx, key, userData, 1*time.Hour)
}

// GetUser retrieves user data from cache
func (u *UserCache) GetUser(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("user:%s", userID)
	return u.redis.Get(ctx, key)
}

// InvalidateUser removes user data from cache
func (u *UserCache) InvalidateUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s", userID)
	return u.redis.Delete(ctx, key)
}

func main() {
	// This is a demo - in a real app, you'd get this from config
	redisURL := "redis://localhost:6379/0"

	redisService := NewRedisService(redisURL)
	defer redisService.Close()

	ctx := context.Background()

	// Test Redis connection
	if err := redisService.Ping(ctx); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return
	}

	// Demo usage
	userCache := NewUserCache(redisService)

	// Cache a user
	err := userCache.CacheUser(ctx, "123", `{"name": "John Doe", "email": "john@example.com"}`)
	if err != nil {
		log.Printf("Failed to cache user: %v", err)
		return
	}

	// Retrieve the user
	userData, err := userCache.GetUser(ctx, "123")
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}

	fmt.Printf("Cached user data: %s\n", userData)

	// Check if user exists
	exists, err := redisService.Exists(ctx, "user:123")
	if err != nil {
		log.Printf("Failed to check if user exists: %v", err)
		return
	}

	fmt.Printf("User exists in cache: %t\n", exists)
}
