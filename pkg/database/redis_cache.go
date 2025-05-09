package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/linkeunid/go-api/pkg/config"
	"go.uber.org/zap"
)

// RedisCacheManager implements the CacheManager interface using Redis
type RedisCacheManager struct {
	client *redis.Client
	logger *zap.Logger
	config *config.Config
}

// NewRedisCacheManager creates a new Redis cache manager
func NewRedisCacheManager(cfg *config.Config, logger *zap.Logger) (*RedisCacheManager, error) {
	// Configure Redis client
	redisOpt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	}

	// Create Redis client
	client := redis.NewClient(redisOpt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis",
		zap.String("host", cfg.Redis.Host),
		zap.Int("port", cfg.Redis.Port),
		zap.Int("db", cfg.Redis.DB),
	)

	return &RedisCacheManager{
		client: client,
		logger: logger,
		config: cfg,
	}, nil
}

// GetCache returns the Redis cache implementation
func (r *RedisCacheManager) GetCache() Cache {
	return r
}

// GetConfig returns the configuration used by this cache manager
func (r *RedisCacheManager) GetConfig() *config.Config {
	return r.config
}

// Get retrieves an item from cache
func (r *RedisCacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	// Add prefix to key
	prefixedKey := r.config.Redis.KeyPrefix + key

	r.logger.Debug("Getting from Redis cache", zap.String("key", prefixedKey))

	// Get value from Redis
	val, err := r.client.Get(ctx, prefixedKey).Result()
	if err != nil {
		if err == redis.Nil {
			r.logger.Debug("Cache miss - key not found", zap.String("key", prefixedKey))
			return fmt.Errorf("key not found: %s", key)
		}
		r.logger.Warn("Redis error", zap.String("key", prefixedKey), zap.Error(err))
		return err
	}

	// Unmarshal JSON data
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		r.logger.Warn("Failed to unmarshal data from Redis", zap.String("key", prefixedKey), zap.Error(err))
		return fmt.Errorf("failed to unmarshal data from Redis: %w", err)
	}

	r.logger.Debug("Cache hit", zap.String("key", prefixedKey))
	return nil
}

// Set stores an item in cache
func (r *RedisCacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Add prefix to key
	prefixedKey := r.config.Redis.KeyPrefix + key

	r.logger.Debug("Setting Redis cache", zap.String("key", prefixedKey), zap.Duration("ttl", expiration))

	// Marshal value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		r.logger.Warn("Failed to marshal data for Redis", zap.String("key", prefixedKey), zap.Error(err))
		return fmt.Errorf("failed to marshal data for Redis: %w", err)
	}

	// Store in Redis
	if err := r.client.Set(ctx, prefixedKey, jsonData, expiration).Err(); err != nil {
		r.logger.Warn("Failed to store in Redis", zap.String("key", prefixedKey), zap.Error(err))
		return fmt.Errorf("failed to store in Redis: %w", err)
	}

	r.logger.Debug("Successfully cached", zap.String("key", prefixedKey))
	return nil
}

// Delete removes an item from cache
func (r *RedisCacheManager) Delete(ctx context.Context, key string) error {
	// Add prefix to key
	prefixedKey := r.config.Redis.KeyPrefix + key

	// Check if the key contains a wildcard
	if strings.Contains(key, "*") {
		// Use pattern matching to delete multiple keys
		r.logger.Debug("Deleting keys with pattern", zap.String("pattern", prefixedKey))

		// Get matching keys
		keys, err := r.client.Keys(ctx, prefixedKey).Result()
		if err != nil {
			r.logger.Warn("Failed to find keys matching pattern", zap.String("pattern", prefixedKey), zap.Error(err))
			return fmt.Errorf("failed to find keys matching pattern: %w", err)
		}

		if len(keys) == 0 {
			// No keys to delete
			return nil
		}

		// Delete all matching keys
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			r.logger.Warn("Failed to delete keys with pattern", zap.String("pattern", prefixedKey), zap.Error(err))
			return fmt.Errorf("failed to delete from Redis: %w", err)
		}

		r.logger.Debug("Successfully deleted keys with pattern",
			zap.String("pattern", prefixedKey),
			zap.Int("count", len(keys)))
		return nil
	}

	// Delete single key from Redis
	if err := r.client.Del(ctx, prefixedKey).Err(); err != nil {
		return fmt.Errorf("failed to delete from Redis: %w", err)
	}

	return nil
}
