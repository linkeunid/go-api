# Linkeun Go API

A comprehensive Go API project with RESTful endpoints, JWT authentication, caching, and database integration.

## Table of Contents

- [Linkeun Go API](#linkeun-go-api)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Environment Setup](#environment-setup)
    - [Running the API](#running-the-api)
      - [With Docker](#with-docker)
      - [Locally](#locally)
    - [API Endpoints](#api-endpoints)
      - [Animals Resource](#animals-resource)
      - [Query Parameters](#query-parameters)
  - [Project Structure](#project-structure)
  - [Authentication](#authentication)
    - [JWT Overview](#jwt-overview)
    - [Token Generation](#token-generation)
    - [Claims Structure](#claims-structure)
    - [Configuration Options](#configuration-options)
      - [Understanding JWT\_ALLOWED\_ISSUERS](#understanding-jwt_allowed_issuers)
    - [Implementation Details](#implementation-details)
      - [Using Authentication](#using-authentication)
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
  - [Database Operations](#database-operations)
    - [Migrations](#migrations)
    - [Seeding](#seeding)
    - [Model Management](#model-management)
  - [Development](#development)
    - [Swagger Documentation](#swagger-documentation)
    - [Development Workflow](#development-workflow)
    - [Using as a Template](#using-as-a-template)
  - [Deployment](#deployment)
    - [Docker](#docker)
    - [Kubernetes](#kubernetes)
  - [License](#license)

## Features

- RESTful API endpoints with CRUD operations
- JWT authentication with role-based access control
- MySQL database integration with GORM ORM
- Redis-based caching system for performance optimization
- Pagination, sorting, and filtering support
- Docker and Kubernetes deployment configurations
- API documentation with Swagger
- Database migrations and seeding
- Environment-specific configurations
- Comprehensive error handling
- Logging with configurable levels and outputs

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

## Project Structure

```
.
├── cmd/                      # Command-line applications
│   ├── api/                  # Main API application
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
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:4445/api/v1/protected-endpoint
```

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

# Roll back the last migration
make migrate-down

# Check migration status
make migrate-status
```

Each migration consists of:
- `[timestamp]_[name].up.sql`: SQL to apply the migration
- `[timestamp]_[name].down.sql`: SQL to roll back the migration

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

# Full setup - new Git repo and remote
make setup-full module=github.com/yourusername/your-project \
  remote=git@github.com:yourusername/your-project.git

# Update dependencies
go mod tidy
```

## Deployment

### Docker

Run with Docker Compose:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

### Kubernetes

Deploy to a Kubernetes cluster:

```bash
# Apply configuration
kubectl apply -f k8s/deployment.yaml
```

## License

This project is licensed under the GNU General Public License v2.0. See the LICENSE file for details. 