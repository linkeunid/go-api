package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/pkg/cache"
	"github.com/linkeunid/go-api/pkg/database"
	"github.com/linkeunid/go-api/pkg/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CacheInfo holds information about cache usage for a query
type CacheInfo struct {
	Status   database.CacheStatus `json:"status"`   // hit, miss, or disabled
	Key      string               `json:"key"`      // cache key used
	Enabled  bool                 `json:"enabled"`  // whether caching is enabled
	TTL      string               `json:"ttl"`      // time-to-live of the cache
	UseCount int                  `json:"useCount"` // number of times used (not implemented yet)
}

// AnimalResult wraps the animal data with cache information
type AnimalResult struct {
	Data      *model.Animal `json:"data"`
	CacheInfo *CacheInfo    `json:"cacheInfo,omitempty"`
}

// AnimalCollectionResult wraps the animal collection with cache information
type AnimalCollectionResult struct {
	Data       []model.Animal     `json:"data"`
	Pagination *pagination.Params `json:"pagination,omitempty"`
	CacheInfo  *CacheInfo         `json:"cacheInfo,omitempty"`
}

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// KeyCustomCacheKey is the context key for custom cache keys
	KeyCustomCacheKey ContextKey = "customCacheKey"
	// KeyQueryParams is the context key for query parameters
	KeyQueryParams ContextKey = "queryParams"
)

// AnimalRepository defines the interface for animal data access
type AnimalRepository interface {
	FindAll(ctx context.Context) (AnimalCollectionResult, error)
	FindAllPaginated(ctx context.Context, params pagination.Params) (AnimalCollectionResult, error)
	FindByID(ctx context.Context, id uint64) (AnimalResult, error)
	Create(ctx context.Context, animal *model.Animal) error
	Update(ctx context.Context, animal *model.Animal) error
	Delete(ctx context.Context, id uint64) error
}

// mysqlAnimalRepository implements AnimalRepository using MySQL with Redis cache
type mysqlAnimalRepository struct {
	db     database.Database
	logger *zap.Logger
	// Store TTL settings
	defaultTTL   string
	paginatedTTL string
}

// NewAnimalRepository creates a new animal repository
func NewAnimalRepository(db database.Database, logger *zap.Logger) AnimalRepository {
	// Get the default TTL from configuration or use sensible defaults
	// The actual TTL is applied in the CachedFind method
	defaultTTL := "30m"
	paginatedTTL := "5m"

	// If db has config, get TTL values from it
	if cacheManager := db.GetCacheManager(); cacheManager != nil {
		if redisMgr, ok := cacheManager.(*database.RedisCacheManager); ok && redisMgr.GetConfig() != nil {
			cfg := redisMgr.GetConfig()

			// Use the REDIS_CACHE_TTL from config (set to 15m in .env)
			defaultTTL = cfg.Redis.CacheTTL.String()

			// Use the REDIS_PAGINATED_TTL from config if defined, otherwise default to 1/3 of CacheTTL
			if cfg.Redis.PaginatedTTL > 0 {
				paginatedTTL = cfg.Redis.PaginatedTTL.String()
				logger.Info("Using configured paginated TTL",
					zap.String("paginatedTTL", paginatedTTL))
			} else {
				// Otherwise use a fraction of the default TTL (1/3)
				paginatedTTL = (cfg.Redis.CacheTTL / 3).String()
				logger.Info("Using calculated paginated TTL (1/3 of default TTL)",
					zap.String("defaultTTL", defaultTTL),
					zap.String("paginatedTTL", paginatedTTL))
			}
		}
	}

	return &mysqlAnimalRepository{
		db:           db,
		logger:       logger,
		defaultTTL:   defaultTTL,
		paginatedTTL: paginatedTTL,
	}
}

// createContextWithCacheKey creates a new context with a cache key
func (r *mysqlAnimalRepository) createContextWithCacheKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, KeyCustomCacheKey, key)
}

