package cache

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Version for cache invalidation when data model changes
const CurrentVersion = "v1"

// GenerateKey creates a structured, deterministic key for caching
func GenerateKey(entity string, params map[string]interface{}) string {
	parts := []string{CurrentVersion, entity}

	// Add sorted params for consistency
	var keys []string
	for k := range params {
		if params[k] != nil && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, params[k]))
	}

	return strings.Join(parts, ":")
}

// GenerateHashKey creates a hash-based key when parameters might be complex
func GenerateHashKey(entity string, params map[string]interface{}) string {
	// Convert params to string for hashing
	var paramValues []string
	for k, v := range params {
		if v != nil && v != "" {
			paramValues = append(paramValues, fmt.Sprintf("%s=%v", k, v))
		}
	}
	sort.Strings(paramValues)

	// Create hash
	h := sha256.New()
	h.Write([]byte(strings.Join(paramValues, ":")))

	return fmt.Sprintf("%s:%s:%x", CurrentVersion, entity, h.Sum(nil)[:8])
}

// Generic key generators for common patterns

// GenerateListKey creates a key for paginated entity lists
func GenerateListKey(entity string, page, limit int, sort, direction string) string {
	params := map[string]interface{}{
		"page":      page,
		"limit":     limit,
		"sort":      sort,
		"direction": direction,
	}
	return GenerateKey(entity+":list", params)
}

// GenerateItemKey creates a key for single entity items
func GenerateItemKey(entity string, id interface{}) string {
	return fmt.Sprintf("%s:%s:item:%v", CurrentVersion, entity, id)
}

// GenerateQueryKey creates a key for custom queries
func GenerateQueryKey(entity string, query string) string {
	// Use hash for query to avoid long keys
	h := sha256.New()
	h.Write([]byte(query))
	return fmt.Sprintf("%s:%s:query:%x", CurrentVersion, entity, h.Sum(nil)[:8])
}
