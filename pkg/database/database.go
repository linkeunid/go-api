package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/linkeunid/go-api/pkg/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CacheStatus represents the cache status for a query
type CacheStatus string

const (
	// CacheHit indicates the data was served from cache
	CacheHit CacheStatus = "hit"
	// CacheMiss indicates the data was fetched from database
	CacheMiss CacheStatus = "miss"
	// CacheDisabled indicates caching is disabled for this query
	CacheDisabled CacheStatus = "disabled"
)

// ContextKey type for context keys
type ContextKey string

const (
	// ContextKeyCacheStatus is the context key for cache status
	ContextKeyCacheStatus ContextKey = "cache_status"
	// ContextKeyCacheKey is the context key for cache key
	ContextKeyCacheKey ContextKey = "cache_key"
)

// ErrRecordNotFound indicates a record was not found
var ErrRecordNotFound = gorm.ErrRecordNotFound

// Database defines the database interface
type Database interface {
	GetDB() *gorm.DB
	CachedFind(ctx context.Context, query *gorm.DB, dest interface{}) error
	GetCacheManager() CacheManager
	GetCacheStatus(ctx context.Context) (CacheStatus, string)
	Close() error
}

// CacheManager defines the cache manager interface
type CacheManager interface {
	GetCache() Cache
	GetConfig() *config.Config
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
}