// createCacheInfo creates a CacheInfo struct from context and TTL
func (r *mysqlAnimalRepository) createCacheInfo(ctx context.Context, ttl string) *CacheInfo {
	status, key := r.db.GetCacheStatus(ctx)
	return &CacheInfo{
		Status:  status,
		Key:     key,
		Enabled: status != database.CacheDisabled,
		TTL:     ttl,
	}
}

// invalidateCache invalidates cache entries for an animal or collection
func (r *mysqlAnimalRepository) invalidateCache(ctx context.Context, itemID uint64, invalidateCollection bool) {
	cacheManager := r.db.GetCacheManager()
	if cacheManager == nil || cacheManager.GetCache() == nil {
		return
	}

	// Invalidate individual cache if itemID is provided
	if itemID > 0 {
		cacheKey := cache.GenerateItemKey("animals", itemID)
		if err := cacheManager.GetCache().Delete(ctx, cacheKey); err != nil {
			r.logger.Warn("Failed to invalidate animal cache", zap.Uint64("id", itemID), zap.Error(err))
		}
	}

	// Invalidate collection cache if requested
	if invalidateCollection {
		listKey := fmt.Sprintf("%s:animals:list", cache.CurrentVersion)
		if err := cacheManager.GetCache().Delete(ctx, listKey+"*"); err != nil {
			r.logger.Warn("Failed to invalidate animal collection cache", zap.Error(err))
		}
	}
}

// FindAll retrieves all animals with caching
func (r *mysqlAnimalRepository) FindAll(ctx context.Context) (AnimalCollectionResult, error) {
	var animals []model.Animal

	// Build the query
	query := r.db.GetDB().Order("created_at DESC")

	// Create a custom cache key
	cacheKey := cache.GenerateListKey("animals", 1, 0, "created_at", "desc")

	// Add the cache key to the context
	ctxWithKey := r.createContextWithCacheKey(ctx, cacheKey)

	// Use cached find with default TTL
	err := r.db.CachedFind(ctxWithKey, query, &animals)

	// Get cache status
	cacheInfo := r.createCacheInfo(ctxWithKey, r.defaultTTL)

	result := AnimalCollectionResult{
		Data:      animals,
		CacheInfo: cacheInfo,
	}

	if err != nil {
		r.logger.Error("Failed to retrieve animals", zap.Error(err))
		return result, err
	}

	return result, nil
}

// FindAllPaginated retrieves paginated animals
func (r *mysqlAnimalRepository) FindAllPaginated(ctx context.Context, params pagination.Params) (AnimalCollectionResult, error) {
	var animals []model.Animal
	result := AnimalCollectionResult{
		Pagination: &params,
	}

	// Get sort field and direction from query parameters
	sortField := "id"      // Default sort field
	sortDirection := "asc" // Default sort direction

	// Query values
	queryParams := make(map[string]string)
	values := ctx.Value(KeyQueryParams)
	if values != nil {
		if existingParams, ok := values.(map[string]string); ok {
			queryParams = existingParams
		}
	}

	// Make sure pagination parameters are included in the query context
	// This will ensure they're part of the cache key
	queryParams["page"] = fmt.Sprintf("%d", params.Page)
	queryParams["limit"] = fmt.Sprintf("%d", params.Limit)

	if field, exists := queryParams["sort"]; exists && field != "" {
		// Basic sanitization to prevent SQL injection
		allowedFields := map[string]bool{"id": true, "name": true, "species": true, "age": true, "created_at": true, "updated_at": true}
		if allowedFields[field] {
			sortField = field
		}
	}

	if dir, exists := queryParams["direction"]; exists {
		if dir == "desc" {
			sortDirection = "desc"
		}
	}

	// Generate a structured cache key using our key generator
	cacheKey := cache.GenerateListKey(
		"animals",
		params.Page,
		params.Limit,
		sortField,
		sortDirection,
	)

	// Create a new context with the structured cache key
	queryCtx := r.createContextWithCacheKey(ctx, cacheKey)

	// Build the base query with explicit model
	baseQuery := r.db.GetDB().Model(&model.Animal{})

	// Apply sorting
	orderClause := fmt.Sprintf("%s %s", sortField, sortDirection)
	baseQuery = baseQuery.Order(orderClause)

	// Count total rows
	var totalRows int64
	if err := baseQuery.Count(&totalRows).Error; err != nil {
		r.logger.Error("Failed to count animals", zap.Error(err))
		return result, err
	}

	// Calculate pagination metadata
	params.CalculatePages(totalRows)

	// Apply pagination limits and get results - Use CachedFind instead of direct Find
	paginatedQuery := baseQuery.Offset(params.GetOffset()).Limit(params.Limit)

	// Use cached find for paginated results with the enhanced context
	err := r.db.CachedFind(queryCtx, paginatedQuery, &animals)

	// Get cache status
	cacheInfo := r.createCacheInfo(queryCtx, r.paginatedTTL)

	result.Data = animals
	result.CacheInfo = cacheInfo

	if err != nil {
		r.logger.Error("Failed to retrieve paginated animals", zap.Error(err))
		return result, err
	}

	return result, nil
}

