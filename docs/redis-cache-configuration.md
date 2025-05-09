# Redis Cache Configuration Guide

This document explains how to configure the Redis caching system for the Go API.

## TTL Configuration Options

The application now supports separate TTL (Time-To-Live) configurations for different types of cached data:

- `REDIS_CACHE_TTL`: The default TTL for most cached items (e.g., single entity lookups)
- `REDIS_PAGINATED_TTL`: A specific TTL for paginated results, which typically change more frequently

### Environment Variables

Add these to your `.env` file:

```
# Redis configuration
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6380
REDIS_DB=0
REDIS_CACHE_TTL=15m          # Default cache expiration time
REDIS_PAGINATED_TTL=5m       # Expiration time specifically for paginated results
REDIS_QUERY_CACHING=true
REDIS_KEY_PREFIX=linkeun_api:
REDIS_POOL_SIZE=10
```

### TTL Duration Format

TTL values use Go's duration format:
- `s`: seconds
- `m`: minutes
- `h`: hours

Examples:
- `15m`: 15 minutes
- `1h`: 1 hour
- `30s`: 30 seconds
- `1h30m`: 1 hour and 30 minutes

## Behavior

1. If `REDIS_PAGINATED_TTL` is set, that exact value will be used for paginated results.
2. If `REDIS_PAGINATED_TTL` is not set, the application will default to using 1/3 of the value specified in `REDIS_CACHE_TTL`.
3. The TTL values are displayed in the `cacheInfo` section of API responses.

## Recommendations

- **Regular items**: Use longer TTL (e.g., `15m` to `1h`) as they change less frequently
- **Paginated results**: Use shorter TTL (e.g., `3m` to `5m`) as they change more often and reflect the most recent data
- **Critical data**: For data that must be very current, use very short TTL (e.g., `30s`) or disable caching

## Debugging Cache Behavior

You can see the TTL values being used by setting `LOG_LEVEL=info` in your `.env` file and checking the application logs when it starts up. Look for log messages like:

```
Using configured paginated TTL paginatedTTL=5m
```

or 

```
Using calculated paginated TTL (1/3 of default TTL) defaultTTL=15m paginatedTTL=5m
``` 