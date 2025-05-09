# Linkeun Go API

This is a Go API project with animals endpoints.

## Features

- RESTful API for animals resource
- CRUD operations (Create, Read, Update, Delete)
- MySQL database integration with GORM
- Redis caching
- Pagination, sorting, and filtering support
- Docker and Kubernetes deployment
- Swagger documentation for API endpoints

## Prerequisites

- Go 1.24.2+
- Docker and Docker Compose
- MySQL 8.0+
- Redis 7.0+

## Project Structure

```
.
├── cmd/api                 # Application entry point
├── config                  # Configuration files
├── docker                  # Docker configuration files
├── internal
│   ├── controller          # HTTP handlers and controllers
│   ├── docs                # Swagger documentation
│   ├── model               # Data models
│   ├── repository          # Data access layer
│   └── service             # Business logic
├── k8s                     # Kubernetes deployment files
├── pkg
│   ├── config              # Configuration package
│   ├── database            # Database utilities
│   ├── middleware          # HTTP middleware
│   ├── pagination          # Pagination utilities
│   ├── response            # HTTP response utilities
│   └── seeder              # Database seeders
├── scripts                 # Helper scripts for development
├── test                    # Test files and utilities
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker build configuration
├── go.mod                  # Go module definition
└── README.md               # This file
```

## Getting Started

### Environment Variables

Create a `.env` file in the root directory with the following environment variables:

#### Application Settings
```
# Application environment (development, production, test)
APP_ENV=development

# Server configuration
PORT=8080                        # API server port
SERVER_READ_TIMEOUT=10s          # HTTP server read timeout
SERVER_WRITE_TIMEOUT=10s         # HTTP server write timeout
SERVER_SHUTDOWN_TIMEOUT=10s      # Graceful shutdown timeout
```

#### Database Configuration
```
# MySQL Database configuration
DB_USER=linkeun                  # Database username
DB_PASSWORD=root                 # Database password
DB_HOST=localhost                # Database host
DB_PORT=3306                     # Database port
DB_NAME=linkeun_go_api           # Database name
DB_PARAMS=charset=utf8mb4&parseTime=True&loc=Local  # Database connection parameters

# Database DSN (Data Source Name) - Alternative to individual connection parameters
# Example: DSN=linkeun:root@tcp(localhost:3306)/linkeun_go_api?charset=utf8mb4&parseTime=True&loc=Local
DSN=                             # Complete database connection string

# Database connection pool settings
DB_MAX_OPEN_CONNS=25             # Maximum number of open connections
DB_MAX_IDLE_CONNS=25             # Maximum number of idle connections
DB_CONN_MAX_LIFETIME=5m          # Maximum connection lifetime
```

#### Redis Configuration
```
# Redis configuration
REDIS_ENABLED=true               # Enable/disable Redis
REDIS_HOST=localhost             # Redis host
REDIS_PORT=6379                  # Redis port
REDIS_PASSWORD=your_password     # Redis password
REDIS_DB=0                       # Redis database number
REDIS_CACHE_TTL=15m              # Default cache expiration time
REDIS_PAGINATED_TTL=5m           # Expiration time for paginated results
REDIS_QUERY_CACHING=true         # Enable query caching
REDIS_KEY_PREFIX=linkeun_api:    # Key prefix for Redis
REDIS_POOL_SIZE=10               # Connection pool size
```

#### Logging Configuration
```
# Logging configuration
LOG_LEVEL=info                   # Log level (debug, info, warn, error)
LOG_FORMAT=json                  # Log format (json, text)
LOG_OUTPUT_PATH=stdout           # Log output path (stdout, file path)
```

#### Viewing Current Configuration

You can view your current environment configuration using the provided Makefile target:

```bash
# Show current environment variables
make env-info
```

This will display the actual values that will be used by the application, either from your .env file or the default values.

### Using as a Template Project

This project is designed to be used as a template for new Go API projects. To set up a new project based on this template:

1. Clone this repository
```bash
git clone https://github.com/linkeunid/go-api.git your-project-name
cd your-project-name
```

2. Use the setup commands to rename the module and reset Git:

```bash
# Basic setup - just rename the module
make setup module=github.com/yourusername/your-project

# Setup with Git remote
make setup-git module=github.com/yourusername/your-project remote=git@github.com:yourusername/your-project.git

# Full setup - rename module, reset Git repository, and set remote
make setup-full module=github.com/yourusername/your-project remote=git@github.com:yourusername/your-project.git
```

These commands will:
- Update the module name in go.mod
- Update all import paths in Go files throughout the project
- Optionally reset the Git repository (remove .git folder and initialize a new one)
- Optionally set a new Git remote URL

The setup process provides a confirmation prompt to ensure you understand the changes that will be made, as these operations cannot be undone.

