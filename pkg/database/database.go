package database

import (
	"context"
	"fmt"
	"time"

	"github.com/linkeunid/go-api/pkg/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrRecordNotFound indicates a record was not found
var ErrRecordNotFound = gorm.ErrRecordNotFound

// Database defines the database interface
type Database interface {
	GetDB() *gorm.DB
	CachedFind(ctx context.Context, query *gorm.DB, dest interface{}) error
	GetCacheManager() CacheManager
	Close() error
}

// CacheManager defines the cache manager interface
type CacheManager interface {
	GetCache() Cache
}

// Cache defines the cache interface
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

// Cacheable is the interface that models must implement to be cacheable
type Cacheable interface {
	CacheEnabled() bool
	CacheTTL() time.Duration
	CacheKey() string
}

// gormDatabase implements the Database interface using GORM
type gormDatabase struct {
	db           *gorm.DB
	cacheManager CacheManager
	logger       *zap.Logger
	config       *config.Config
}

// NewDatabase creates a new database instance
func NewDatabase(cfg *config.Config, logger *zap.Logger, db *gorm.DB, cacheManager CacheManager) Database {
	return &gormDatabase{
		db:           db,
		cacheManager: cacheManager,
		logger:       logger,
		config:       cfg,
	}
}

// GetDB returns the GORM DB instance
func (d *gormDatabase) GetDB() *gorm.DB {
	return d.db
}

// CachedFind performs a find operation with caching
func (d *gormDatabase) CachedFind(ctx context.Context, query *gorm.DB, dest interface{}) error {
	// If caching is not enabled, just perform the query
	if d.cacheManager == nil || d.cacheManager.GetCache() == nil || !d.config.Redis.Enabled || !d.config.Redis.QueryCache {
		return query.Find(dest).Error
	}

	// Try to get the item from cache first
	cacheKey := generateCacheKey(query)
	err := d.cacheManager.GetCache().Get(ctx, cacheKey, dest)
	if err == nil {
		return nil // Cache hit, return cached data
	}

	// Cache miss, perform database query
	if err := query.Find(dest).Error; err != nil {
		return err
	}

	// Store result in cache
	cacheTTL := d.config.Redis.CacheTTL
	if cacheable, ok := dest.(Cacheable); ok && cacheable.CacheEnabled() {
		cacheTTL = cacheable.CacheTTL()
	}

	if err := d.cacheManager.GetCache().Set(ctx, cacheKey, dest, cacheTTL); err != nil {
		d.logger.Warn("Failed to cache query result", zap.String("key", cacheKey), zap.Error(err))
	}

	return nil
}

// GetCacheManager returns the cache manager
func (d *gormDatabase) GetCacheManager() CacheManager {
	return d.cacheManager
}

// Close closes the database connection
func (d *gormDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	return sqlDB.Close()
}

// generateCacheKey generates a cache key for a GORM query
func generateCacheKey(query *gorm.DB) string {
	stmt := query.Statement
	if stmt == nil {
		return "query:unknown"
	}

	// Generate a key based on the table name and where conditions
	return fmt.Sprintf("query:%s:%v", stmt.Table, stmt.Clauses)
}
