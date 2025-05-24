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
    - [Log Rotation](#log-rotation)
    - [Testing Log Rotation](#testing-log-rotation)
    - [Logging Best Practices](#logging-best-practices)
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

# Logging configuration
LOG_LEVEL=info                  # Options: debug, info, warn, error
LOG_FORMAT=json                 # Options: json, console
LOG_OUTPUT_PATH=stdout          # Options: stdout, stderr
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

| Alias | Full Command       | Description                   |
| ----- | ------------------ | ----------------------------- |
| `sd`  | `seed`             | Run all database seeders      |
| `tr`  | `truncate`         | Truncate specific table       |
| `um`  | `update-model-map` | Update model map for database |

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
make gt     # Standard token
make gta    # Admin token
make gtu id=123  # Token for specific user ID
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
        E -->|make migrate| F[Seed Database]
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
        Q2[Database Operations] -->|make um, make sm| J
        Q3[Monitor Containers] -->|make fps, make docker-logs / make dlogs| J
        Q4[Flush Cache] -->|make flush-redis / make fr| J
    end
```

<details>
<summary>🏗️ Project Structure</summary>

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

The API implements a configurable logging system with the following features:

- Multiple log outputs (console and/or file)
- Log rotation for file-based logs
- Configurable log levels and formats
- Retention policies for old log files

### Logging Configuration

Configure logging in your `.env` file:

```
# Basic logging configuration
LOG_LEVEL=info                  # Options: debug, info, warn, error
LOG_FORMAT=json                 # Options: json, console
LOG_OUTPUT_PATH=stdout          # Options: stdout, stderr

# File logging with rotation
LOG_FILE_PATH=./logs/app.log    # Path to log file (empty = disable file logging)
LOG_FILE_MAX_SIZE=100           # Maximum size in megabytes before rotation
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

### Log Rotation

When file logging is enabled, logs will automatically rotate when:

- The file reaches the specified maximum size (`LOG_FILE_MAX_SIZE` in megabytes)
- The maximum number of backup files is limited by `LOG_FILE_MAX_BACKUPS`
- Old log files older than `LOG_FILE_MAX_AGE` days are automatically deleted

### Testing Log Rotation

The project includes a script to test log rotation functionality:

```bash
# Run the log rotation test
make test-log-rotation
# or
./scripts/test-log-rotation.sh
```

This script:
1. Creates a temporary environment with small log file size limits
2. Runs the API with file logging enabled
3. Generates enough log entries to trigger rotation
4. Shows the resulting log files
5. Restores the original environment configuration

### Logging Best Practices

For optimal logging:

1. **Choose Appropriate Log Levels**:
   - Use `debug` for detailed troubleshooting
   - Use `info` for general operational information
   - Use `warn` for concerning but non-critical issues
   - Use `error` for serious problems

2. **Configure Rotation Settings**:
   - Set reasonable file size limits based on disk space
   - Adjust retention period based on compliance requirements
   - Enable compression for long-term storage

3. **Production Settings**:
   - In production, consider using a log aggregation system
   - Set higher log levels (`warn` or `error`) to reduce disk I/O
   - Ensure log directories have appropriate permissions

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