3. After setup, update dependencies:
```bash
go mod tidy
```

For verbose output with detailed changes during setup, add the -v flag:
```bash
go run ./cmd/setup-project -module github.com/yourusername/your-project -v
```

Setup modes available:
- Basic mode: Only updates module name and import paths
- Git remote mode: Updates module name and sets Git remote URL
- Full setup: Updates module name, resets Git repository, and sets Git remote URL

Each mode is designed for different use cases:
- Basic mode is useful when you want to keep Git history but change the module name
- Git remote mode is useful when you want to keep Git history but change the remote repository
- Full setup is useful when starting a completely new project from the template

### Running with Docker

To run the application using Docker:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

### Running Locally

To run the application locally:

```bash
# Download dependencies
go mod download

# Build the application
go build -o bin/api ./cmd/api

# Run the application
./bin/api
```

### Database Migrations

The API uses a migration system to manage database schema changes. Migrations are stored in the `migrations` directory as SQL files.

```bash
# Run all pending migrations
make migrate

# Create a new empty migration (replace migration_name with a descriptive name)
make migrate-create name=add_new_field_to_animals

# Create a migration directly from a model (automatically generates SQL)
make migrate-from-model model=animal

# List available models for migrations
make migrate-list-models

# Roll back the last migration
make migrate-down

# Check current migration status
make migrate-status

# Reset all migrations (use with caution!)
make migrate-reset
```

Each migration consists of two files:
- `[timestamp]_[name].up.sql`: Contains SQL to apply the migration
- `[timestamp]_[name].down.sql`: Contains SQL to roll back the migration

Model-based migrations automatically generate SQL from your Go model structs using GORM tags to determine column types, constraints, and indexes.

### Database Seeding

The API includes a seeding system to populate the database with test data:

```bash
# Run all database seeders
make seed

# Run specific seeder
make seed-animal
make seed-flower

# Run seeders with custom record count
make seed-count count=500
```

### Model Map Management

The API uses a model map system to maintain a registry of available models for database operations. This system automatically detects models with a `TableName()` method.

```bash
# Update model map by adding new models
make update-model-map
# or using the shorthand
make um

# Clean model map by removing models that no longer exist
make clean-model-map
# or using the shorthand
make cm

# Sync model map by adding new models and removing deleted ones
make sync-model-map
# or using the shorthand
make sm
```

These commands help maintain the model map when adding or removing models:
- `update-model-map`: Adds new models but doesn't remove deleted ones
- `clean-model-map`: Removes models that no longer exist in the filesystem
- `sync-model-map`: Both adds new models and removes deleted ones

The model map system provides the following benefits:
- Automatic detection of model structs with TableName() methods
- Support for model registration without manual updates
- Consistent table naming through the application
- Simplified database operations through the cmd/db utility

For example, when you create a new model:

```go
// internal/model/flower.go
package model

type Flower struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:100;not null"`
    Color     string `gorm:"size:50;not null"`
    Fragrance string `gorm:"size:50"`
}

// TableName returns the table name for this model
func (Flower) TableName() string {
    return "flowers"
}
```

You can then run `make um` to add it to the model map automatically.

### Generating Swagger Documentation

The API includes Swagger documentation. To generate or update the documentation:

```bash
# Generate/update Swagger documentation
make swagger

# Run standalone Swagger UI server
make swagger-ui
```

The Swagger UI will be available at http://localhost:8090/swagger/

## API Endpoints

### Animals

| Method | Endpoint            | Description                 |
| ------ | ------------------- | --------------------------- |
| GET    | /api/v1/animals     | Get all animals (paginated) |
| GET    | /api/v1/animals/:id | Get a specific animal by ID |
| POST   | /api/v1/animals     | Create a new animal         |
| PUT    | /api/v1/animals/:id | Update an existing animal   |
| DELETE | /api/v1/animals/:id | Delete an animal            |

### Query Parameters

For the GET /api/v1/animals endpoint:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field (id, name, species, age, created_at, updated_at)
- `direction`: Sort direction (asc, desc)

### Example Requests

#### Get Animals (Paginated)

```
GET /api/v1/animals?page=1&limit=10&sort=name&direction=asc
```

#### Create Animal

```
POST /api/v1/animals
Content-Type: application/json

{
  "name": "Fluffy",
  "species": "Cat",
  "age": 3
}
```

#### Update Animal

```
PUT /api/v1/animals/1
Content-Type: application/json

