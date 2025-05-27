# Linkeun Go API

A comprehensive Go API project with RESTful endpoints, JWT authentication, caching, and database integration.

<details>
<summary>📋 Table of Contents</summary>

- [Linkeun Go API](#linkeun-go-api)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Environment Setup](#environment-setup)
    - [Running the API](#running-the-api)
      - [With Docker](#with-docker)
      - [Locally](#locally)
    - [Command Aliases](#command-aliases)
      - [Available Aliases](#available-aliases)
      - [Commonly Used Aliases](#commonly-used-aliases)
      - [Docker Aliases](#docker-aliases)
      - [Database Aliases](#database-aliases)
      - [Project Setup Aliases](#project-setup-aliases)
      - [Examples](#examples)
    - [API Endpoints](#api-endpoints)
      - [Animals Resource](#animals-resource)
      - [Query Parameters](#query-parameters)
  - [Development Flow Diagram](#development-flow-diagram)
  - [Project Structure](#project-structure)
  - [Authentication](#authentication)
    - [JWT Overview](#jwt-overview)
    - [Token Generation](#token-generation)
    - [Claims Structure](#claims-structure)
    - [Configuration Options](#configuration-options)
      - [Understanding JWT\_ALLOWED\_ISSUERS](#understanding-jwt_allowed_issuers)
    - [Implementation Details](#implementation-details)
      - [Using Authentication](#using-authentication)
      - [Available API Endpoints](#available-api-endpoints)
      - [Role-Based Access Control](#role-based-access-control)
      - [Error Handling](#error-handling)
    - [Security Best Practices](#security-best-practices)
  - [Caching System](#caching-system)
    - [Redis Configuration](#redis-configuration)
    - [Caching Features](#caching-features)
      - [Cache Information in Responses](#cache-information-in-responses)
      - [Pagination with Caching](#pagination-with-caching)
      - [Cache TTL Strategy](#cache-ttl-strategy)
      - [Cache Invalidation](#cache-invalidation)
    - [Caching Best Practices](#caching-best-practices)
  - [Logging System](#logging-system)
    - [Logging Configuration](#logging-configuration)
      - [Understanding LOG\_LEVEL](#understanding-log_level)
    - [Log Output Options](#log-output-options)
    - [Log Rotation Types](#log-rotation-types)
      - [Daily Rotation (Default)](#daily-rotation-default)
      - [Size-based Rotation](#size-based-rotation)
      - [Retention Settings](#retention-settings)
    - [Testing Log Rotation](#testing-log-rotation)
      - [Testing Size-based Rotation (Default)](#testing-size-based-rotation-default)
      - [Test Script Features](#test-script-features)
      - [Usage Examples](#usage-examples)
    - [Logging Best Practices](#logging-best-practices)
    - [Migration from Size-based to Daily Rotation](#migration-from-size-based-to-daily-rotation)
  - [Database Operations](#database-operations)
    - [Migrations](#migrations)
    - [Seeding](#seeding)
    - [Seeder Management](#seeder-management)
    - [Model Management](#model-management)
  - [Development](#development)
    - [Swagger Documentation](#swagger-documentation)
    - [Development Workflow](#development-workflow)
    - [Using as a Template](#using-as-a-template)
      - [What the Setup Tool Updates](#what-the-setup-tool-updates)
      - [Example Transformation](#example-transformation)
  - [Deployment](#deployment)
    - [Docker](#docker)
    - [Kubernetes](#kubernetes)
      - [Minikube Deployment](#minikube-deployment)
        - [Prerequisites](#prerequisites-1)
        - [Setup Minikube](#setup-minikube)
        - [Automated Deployment](#automated-deployment)
        - [Manual Deployment](#manual-deployment)
        - [Accessing the Application](#accessing-the-application)
        - [Kubernetes Resources](#kubernetes-resources)
        - [Monitoring and Debugging](#monitoring-and-debugging)
      - [Production Deployment Considerations](#production-deployment-considerations)
  - [API Testing](#api-testing)
    - [Running the Tests](#running-the-tests)
    - [Test Coverage](#test-coverage)
    - [Testing Strategy](#testing-strategy)
    - [Extending Tests](#extending-tests)
  - [License](#license)

</details>

## Features

- RESTful API endpoints with CRUD operations
- JWT authentication with role-based access control
- MySQL database integration with GORM ORM
- Redis-based caching system for performance optimization
- Pagination, sorting, and filtering support
- Docker and Kubernetes deployment configurations
- API documentation with Swagger
- Database migrations and seeding with automatic seeder registration
- Environment-specific configurations
- Comprehensive error handling
- Advanced logging system with daily/size-based rotation and configurable outputs

## Getting Started

### Prerequisites

- Go 1.24.2+
- Docker and Docker Compose (for containerized deployment)
- MySQL 8.0+
- Redis 7.0+

### Environment Setup

Create a `.env` file in the root directory with the following variables:

```bash
# Application environment
APP_ENV=development              # Options: development, test, production

# Server configuration
PORT=8080                        
SERVER_READ_TIMEOUT=10s          
SERVER_WRITE_TIMEOUT=10s         
SERVER_SHUTDOWN_TIMEOUT=10s      

# Logging configuration
LOG_LEVEL=info                  # Options: debug, info, warn, error
LOG_FORMAT=json                 # Options: json, console
LOG_OUTPUT_PATH=stdout          # Options: stdout, stderr
LOG_ROTATION_TYPE=daily         # Options: daily, size (default: daily)
LOG_FILE_PATH=./logs/app.log    # Path to log file (empty = disable file logging)
LOG_FILE_MAX_SIZE=100           # Maximum size of log files in megabytes before rotation
LOG_FILE_MAX_BACKUPS=3          # Maximum number of old log files to retain
LOG_FILE_MAX_AGE=28             # Maximum number of days to retain old log files
LOG_FILE_COMPRESS=true          # Whether to compress rotated log files

# Database configuration
DB_USER=linkeun                  
DB_PASSWORD=root                 
DB_HOST=localhost                
DB_PORT=3306                     
DB_NAME=linkeun_go_api           
DB_PARAMS=charset=utf8mb4&parseTime=True&loc=Local

# Redis configuration
REDIS_ENABLED=true               
REDIS_HOST=localhost             
REDIS_PORT=6379                  
REDIS_PASSWORD=your_password     

# Authentication
AUTH_ENABLED=true                
JWT_SECRET=your-secret-key       
JWT_EXPIRATION=24h               
JWT_ALLOWED_ISSUERS=linkeun-go-api
```

View current environment settings:
```bash
make env-info
```

### Running the API

#### With Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

#### Locally

```bash
# Download dependencies
go mod download

# Build the application
go build -o bin/api ./cmd/api

# Run the application
./bin/api
```

<details>
<summary>⚡ Command Aliases</summary>

### Command Aliases

The project includes numerous command aliases to make development more efficient. These aliases are shortcuts for commonly used commands.

#### Available Aliases

View all available command aliases:

```bash
make help
```

#### Commonly Used Aliases

| Alias | Full Command           | Description                            |
| ----- | ---------------------- | -------------------------------------- |
| `d`   | `dev`                  | Run development server with hot reload |
| `r`   | `run`                  | Run the application                    |
| `s`   | `swagger`              | Generate Swagger documentation         |
| `su`  | `swagger-ui`           | Run Swagger UI server                  |
| `t`   | `test`                 | Run tests                              |
| `l`   | `lint`                 | Lint code                              |
| `fr`  | `flush-redis`          | Flush Redis cache                      |
| `gt`  | `generate-token`       | Generate JWT token                     |
| `gtu` | `generate-token-user`  | Generate JWT token for specific user   |
| `gta` | `generate-token-admin` | Generate admin JWT token               |

#### Docker Aliases

| Alias   | Full Command  | Description                         |
| ------- | ------------- | ----------------------------------- |
| `dup`   | `docker-up`   | Start all containers                |
| `ddown` | `docker-down` | Stop all containers                 |
| `ddb`   | `docker-db`   | Start only database containers      |
| `fps`   | `fancy-ps`    | Show fancy container status details |

#### Database Aliases

| Alias | Full Command             | Description                        |
| ----- | ------------------------ | ---------------------------------- |
| `sd`  | `seed`                   | Run all database seeders           |
| `tr`  | `truncate`               | Truncate specific table            |
| `mam` | `migrate-all-models`     | Create migrations from all models  |
| `um`  | `update-model-map`       | Update model map for database      |
| `sm`  | `sync-model-map`         | Sync model map                     |
| `smm` | `sync-model-map`         | Sync model map (alternative alias) |
| `usr` | `update-seeder-registry` | Update seeder registry             |
| `csr` | `clean-seeder-registry`  | Clean seeder registry              |
| `ss`  | `sync-seeder`            | Sync seeder registry               |

#### Project Setup Aliases

| Alias     | Full Command | Description                        |
| --------- | ------------ | ---------------------------------- |
| `setup-s` | `setup`      | Setup project with new module name |
| `setup-g` | `setup-git`  | Setup project with git remote      |
| `setup-f` | `setup-full` | Full setup with new git repo       |

#### Examples

Start the development server:
```bash
make d
```

Run tests:
```bash
make t
```

Generate and view Swagger documentation:
```bash
make s
make su
```

Work with Docker containers:
```bash
make dup    # Start all containers
make fps    # Check container status
make ddown  # Stop all containers
```

Generate a JWT token for testing:
```bash
# Generate a token with default values
make gt
# or
make generate-token

# Generate an admin token
make gta
# or
make generate-token-admin

# Generate a token for a specific user ID
make gtu id=123
# or
make generate-token-user id=123

# Force token generation in any environment (for emergencies)
make gtf
# or
make generate-token-force
```

Work with database migrations:
```bash
make mam    # Create migrations for all models
make um     # Update model map
make sm     # Sync model map
```

Work with database seeders:
```bash
make ss     # Sync seeder registry (add new, remove deleted)
make usr    # Update seeder registry (add new seeders only)
make csr    # Clean seeder registry (remove deleted seeders only)
```

</details>

### API Endpoints

#### Animals Resource

| Method | Endpoint            | Description                 |
| ------ | ------------------- | --------------------------- |
| GET    | /api/v1/animals     | Get all animals (paginated) |
| GET    | /api/v1/animals/:id | Get a specific animal by ID |
| POST   | /api/v1/animals     | Create a new animal         |
| PUT    | /api/v1/animals/:id | Update an existing animal   |
| DELETE | /api/v1/animals/:id | Delete an animal            |

#### Query Parameters

For paginated endpoints:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field (e.g., id, name, created_at)
- `direction`: Sort direction (asc, desc)

## Development Flow Diagram

The following diagram illustrates the development workflow from initial setup through to deployment, highlighting the key commands and their aliases used at each stage:

```mermaid
graph TD
    subgraph "Setup"
        A[Clone Repository] --> B[Configure Environment]
        B --> C[Initialize Project with Aliases]
        C -->|make init / make i| D[Generate Swagger Docs]
    end
    
    subgraph "Database Setup"
        D --> E[Run Migrations]
        E -->|make migrate| E2[Create Model Migrations]
        E2 -->|make mam / make migrate-all-models| F[Seed Database]
        F -->|make seed / make sd| G[Ready for Development]
    end
    
    subgraph "Development Cycle"
        G --> H[Code Changes]
        H --> I[Run Development Server]
        I -->|make dev / make d| J[Test API]
        J --> K[Generate & View Documentation]
        K -->|make swagger / make s| K2[View Swagger UI]
        K2 -->|make swagger-ui / make su| H
    end
    
    subgraph "Testing"
        H --> L1[Run Tests]
        L1 -->|make test / make t| L2[Format Code]
        L2 -->|make fmt| L3[Lint Code]
        L3 -->|make lint / make l| H
    end
    
    subgraph "Deployment"
        H --> M1[Build]
        M1 -->|make build| M2[Deployment Options]
        M2 --> N[Docker Container]
        N -->|make docker-up / make dup| O[Production]
        M2 --> P[Kubernetes]
        P -->|k8s/deploy-minikube.sh| O
    end
    
    subgraph "Helper Commands"
        Q1[Generate JWT Token] -->|make gt, make gtu, make gta| J
        Q2[Database Operations] -->|make um, make sm, make smm| J
        Q3[Migrate All Models] -->|make mam| J
        Q4[Seeder Management] -->|make ss, make usr, make csr| J
        Q5[Monitor Containers] -->|make fps, make docker-logs / make dlogs| J
        Q6[Flush Cache] -->|make flush-redis / make fr| J
    end
```

<details>
<summary>🏗️ Project Structure</summary>

## Project Structure

```
.
├── cmd/                      # Command-line applications
│   ├── api/                  # Main API application
│   ├── seeder-mapper/        # Automatic seeder registration utility
│   └── token-generator/      # JWT token generation utility
├── internal/                 # Private application code
│   ├── bootstrap/            # Application bootstrapping
│   ├── controller/           # HTTP request handlers
│   ├── model/                # Data models
│   ├── repository/           # Data access layer
│   └── service/              # Business logic
├── pkg/                      # Public packages
│   ├── auth/                 # Authentication services
│   ├── config/               # Configuration utilities
│   ├── middleware/           # HTTP middleware
│   ├── response/             # HTTP response utilities
│   └── ...                   # Other utility packages
├── scripts/                  # Helper scripts
├── migrations/               # Database migrations
├── docker/                   # Docker configurations
├── k8s/                      # Kubernetes manifests
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker build configuration
├── go.mod                    # Go module definition
└── README.md                 # Documentation
```

</details>

<details>
<summary>🔐 Authentication</summary>

## Authentication

### JWT Overview

The API uses JSON Web Tokens (JWT) for authentication with the following features:

- Token-based authentication with secure validation
- Role-based access control for protected endpoints
- Configurable token expiration and issuer validation
- Environment restrictions for token generation
- Comprehensive claim validation

### Token Generation

You can generate test tokens using the provided utility. **Note: Token generation is only available in development and test environments by default.**

```bash
# Generate a token with default values
make gt
# or
make generate-token

# Generate an admin token
make gta
# or
make generate-token-admin

# Generate a token for a specific user ID
make gtu id=123
# or
make generate-token-user id=123

# Force token generation in any environment (for emergencies)
make gtf
# or
make generate-token-force
```

### Claims Structure

JWT tokens use the following claims structure:

```json
{
  // Custom claims
  "username": "johndoe",        // Username (string)
  "role": "admin",              // User role (string)
  "email": "john@example.com",  // User email (string)

  // Standard JWT claims
  "iss": "linkeun-go-api",      // Issuer
  "sub": "123",                 // Subject (user ID as string)
  "exp": 1673667272,            // Expiration Time (Unix timestamp)
  "iat": 1673580872             // Issued At (Unix timestamp)
}
```

Accessing claims in your code:

```go
// From request context after authentication
userID := r.Context().Value(middleware.KeyUserID).(uint64)  // Parsed from 'sub' claim
username := r.Context().Value(middleware.KeyUsername).(string)
role := r.Context().Value(middleware.KeyUserRole).(string)
email := r.Context().Value(middleware.KeyUserEmail).(string)
```

### Configuration Options

Configure authentication via environment variables:

```
AUTH_ENABLED=true                # Enable/disable authentication
JWT_SECRET=your-secret-key       # Secret key for JWT signing
JWT_EXPIRATION=24h               # Token expiration time
JWT_ALLOWED_ISSUERS=linkeun-go-api,other-trusted-issuer
```

#### Understanding JWT_ALLOWED_ISSUERS

This variable contains a comma-separated list of trusted token issuers:

- **Purpose**: Controls which systems can issue accepted tokens
- **Format**: Comma-separated names (no spaces)
- **Default**: Only accepts tokens from `linkeun-go-api`
- **Examples**:
  - Single API: `JWT_ALLOWED_ISSUERS=linkeun-go-api`
  - Multiple services: `JWT_ALLOWED_ISSUERS=linkeun-go-api,auth-service,admin-portal`

### Implementation Details

#### Using Authentication

Protected endpoints require a JWT token in the Authorization header:

```bash
# Access protected endpoint
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:8080/api/v1/protected/

# Access admin-only endpoint
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:8080/api/v1/protected/admin/
```

#### Available API Endpoints

| Endpoint                     | Auth Required | Role Required | Description                       |
| ---------------------------- | ------------- | ------------- | --------------------------------- |
| GET /health                  | No            | None          | Health check endpoint             |
| GET /swagger/                | No            | None          | Swagger UI (dev mode only)        |
| GET /api/v1/public/          | No            | None          | Public API endpoint               |
| GET /api/v1/protected/       | Yes           | Any           | Protected endpoint with user info |
| GET /api/v1/protected/admin/ | Yes           | Admin         | Admin-only protected endpoint     |
| GET /api/v1/animals          | No*           | None          | List all animals                  |
| GET /api/v1/animals/:id      | No*           | None          | Get animal by ID                  |
| POST /api/v1/animals         | No*           | None          | Create a new animal               |
| PUT /api/v1/animals/:id      | No*           | None          | Update an animal                  |
| DELETE /api/v1/animals/:id   | No*           | None          | Delete an animal                  |

*Note: Animal endpoints may require authentication depending on your configuration.

#### Role-Based Access Control

Protect routes with role requirements:

```go
// Authenticate any user
r.Use(authMiddleware.Authenticate)

// Require specific role(s)
r.Use(authMiddleware.RequireRole("admin"))
r.Use(authMiddleware.RequireRole("admin", "manager"))
```

#### Error Handling

The middleware provides specific error messages:

- **Missing Token**: "Authorization header is required"
- **Invalid Format**: "Invalid token format, expected 'Bearer <token>'"
- **Expired Token**: "Token has expired"
- **Invalid Token**: "Invalid token"
- **Invalid Issuer**: "Invalid token issuer"

### Security Best Practices

This implementation follows these security practices:

1. **Secret Management**:
   - Environment variables for secrets
   - Different secrets per environment
   - Regular secret rotation

2. **Token Validation**:
   - Signature verification
   - Expiration validation
   - Issuer validation
   - Environment restrictions

3. **Claims Best Practices**:
   - Standard JWT claims (iss, sub, exp, iat)
   - Minimal custom claims
   - No sensitive data in tokens

4. **Security Headers**:
   - Authorization header (not cookies)
   - Bearer authentication scheme
   - Clear error messages without exposing internals

For more details, see the [OWASP JWT Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html).

</details>

<details>
<summary>🚀 Caching System</summary>

## Caching System

The API implements a Redis-based caching system to improve performance and reduce database load.

### Redis Configuration

Configure caching in your `.env` file:

```
REDIS_ENABLED=true               # Enable/disable Redis
REDIS_HOST=localhost             # Redis host
REDIS_PORT=6379                  # Redis port
REDIS_PASSWORD=your_password     # Redis password
REDIS_CACHE_TTL=15m              # Default cache expiration
REDIS_PAGINATED_TTL=5m           # Paginated results expiration
REDIS_QUERY_CACHING=true         # Enable query caching
REDIS_KEY_PREFIX=linkeun_api:    # Key prefix
```

### Caching Features

#### Cache Information in Responses

API responses include cache details:

```json
"cacheInfo": {
  "status": "hit",              // hit, miss, or disabled
  "key": "query:animals:...",   // Cache key
  "enabled": true,              // Caching status
  "ttl": "30m",                 // Time-to-live
  "useCount": 0                 // Usage statistics
}
```

#### Pagination with Caching

The system ensures proper caching for paginated results:

- Each page has its own cache entry with unique keys
- Pagination parameters are included in cache keys
- Cache invalidation works across all pages

#### Cache TTL Strategy

Different types of queries have different TTL values:

- Single items: Default 15 minutes (`REDIS_CACHE_TTL`)
- Paginated results: Default 5 minutes (`REDIS_PAGINATED_TTL`)

#### Cache Invalidation

Automatic cache invalidation when data changes:

- Individual items invalidated on update/delete
- Collection cache invalidated when items change

### Caching Best Practices

For optimal performance:

1. **Configure TTL Values**:
   - Long TTL for static data
   - Short TTL for frequently changing data
   - No caching for real-time critical data

2. **Monitor Cache Efficiency**:
   - Check hit/miss ratios in responses
   - Use debug logging for cache behavior
   - Clear cache with `make flush-redis` during testing

3. **Redis Security**:
   - Use password authentication
   - Consider network security measures
   - Rotate credentials periodically

</details>

<details>
<summary>📝 Logging System</summary>

## Logging System

The API implements a comprehensive logging system with the following features:

- **Multiple rotation types**: Daily rotation (default) and size-based rotation
- **Multiple log outputs**: Console and/or file with thread-safe operations
- **Flexible log formats**: JSON (structured) and console (human-readable)
- **Configurable log levels**: Debug, info, warn, error with hierarchical filtering
- **Retention policies**: Automatic cleanup of old log files
- **Environment-aware defaults**: Different configurations for development and production

### Logging Configuration

Configure logging in your `.env` file:

```bash
# Basic logging configuration
LOG_LEVEL=info                  # Options: debug, info, warn, error
LOG_FORMAT=json                 # Options: json, console
LOG_OUTPUT_PATH=stdout          # Options: stdout, stderr

# Log rotation configuration
LOG_ROTATION_TYPE=daily         # Options: daily, size (default: daily)

# File logging settings
LOG_FILE_PATH=./logs/app.log    # Path to log file (empty = disable file logging)
LOG_FILE_MAX_SIZE=100           # Maximum size in megabytes before rotation (size-based only)
LOG_FILE_MAX_BACKUPS=3          # Maximum number of old log files to retain
LOG_FILE_MAX_AGE=28             # Maximum number of days to retain old log files
LOG_FILE_COMPRESS=true          # Whether to compress rotated log files
```

#### Understanding LOG_LEVEL

The `LOG_LEVEL` setting controls which messages are displayed in your logs, using a hierarchical approach:

| LOG_LEVEL | Debug messages | Info messages | Warning messages | Error messages |
| --------- | -------------- | ------------- | ---------------- | -------------- |
| `debug`   | ✅ Shown        | ✅ Shown       | ✅ Shown          | ✅ Shown        |
| `info`    | ❌ Hidden       | ✅ Shown       | ✅ Shown          | ✅ Shown        |
| `warn`    | ❌ Hidden       | ❌ Hidden      | ✅ Shown          | ✅ Shown        |
| `error`   | ❌ Hidden       | ❌ Hidden      | ❌ Hidden         | ✅ Shown        |

**Guidelines for choosing a level:**

- **Development environments**: Use `debug` to see all logs including detailed debugging information
- **Testing environments**: Use `info` to see normal operational logs plus warnings and errors
- **Production environments**: Use `warn` or `error` to reduce log volume and focus on important issues

**Examples:**

```
# Show all possible logs (development)
LOG_LEVEL=debug

# Show operational logs, warnings and errors (testing)
LOG_LEVEL=info

# Show only warnings and errors (production)
LOG_LEVEL=warn

# Show only errors (production with minimal logging)
LOG_LEVEL=error
```

Each log entry includes a level indicator (L) in its output, like:
```json
{"L":"INFO","T":"2025-05-09T22:15:06.558+0700","C":"bootstrap/server.go:173","M":"Swagger UI available at","url":"http://localhost:8080/swagger/"}
```

### Log Output Options

You can configure the logger to output to:

1. **Standard Output/Error Only**: Set `LOG_OUTPUT_PATH=stdout` and leave `LOG_FILE_PATH` empty
2. **File Only**: Set `LOG_FILE_PATH` to a valid path and `LOG_OUTPUT_PATH` to anything except stdout/stderr
3. **Both Console and File**: Set both `LOG_FILE_PATH` to a valid path and `LOG_OUTPUT_PATH=stdout` or `LOG_OUTPUT_PATH=stderr`

### Log Rotation Types

The logging system supports two rotation strategies:

#### Daily Rotation (Default)

Creates a new log file each day with date-based naming:

```
logs/
├── app-2024-01-15.log    # Today's logs
├── app-2024-01-14.log    # Yesterday's logs
├── app-2024-01-13.log    # Previous day's logs
└── app-2024-01-12.log    # Older daily logs
```

**Benefits:**
- Easy to find logs by specific date
- Clear temporal organization for debugging
- Automatic midnight rotation

**Configuration:**
```bash
LOG_ROTATION_TYPE=daily
```

#### Size-based Rotation

Rotates when log files reach the specified size limit:

```
logs/
├── app.log                           # Current active log file
├── app.log.2024-01-15T10-30-00.000   # Rotated backup 1
├── app.log.2024-01-14T15-45-20.000   # Rotated backup 2
└── app.log.2024-01-13T09-12-45.000   # Rotated backup 3
```

**Benefits:**
- Predictable file sizes for storage planning
- Works well with existing log management tools
- Suitable for high-volume logging scenarios

**Configuration:**
```bash
LOG_ROTATION_TYPE=size
LOG_FILE_MAX_SIZE=100  # Rotate when file reaches 100MB
```

#### Retention Settings

Both rotation types support the same retention policies:

- `LOG_FILE_MAX_BACKUPS`: Maximum number of old log files to retain
- `LOG_FILE_MAX_AGE`: Maximum number of days to retain old log files
- `LOG_FILE_COMPRESS`: Whether to compress rotated log files

### Testing Log Rotation

The project includes a comprehensive script to test both rotation types:

#### Testing Size-based Rotation (Default)
```bash
# Test size-based rotation (quick test, triggers rotation)
make test-log-rotation
# or
make test-log-rotation type=size
# or directly
./scripts/test-log-rotation.sh size
```$$

#### Testing Daily Rotation
```bash
# Test daily rotation setup (verifies daily file creation)
make test-log-rotation type=daily
# or directly
./scripts/test-log-rotation.sh daily
```

#### Test Script Features

**Size-based rotation test:**
1. Creates temporary environment with 1MB file size limit
2. Runs API with size-based rotation enabled
3. Generates log entries to trigger rotation
4. Verifies backup file creation
5. Shows resulting rotated files

**Daily rotation test:**
1. Creates temporary environment with daily rotation
2. Runs API and verifies daily log file creation (`app-YYYY-MM-DD.log`)
3. Generates log entries to populate daily log
4. Shows current daily log file size
5. Provides guidance for testing actual midnight rotation

**Both tests:**
- Backup and restore original `.env` file
- Clean process termination
- Port availability checking
- Comprehensive error handling

#### Usage Examples
```bash
# Quick size-based rotation test
make test-log-rotation

# Test daily rotation setup
make test-log-rotation type=daily

# Get help for test options
./scripts/test-log-rotation.sh --help
```

### Logging Best Practices

For optimal logging:

1. **Choose Appropriate Log Levels**:
   - **Development**: Use `debug` to see all logs including detailed debugging information
   - **Testing**: Use `info` to see normal operational logs plus warnings and errors
   - **Production**: Use `warn` or `error` to reduce log volume and focus on important issues

2. **Select the Right Rotation Type**:
   - **Daily Rotation**: Best for date-specific debugging and temporal organization
   - **Size-based Rotation**: Best for predictable storage usage and high-volume scenarios
   - **Consider your debugging patterns**: Choose daily if you often need to find logs by date

3. **Configure Retention Settings**:
   - Set `LOG_FILE_MAX_AGE` based on compliance and debugging requirements
   - Adjust `LOG_FILE_MAX_BACKUPS` based on available disk space
   - Enable compression (`LOG_FILE_COMPRESS=true`) for long-term storage

4. **Environment-specific Settings**:
   - **Development**: Use console output with daily rotation for easy debugging
   - **Production**: Use file-only output with appropriate log levels and retention
   - **Consider log aggregation systems** for production environments

5. **Performance Considerations**:
   - Use higher log levels in production to reduce I/O
   - Enable compression for rotated files to save disk space
   - Ensure log directories have appropriate permissions and sufficient space

### Migration from Size-based to Daily Rotation

If you're upgrading from a previous version that used only size-based rotation:

**To keep the existing behavior:**
```bash
# Explicitly set size-based rotation
LOG_ROTATION_TYPE=size
```

**To migrate to daily rotation (recommended):**
```bash
# Set daily rotation (or simply remove the variable, as daily is the default)
LOG_ROTATION_TYPE=daily
```

**What changes:**
- **File naming**: From `app.log` → `app-2024-01-15.log`
- **Rotation trigger**: From file size → date change at midnight
- **Organization**: Logs grouped by date instead of size

**Backward compatibility:** All existing configuration variables are still supported.

</details>

<details>
<summary>💾 Database Operations</summary>

## Database Operations

### Migrations

Manage database schema changes:

```bash
# Run all pending migrations
make migrate

# Create a new migration
make migrate-create name=add_new_field

# Create a migration from a model
make migrate-from-model model=animal

# Create migrations from all available models (skip existing tables)
make migrate-all-models
# or use alias
make mam

# Roll back the last migration
make migrate-down

# Check migration status
make migrate-status

# List available models
make migrate-list-models
```

Each migration consists of:
- `[timestamp]_[name].up.sql`: SQL to apply the migration
- `[timestamp]_[name].down.sql`: SQL to roll back the migration

**Migrate All Models Feature:**
The `migrate-all-models` command automatically:
- 🔍 Discovers all available models in the registry
- 🔄 Creates migrations for each model
- ⏭️ Skips models whose tables already exist
- 📊 Provides a summary of operations (created/skipped)
- 🎨 Shows colorful progress with `[1/5]`, `[2/5]` format

### Seeding

Populate the database with test data:

```bash
# Run all seeders
make seed

# Run specific seeder
make seed-animal

# Run with custom count
make seed-count count=500
```

### Seeder Management

The API provides automatic seeder registration to eliminate manual registry updates:

```bash
# Automatically register all seeders (recommended)
make sync-seeder
# or use alias
make ss

# Only add new seeders
make update-seeder-registry
# or use alias
make usr

# Only remove deleted seeders
make clean-seeder-registry
# or use alias
make csr
```

**Automatic Seeder Registration Features:**
- 🔍 **Auto-Discovery**: Scans `pkg/seeder/` for `*_seeder.go` files
- 🔄 **Interface Validation**: Ensures seeders implement the `Seeder` interface
- ✅ **Smart Registration**: Automatically updates `registerSeeders` function
- 🧹 **Cleanup**: Removes references to deleted seeders
- 📊 **Detailed Reporting**: Shows what was added/removed with counts

**Seeder Requirements:**
For automatic detection, seeders must:
1. Be in `pkg/seeder/` directory with `*_seeder.go` naming
2. Have struct name ending with "Seeder" (e.g., `UserSeeder`)
3. Implement `Seed(ctx context.Context) error` and `GetName() string` methods
4. Have constructor function `New{Name}Seeder`

**Example Seeder Structure:**
```go
// pkg/seeder/product_seeder.go
type ProductSeeder struct {
    db     database.Database
    logger *zap.Logger
    count  int
}

func NewProductSeeder(db database.Database, logger *zap.Logger, count int) *ProductSeeder {
    return &ProductSeeder{db: db, logger: logger, count: count}
}

func (s *ProductSeeder) GetName() string { return "product" }
func (s *ProductSeeder) Seed(ctx context.Context) error { /* implementation */ }
```

Once created, run `make sync-seeder` to automatically register it!

### Model Management

The API maintains a registry of models:

```bash
# Add new models
make update-model-map
# or
make um

# Remove deleted models
make clean-model-map
# or
make cm

# Both add and remove models
make sync-model-map
# or
make sm
```

</details>

<details>
<summary>🛠️ Development</summary>

## Development

### Swagger Documentation

Generate and view API documentation:

```bash
# Generate Swagger docs
make swagger

# Run Swagger UI server
make swagger-ui
```

Access the Swagger UI at http://localhost:8090/swagger/

### Development Workflow

The project follows a streamlined development workflow:

1. **Initial Setup**: Clone repo, configure environment
2. **Database Setup**: Run migrations, seed test data
3. **Development Cycle**: Code, test, document
4. **Deployment**: Build and deploy via Docker or Kubernetes

### Using as a Template

This project can be used as a template for new Go APIs:

```bash
# Clone the repository
git clone https://github.com/linkeunid/go-api.git your-project-name
cd your-project-name

# Basic setup - rename module
make setup module=github.com/yourusername/your-project
# or using alias
make setup-s module=github.com/yourusername/your-project

# Setup with Git remote
make setup-git module=github.com/yourusername/your-project \
  remote=git@github.com:yourusername/your-project.git
# or using alias
make setup-g module=github.com/yourusername/your-project \
  remote=git@github.com:yourusername/your-project.git

# Full setup - new Git repo and remote
make setup-full module=github.com/yourusername/your-project \
  remote=git@github.com:yourusername/your-project.git
# or using alias
make setup-f module=github.com/yourusername/your-project \
  remote=git@github.com:yourusername/your-project.git

# Update dependencies
go mod tidy
```

#### What the Setup Tool Updates

The setup-project tool automatically updates the following components to match your new project name:

**📝 Go Module & Imports:**
- Updates `go.mod` module name
- Updates all import paths in Go files

**🐳 Docker Configuration:**
- **Service Names**: `api` → `your-project`, `mysql` → `your-project-mysql`, `redis` → `your-project-redis`
- **Container Names**: `go-api` → `your-project`, `go-mysql` → `your-project-mysql`, `linkeun-redis` → `your-project-redis`
- **Network**: `linkeun-network` → `your-project-network`
- **Volumes**: `mysql_data` → `your-project_mysql_data`, `redis_data` → `your-project_redis_data`
- **Service References**: Updates `depends_on` and environment variable references

**🔧 Git Repository:**
- Optionally resets Git history and creates a new repository
- Sets up new Git remote origin

#### Example Transformation

When you run:
```bash
make setup-f module=github.com/mycompany/awesome-api \
  remote=git@github.com:mycompany/awesome-api.git
```

**Before:**
```yaml
services:
  api:
    container_name: go-api
    depends_on:
      mysql:
        condition: service_healthy
  mysql:
    container_name: go-mysql
  redis:
    container_name: linkeun-redis
networks:
  linkeun-network:
volumes:
  mysql_data:
  redis_data:
```

**After:**
```yaml
services:
  awesome-api:
    container_name: awesome-api
    depends_on:
      awesome-api-mysql:
        condition: service_healthy
  awesome-api-mysql:
    container_name: awesome-api-mysql
  awesome-api-redis:
    container_name: awesome-api-redis
networks:
  awesome-api-network:
volumes:
  awesome-api_mysql_data:
  awesome-api_redis_data:
```

**⚠️ Important Notes:**
- This operation cannot be undone - make sure to backup your project first
- After running the setup, you may need to update your `.env` file if you have custom service configurations
- The tool will show you a detailed preview of all changes before proceeding

</details>

<details>
<summary>🚀 Deployment</summary>

## Deployment

### Docker

```bash
# Build the Docker image
docker build -t go-api:latest .

# Run the container
docker run -p 8080:8080 go-api:latest
```

### Kubernetes

This project includes a complete Kubernetes deployment setup for deploying to Minikube or a production Kubernetes cluster.

#### Minikube Deployment

##### Prerequisites

- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Docker](https://docs.docker.com/get-docker/)

##### Setup Minikube

1. Start Minikube:

```bash
minikube start
```

2. Enable the Ingress addon:

```bash
minikube addons enable ingress
```

##### Automated Deployment

Use the provided deployment script to deploy the application to Minikube:

```bash
# Deploy the application
./k8s/deploy-minikube.sh

# Clean up resources
./k8s/cleanup-minikube.sh
```

##### Manual Deployment

1. Build and load the Docker image:

```bash
# Build the image
docker build -t go-api:latest .

# Load the image into Minikube
minikube image load go-api:latest
```

2. Apply Kubernetes manifests:

```bash
# Apply all manifests at once
kubectl apply -k k8s/

# Or apply each manifest separately
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/mysql.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/deployment.yaml
```

##### Accessing the Application

1. Get the Minikube IP:

```bash
minikube ip
```

2. Add an entry to your hosts file:

```
# /etc/hosts (Linux/Mac) or C:\Windows\System32\drivers\etc\hosts (Windows)
<minikube-ip> go-api.local
```

3. Access the application:

```
http://go-api.local
```

##### Kubernetes Resources

The Kubernetes deployment includes:

- **API Deployment** - The main Go API application
- **MySQL Database** - Persistent MySQL instance
- **Redis** - For caching and session management
- **ConfigMap** - For application configuration
- **Secrets** - For sensitive data
- **Services** - For internal and external networking
- **Ingress** - For external access

##### Monitoring and Debugging

```bash
# Check pod status
kubectl get pods

# View pod logs
kubectl logs <pod-name>

# Execute commands in a pod
kubectl exec -it <pod-name> -- /bin/sh

# Port forwarding for local debugging
kubectl port-forward service/go-api 8080:80
```

#### Production Deployment Considerations

For production deployments, consider:

1. Setting up production secrets management (e.g., Vault or Kubernetes Secrets Store CSI Driver)
2. Configuring proper resource limits and requests
3. Implementing horizontal pod autoscaling
4. Setting up proper database backups and replication
5. Configuring a proper Ingress controller with TLS
6. Implementing proper monitoring and alerting

</details>

<details>
<summary>🧪 API Testing</summary>

## API Testing

The project includes comprehensive unit tests for the REST API components, focusing on:

- Controller Layer: Tests for all API endpoints and edge cases
- Service Layer: Tests for business logic handling
- Repository Mocking: Simulated data store interactions 

### Running the Tests

The project includes several testing commands:

```bash
# Run all tests
make test

# Run all tests with coverage
make test-coverage

# Run only API-related tests
make test-api

# Run only service-related tests  
make test-service

# Run only repository-related tests
make test-repository

# Run tests with race detection
make test-race
```

### Test Coverage

The test suite achieves high coverage percentages for critical components:
- Controller Layer: 92.6% coverage
- Service Layer: 100% coverage

Coverage reports are generated in HTML format and can be found in the `./coverage` directory after running the tests.

### Testing Strategy

The tests use the following strategies:
- **Mocking**: Service dependencies are mocked to isolate units
- **Table-Driven Tests**: Tests cover multiple test cases efficiently
- **Edge Cases**: Tests handle error conditions and invalid inputs
- **Assertions**: Tests verify correct behavior and outputs

### Extending Tests

When adding new API endpoints, follow the same pattern to extend the tests:
1. Create controller tests for HTTP request/response handling
2. Create service tests for business logic
3. Create repository tests or mocks for data access
4. Run tests to ensure coverage is maintained

</details>

## License

This project is licensed under the GNU General Public License v2.0. See the LICENSE file for details.
