package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

// setupRedisContainer starts a Redis testcontainer and returns the service and cleanup function
func setupRedisContainer(t *testing.T) (*RedisService, func()) {
	ctx := context.Background()

	// Start Redis container
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err, "Failed to start Redis container")

	// Get connection string
	connectionString, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get Redis connection string")

	// Create Redis service
	redisService := NewRedisService(connectionString)

	// Test connection
	err = redisService.Ping(ctx)
	require.NoError(t, err, "Failed to ping Redis")

	// Cleanup function
	cleanup := func() {
		redisService.Close()
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	}

	return redisService, cleanup
}

func TestRedisService_BasicOperations(t *testing.T) {
	redisService, cleanup := setupRedisContainer(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test:key"
		value := "test-value"

		// Set value
		err := redisService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Get value
		retrievedValue, err := redisService.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, retrievedValue)
	})

	t.Run("Set with expiration", func(t *testing.T) {
		key := "test:expiring"
		value := "expiring-value"
		expiration := 100 * time.Millisecond

		// Set value with expiration
		err := redisService.Set(ctx, key, value, expiration)
		assert.NoError(t, err)

		// Immediately check if exists
		exists, err := redisService.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Check if expired
		exists, err = redisService.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "test:delete"
		value := "delete-me"

		// Set value
		err := redisService.Set(ctx, key, value, 0)
		assert.NoError(t, err)

		// Verify it exists
		exists, err := redisService.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Delete
		err = redisService.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify it's gone
		exists, err = redisService.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		_, err := redisService.Get(ctx, "non-existent-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis: nil")
	})
}