{
  "name": "Fluffy",
  "species": "Cat",
  "age": 4
}
```

#### Delete Animal

```
DELETE /api/v1/animals/1
```

## Kubernetes Deployment

To deploy the application to Kubernetes:

```bash
# Apply deployment configuration
kubectl apply -f k8s/deployment.yaml
```

## Development Flow

The diagram below illustrates the development workflow for this project. It shows the four main phases: initial setup, database management, development cycle, and deployment. Each phase has specific commands (shown in parentheses) to help streamline the development process. This workflow is designed to make it easy to use this project as a template for new Go API projects while maintaining consistency in development practices.

<div align="center">
  <img src="docs/images/development-flow.svg" alt="Development Flow Diagram" width="800">
</div>

## License

This project is licensed under the GNU General Public License v2.0. See the LICENSE file for details.

## Redis Caching System

The API implements a Redis-based caching system to improve performance and reduce database load. The caching system is integrated with the repository layer and provides cache information in API responses.

### Cache Configuration

Cache behavior can be configured in the `.env` file:

```
# Redis configuration
REDIS_ENABLED=true               # Enable/disable Redis
REDIS_HOST=localhost             # Redis host
REDIS_PORT=6379                  # Redis port
REDIS_PASSWORD=your_password     # Redis password
REDIS_DB=0                       # Redis database number
REDIS_CACHE_TTL=15m              # Default cache expiration time
REDIS_PAGINATED_TTL=5m           # Expiration time for paginated results
REDIS_QUERY_CACHING=true         # Enable query caching
REDIS_KEY_PREFIX=linkeun_api:    # Key prefix for Redis
REDIS_POOL_SIZE=10               # Connection pool size
```

### Cache Information in API Responses

API responses include cache information in the `cacheInfo` property, which provides insights into the caching behavior:

```json
"cacheInfo": {
  "status": "hit",              // Cache status: hit, miss, or disabled
  "key": "query:animals:...",   // Cache key used for this query
  "enabled": true,              // Whether caching is enabled
  "ttl": "30m",                 // Time-to-live for this cache entry
  "useCount": 0                 // Usage statistics (not implemented yet)
}
```

#### Cache Status Types

- `hit`: The data was successfully retrieved from the cache. This indicates:
  - The requested data was previously stored in Redis
  - The TTL (expiration time) has not yet been reached
  - The system saved a database query, improving response time
  - The cache key correctly matched the query parameters

- `miss`: The data was not found in the cache and was fetched from the database. This can be due to:
  - The data has not been cached yet (first request)
  - The cache entry expired (TTL reached)
  - The cache was invalidated due to data updates

- `disabled`: Caching is disabled for this query or globally. This can happen when:
  - The `REDIS_ENABLED=false` setting is used in your configuration
  - The `REDIS_QUERY_CACHING=false` setting is used in your configuration
  - Redis connection is unavailable or failed
  - The specific model has caching disabled in its implementation
  - The query type is not suitable for caching (e.g., writes, complex joins)

The cache status in API responses helps troubleshoot and monitor the caching system's effectiveness. A high number of "miss" responses may indicate that the TTL is too short, or that data is being modified frequently.

#### Interpreting Cache Information

The `cacheInfo` object in API responses provides valuable insights for performance analysis and troubleshooting:

```json
"cacheInfo": {
  "status": "hit",              // Cache status (hit, miss, disabled)
  "key": "query:animals:...",   // The exact Redis key used
  "enabled": true,              // Whether caching is enabled for this request
  "ttl": "30m",                 // Time until this cache entry expires
  "useCount": 0                 // Usage statistics (planned feature)
}
```

This information can be used to:

1. **Performance Analysis**:
   - Compare response times between "hit" and "miss" statuses to quantify cache benefits
   - Monitor hit/miss ratios over time to optimize TTL settings
   - Identify queries that might benefit from longer or shorter TTL values

2. **Debugging**:
   - Verify that the expected cache keys are being generated
   - Confirm that cache invalidation is working correctly after updates
   - Check if caching is properly enabled/disabled based on configuration

3. **Cache Optimization**:
   - Identify frequently accessed data that should have longer TTL values
   - Detect patterns where certain queries almost always result in "miss"
   - Determine if the cache key generation strategy is effective

4. **Monitoring**:
   - Track cache hit ratios in your monitoring system
   - Set up alerts for unexpected cache disabled states
   - Measure cache efficiency across different API endpoints

The `key` field is particularly useful for Redis CLI debugging, as you can use the exact key to examine or manipulate the cache entry directly.

### Cache TTL (Time-To-Live)

Different types of queries have different cache TTL values:

- Single item queries (e.g., GetByID): Default is 15 minutes or as configured in `REDIS_CACHE_TTL`
- Paginated queries: Default is 5 minutes or as configured in `REDIS_PAGINATED_TTL`

All TTL values are dynamically loaded from environment variables - there are no hardcoded TTL values in the code. The system respects both the `REDIS_CACHE_TTL` and `REDIS_PAGINATED_TTL` settings from the `.env` file.

If `REDIS_PAGINATED_TTL` is not specified, the system will default to using 1/3 of the value specified in `REDIS_CACHE_TTL`.

This dual-TTL approach ensures that frequently changing data (like paginated lists) is refreshed more often than relatively static data, while giving you explicit control over both settings.

### Cache Invalidation

The system automatically invalidates cache entries when the corresponding data is modified:

- When an item is created, updated, or deleted, its individual cache entry is invalidated
- When the collection is modified, the collection cache is also invalidated

This ensures that API responses always reflect the most current data.

### Performance Benefits

The Redis caching system provides significant performance improvements:

- **Reduced Database Load**: By serving data from cache, the number of database queries is reduced
- **Faster Response Times**: Redis in-memory storage provides much faster data access than database queries
- **Improved Scalability**: The system can handle more concurrent users with the same database resources

### Cache Optimization Strategies

For optimal performance, consider these strategies:

1. **Adjust TTL Values**: Configure appropriate TTL values based on data change frequency:
   - Long TTL (hours) for rarely changing data
   - Short TTL (minutes) for frequently changing data
   - No caching for real-time critical data

2. **Cache Key Design**: The system uses query-based cache keys, which automatically:
   - Include relevant query parameters in the cache key
   - Create distinct keys for different queries on the same resource
   - Handle query variations like sorting and filtering

3. **Selective Caching**: Not all data should be cached:
   - Enable `CacheEnabled()` for frequently accessed models
   - Set larger TTL for static reference data
   - Use shorter TTL for user-specific or frequently updated data

### Best Practices for Cache Usage

To get the most out of the Redis caching system in this API, follow these best practices:

1. **Carefully Consider Cache Invalidation**:
   - Always invalidate cache entries when the underlying data changes
   - The repository layer already handles invalidation for CRUD operations
   - For custom repositories, ensure you implement similar invalidation logic

2. **Balance TTL Settings**:
   - Too short: High database load but more current data
   - Too long: Better performance but potentially stale data
   - Consider data change frequency when setting TTL values

3. **Monitor Cache Efficiency**:
   - Regularly check hit/miss ratios in API responses
   - Consider implementing metrics collection for cache performance
   - Use the `cacheInfo` data to tune your caching strategy

4. **Handle Cache Failures Gracefully**:
   - The system already falls back to database queries when Redis is unavailable
   - Consider adding circuit breaker patterns for repeated cache failures
   - Log and monitor Redis connection issues

5. **Cache Key Design**:
   - The current implementation uses query parameters in cache keys
   - For custom implementations, ensure keys are unique to the query
   - Avoid overly complex keys that might cause overhead

6. **Secure Your Redis Instance**:
   - Always use password authentication (set in `REDIS_PASSWORD`)
   - Consider network-level security (firewall, private subnet)
   - Periodically rotate Redis credentials

7. **Respect Data Privacy**:
   - Be cautious about caching sensitive or personal data
   - Consider implementing encryption for sensitive cached data
   - Ensure compliance with relevant privacy regulations

### Monitoring and Debugging

#### Examining Cache Behavior

1. Enable debug logging by setting `LOG_LEVEL=debug` in your `.env` file
2. Make API requests and observe the logs for cache-related messages:
   - "Cache hit" - Data was successfully retrieved from cache
   - "Cache miss" - Data was not found in cache and was fetched from database
   - "Setting Redis cache" - Data is being stored in cache
   - "Failed to store in Redis" - An error occurred when storing data

#### Monitoring Cache Efficiency

The `cacheInfo` field in API responses helps monitor cache efficiency:

- A high ratio of "hit" to "miss" indicates efficient caching
- Consistent "miss" responses for the same query may indicate cache invalidation issues
- Frequent cache invalidations might suggest that the TTL is too long for that data type
- "disabled" status when expecting caching indicates configuration issues

#### Common Cache Issues

- **"NOAUTH Authentication required"**: Check that the Redis password is correctly set in your `.env` file
- **"Connection refused"**: Ensure Redis server is running and accessible at the configured host/port
- **Empty cache results**: Verify that cache serialization/deserialization is working correctly

#### Redis CLI Commands for Debugging

Connect to your Redis instance using the Redis CLI:

```bash
# Connect to Redis server
redis-cli -h localhost -p 6379 -a your_password

# List all keys with the API prefix
keys linkeun_api:*

# Examine a specific cache entry
get linkeun_api:query:animals:1

# Delete a specific cache entry
del linkeun_api:query:animals:1

# Clear all cache entries for the API
flushdb
```

Note: Use `flushdb` with caution in shared Redis instances, as it clears all keys in the current database. 