// FindByID retrieves an animal by ID with caching
func (r *mysqlAnimalRepository) FindByID(ctx context.Context, id uint64) (AnimalResult, error) {
	if id == 0 {
		return AnimalResult{}, errors.New("invalid ID")
	}

	var animal model.Animal
	result := AnimalResult{}

	// Build the query
	query := r.db.GetDB().Where("id = ?", id)

	// Generate a structured cache key for the item
	cacheKey := cache.GenerateItemKey("animals", id)

	// Add the cache key to the context
	ctxWithKey := r.createContextWithCacheKey(ctx, cacheKey)

	// Use cached find
	err := r.db.CachedFind(ctxWithKey, query, &animal)

	// Get cache status
	cacheInfo := r.createCacheInfo(ctxWithKey, r.defaultTTL)
	result.CacheInfo = cacheInfo

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, nil // Return empty result for not found
		}
		r.logger.Error("Failed to retrieve animal by ID", zap.Uint64("id", id), zap.Error(err))
		return result, err
	}

	// Check if the record was actually found (GORM might not return ErrRecordNotFound)
	if animal.ID == 0 {
		return result, nil // Return empty result for not found
	}

	result.Data = &animal
	return result, nil
}

// Create saves a new animal
func (r *mysqlAnimalRepository) Create(ctx context.Context, animal *model.Animal) error {
	// Create the record (ID will be auto-generated by the database)
	if err := r.db.GetDB().Create(animal).Error; err != nil {
		r.logger.Error("Failed to create animal", zap.Error(err))
		return err
	}

	// Invalidate the collection cache
	r.invalidateCache(ctx, 0, true)

	return nil
}

// Update updates an existing animal
func (r *mysqlAnimalRepository) Update(ctx context.Context, animal *model.Animal) error {
	if animal.ID == 0 {
		return errors.New("invalid ID")
	}

	if err := r.db.GetDB().Save(animal).Error; err != nil {
		r.logger.Error("Failed to update animal", zap.Uint64("id", animal.ID), zap.Error(err))
		return err
	}

	// Invalidate both individual and collection caches
	r.invalidateCache(ctx, animal.ID, true)

	return nil
}

// Delete removes an animal
func (r *mysqlAnimalRepository) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	if err := r.db.GetDB().Delete(&model.Animal{}, id).Error; err != nil {
		r.logger.Error("Failed to delete animal", zap.Uint64("id", id), zap.Error(err))
		return err
	}

	// Invalidate both individual and collection caches
	r.invalidateCache(ctx, id, true)

	return nil
}