// gormDatabase implements the Database interface using GORM
type gormDatabase struct {
	db           *gorm.DB
	cacheManager CacheManager
	logger       *zap.Logger
	config       *config.Config
	cacheContext context.Context
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
	// Create a new context to store cache status
	cacheCtx := context.WithValue(context.Background(), ContextKeyCacheStatus, CacheDisabled)
	cacheCtx = context.WithValue(cacheCtx, ContextKeyCacheKey, "")

	// Store the context in a package-level variable to make it accessible
	d.cacheContext = cacheCtx

	// If caching is not enabled, just perform the query and mark as disabled
	if d.cacheManager == nil || d.cacheManager.GetCache() == nil || !d.config.Redis.Enabled || !d.config.Redis.QueryCache {
		d.logger.Debug("Cache disabled")
		return query.Find(dest).Error
	}

	// Check if there's a custom cache key in the context
	var cacheKey string

	// First, check if a key is already set with our exact ContextKey
	if customKey := ctx.Value(ContextKeyCacheKey); customKey != nil {
		if key, ok := customKey.(string); ok && key != "" {
			cacheKey = key
			d.logger.Debug("Using cache key from context", zap.String("key", cacheKey), zap.String("source", "ContextKeyCacheKey"))
		}
	}

	// If no key found yet, try using repository's custom key
	if cacheKey == "" {
		if customKey := ctx.Value(ContextKey("customCacheKey")); customKey != nil {
			if key, ok := customKey.(string); ok && key != "" {
				cacheKey = key
				d.logger.Debug("Using custom cache key", zap.String("key", cacheKey), zap.String("source", "customCacheKey"))
			}
		}
	}

	// If still no key, try using string key name as fallback
	if cacheKey == "" {
		if customKey := ctx.Value("customCacheKey"); customKey != nil {
			if key, ok := customKey.(string); ok && key != "" {
				cacheKey = key
				d.logger.Debug("Using custom cache key from string key", zap.String("key", cacheKey), zap.String("source", "string"))
			}
		}
	}

	// If still no key, generate one from the query as last resort
	if cacheKey == "" {
		cacheKey = generateCacheKey(query)
		d.logger.Debug("Generated cache key from query", zap.String("key", cacheKey), zap.String("source", "generated"))
	}

	d.cacheContext = context.WithValue(d.cacheContext, ContextKeyCacheKey, cacheKey)

	// Try to get the item from cache first
	err := d.cacheManager.GetCache().Get(ctx, cacheKey, dest)
	if err == nil {
		// Cache hit
		d.cacheContext = context.WithValue(d.cacheContext, ContextKeyCacheStatus, CacheHit)
		d.logger.Debug("Cache hit", zap.String("key", cacheKey))
		return nil
	}

	// Cache miss
	d.cacheContext = context.WithValue(d.cacheContext, ContextKeyCacheStatus, CacheMiss)
	d.logger.Debug("Cache miss", zap.String("key", cacheKey))

	// Perform database query
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

// GetCacheStatus returns the cache status and key for the current request
func (d *gormDatabase) GetCacheStatus(ctx context.Context) (CacheStatus, string) {
	// Use the stored context if available
	if d.cacheContext != nil {
		status, ok := d.cacheContext.Value(ContextKeyCacheStatus).(CacheStatus)
		if ok {
			key, _ := d.cacheContext.Value(ContextKeyCacheKey).(string)
			return status, key
		}
	}

	// Fall back to passed context if needed
	status, ok := ctx.Value(ContextKeyCacheStatus).(CacheStatus)
	if !ok {
		return CacheDisabled, ""
	}

	key, ok := ctx.Value(ContextKeyCacheKey).(string)
	if !ok {
		key = ""
	}

	return status, key
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

	// Generate a more stable cache key
	// Extract useful information from query while avoiding memory addresses
	var tableName string
	var whereClauses []string
	var limitValue int
	var orderByValues []string
	var paginationInfo []string

	// Get table name
	tableName = stmt.Table

	// Extract WHERE conditions
	if whereClause, ok := stmt.Clauses["WHERE"]; ok && whereClause.Expression != nil {
		// Simplify WHERE conditions
		whereClauses = append(whereClauses, fmt.Sprintf("%v", whereClause.Expression))
	}

	// Extract LIMIT - handle safely without type assertions
	if limitClause, ok := stmt.Clauses["LIMIT"]; ok && limitClause.Expression != nil {
		// Just use the string representation instead of trying to cast
		if limit := fmt.Sprintf("%v", limitClause.Expression); limit != "" {
			limitValue = 1 // Just indicate that a limit exists
		}
	}

	// Extract ORDER BY - handle safely without type assertions
	if orderByClause, ok := stmt.Clauses["ORDER BY"]; ok && orderByClause.Expression != nil {
		// Just use the string representation
		orderStr := fmt.Sprintf("%v", orderByClause.Expression)
		if orderStr != "" {
			orderByValues = append(orderByValues, orderStr)
		}
	}

	// Include page and limit parameters if found in the request context
	if ctx := query.Statement.Context; ctx != nil {
		// Try multiple possible context keys for query parameters
		var queryParams map[string]string

		// Check both string and custom key types
		for _, keyName := range []string{"queryParams", "KeyQueryParams"} {
			if values := ctx.Value(keyName); values != nil {
				if params, ok := values.(map[string]string); ok {
					queryParams = params
					break
				}
			}
		}

		if queryParams != nil {
			if page, ok := queryParams["page"]; ok {
				paginationInfo = append(paginationInfo, fmt.Sprintf("page=%s", page))
			}
			if limit, ok := queryParams["limit"]; ok {
				paginationInfo = append(paginationInfo, fmt.Sprintf("limit=%s", limit))
			}
			if sort, ok := queryParams["sort"]; ok {
				paginationInfo = append(paginationInfo, fmt.Sprintf("sort=%s", sort))
			}
			if direction, ok := queryParams["direction"]; ok {
				paginationInfo = append(paginationInfo, fmt.Sprintf("direction=%s", direction))
			}
		}
	}

	// Build the key with stable components
	key := fmt.Sprintf("query:%s", tableName)

	if len(whereClauses) > 0 {
		key += fmt.Sprintf(":where(%s)", strings.Join(whereClauses, ","))
	}

	if limitValue > 0 {
		key += ":limit"
	}

	if len(orderByValues) > 0 {
		key += fmt.Sprintf(":orderby(%s)", strings.Join(orderByValues, ","))
	}

	// Add pagination info to the key
	if len(paginationInfo) > 0 {
		key += fmt.Sprintf(":pagination(%s)", strings.Join(paginationInfo, ","))
	}

	return key
}
