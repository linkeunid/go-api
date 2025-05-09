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
REDIS_ENABLED=false              # Enable/disable Redis
REDIS_HOST=localhost             # Redis host
REDIS_PORT=6379                  # Redis port
REDIS_PASSWORD=                  # Redis password
REDIS_DB=0                       # Redis database number
REDIS_CACHE_TTL=15m              # Cache expiration time
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