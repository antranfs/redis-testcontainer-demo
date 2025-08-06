# Redis Testcontainer Demo - Codebase Overview

This is a **Go-based demonstration project** that showcases how to use Redis with Testcontainers for integration testing. The project implements a simple Redis service wrapper with user caching capabilities and comprehensive testing.

## 🏗️ Project Structure

The project consists of 4 key files:

- **`go.mod`** - Go module definition and dependencies
- **`main.go`** - Core application with Redis service implementation 
- **`redis_service_test.go`** - Integration tests using Testcontainers
- **`go.sum`** - Dependency checksums

## 🔧 Core Components

### 1. **RedisService** (`main.go:13-57`)
A wrapper around the Redis client providing essential operations:
- `NewRedisService()` - Creates Redis client from connection URL
- `Set()` - Stores key-value pairs with optional expiration
- `Get()` - Retrieves values by key
- `Delete()` - Removes keys from Redis
- `Exists()` - Checks key existence
- `Ping()` - Tests Redis connectivity
- `Close()` - Closes Redis connection

### 2. **UserCache** (`main.go:59-85`)
A higher-level service for user data caching:
- `CacheUser()` - Stores user data with 1-hour TTL
- `GetUser()` - Retrieves cached user data
- `InvalidateUser()` - Removes user from cache
- Uses prefixed keys (`user:{userID}`) for organization

### 3. **Demo Application** (`main.go:87-129`)
Simple demonstration that:
- Connects to Redis at `redis://localhost:6379/0`
- Caches sample user data
- Retrieves and displays cached data
- Checks key existence

## 🧪 Testing Strategy

**Integration Testing with Testcontainers** (`redis_service_test.go`):
- **`setupRedisContainer()`** - Spins up Redis 7 Alpine container for testing
- **`TestRedisService_BasicOperations()`** - Comprehensive test suite covering:
  - Basic set/get operations
  - Key expiration functionality 
  - Delete operations
  - Error handling for non-existent keys

## 📦 Key Dependencies

- **`github.com/redis/go-redis/v9`** - Redis client library
- **`github.com/testcontainers/testcontainers-go/modules/redis`** - Redis Testcontainer module
- **`github.com/stretchr/testify`** - Testing assertions and requirements

## 🎯 Use Cases

This codebase demonstrates:
- **Redis integration patterns** in Go applications
- **Testcontainer usage** for integration testing
- **Service layer abstraction** over Redis operations
- **User caching implementation** with TTL
- **Error handling** and connection management

## 🚀 Running the Code

The project is designed to:
1. **Run standalone**: `go run main.go` (requires Redis at localhost:6379)
2. **Run tests**: `go test ./...` (automatically spins up Redis container)

This is an excellent reference implementation for developers learning Redis integration with Go and modern testing practices using containers.