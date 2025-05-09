.PHONY: swagger swagger-ui dev test lint fmt help

# Default target - show help
help:
	@echo "✨ Linkeun Go API - Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run development server with hot reload"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make swagger-ui     - Run Swagger UI server"
	@echo "  make init           - Initialize the project"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"
	@echo "  make mocks          - Generate mocks for testing"
	@echo ""
	@echo "Database migrations:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-create name=NAME - Create a new migration"
	@echo "  make migrate-from-model model=NAME - Create a migration from a model"
	@echo "  make migrate-list-models - List available models for migrations"
	@echo "  make migrate-down   - Roll back the last migration"
	@echo "  make migrate-status - Show current migration status"
	@echo "  make migrate-reset  - Reset all migrations"
	@echo ""
	@echo "  make env-info       - Show environment variables used by the application"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-db      - Start only database containers (MySQL and Redis)"
	@echo "  make docker-up      - Start all containers"
	@echo "  make docker-down    - Stop all containers"
	@echo "  make docker-rebuild - Rebuild and restart all containers"
	@echo "  make docker-logs    - View logs from all containers"
	@echo "  make docker-ps      - List running containers"
	@echo "  make fancy-ps       - Show fancy container status with colors and details 🌈"
	@echo "  make docker-clean   - Remove all containers, volumes, and images"
	@echo ""
	@echo "Aliases:"
	@echo "  make s              - Generate Swagger documentation"
	@echo "  make su             - Run Swagger UI server"
	@echo "  make r              - Run the application"
	@echo "  make d              - Run development server with hot reload"
	@echo "  make t              - Run tests"
	@echo "  make l              - Lint code"
	@echo "  make ddb            - Start database containers"
	@echo "  make dup            - Start all containers"
	@echo "  make ddown          - Stop all containers"
	@echo "  make fps            - Show fancy container status"
	@echo "  make ei             - Show environment info"

# Helper function to get env variable with default value
# Usage: $(call get_env,VARIABLE_NAME,DEFAULT_VALUE)
define get_env
$(shell if [ -f .env ]; then grep -E "^$(1)=" .env | cut -d= -f2 || echo "$(2)"; else echo "$(2)"; fi)
endef

# Build the application
build:
	@echo "🔨 Building application..."
	@go build -o bin/api ./cmd/api
	@echo "✅ Build complete: ./bin/api"

# Run the application
run:
	@echo "🚀 Starting application..."
	@go run ./cmd/api

# Development mode with hot reload (alias)
dev:
	@echo "🔄 Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
	else \
		echo "⚠️ Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air -c .air.toml; \
	fi

# Generate Swagger documentation
swagger:
	@echo "📝 Generating Swagger documentation..."
	@./scripts/swagger.sh
	@echo "✅ Swagger documentation generated"
	@echo "ℹ️ Note: You may see warnings about 'no Go files in root directory' - this is normal and can be ignored"

# Run Swagger UI
swagger-ui: swagger
	@echo "🌐 Starting Swagger UI server..."
	@./scripts/swagger-ui.sh

# Install Swagger tools
swagger-tools:
	@echo "⚙️ Installing Swagger tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Swagger tools installed"

# Initialize the project (download dependencies, generate swagger, etc.)
init: swagger-tools swagger
	@echo "🔧 Initializing project..."
	@echo "✅ Project initialized successfully"

# Clean the project
clean:
	@echo "🧹 Cleaning project..."
	@rm -rf bin/
	@rm -rf internal/docs/swaggerdocs/
	@echo "✅ Project cleaned"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated at coverage.html"

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Lint code
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️ golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
		golangci-lint run ./...; \
	fi

# Generate mocks for testing
mocks:
	@echo "🧩 Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		go generate ./...; \
	else \
		echo "⚠️ mockgen not found. Installing..."; \
		go install github.com/golang/mock/mockgen@latest; \
		go generate ./...; \
	fi
	@echo "✅ Mocks generated"

# Run database migrations
migrate:
	@echo "🗃️ Running database migrations..."
	@go run ./cmd/migrate -up
	@echo "✅ Migrations completed"

# Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Migration name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "🗃️ Creating new migration: $(name)..."
	@go run ./cmd/migrate -create $(name)
	@echo "✅ Migration files created"

# Create a migration from a model
migrate-from-model:
	@if [ -z "$(model)" ]; then \
		echo "❌ Model name is required. Usage: make migrate-from-model model=animal"; \
		echo "Available models:"; \
		go run ./cmd/migrate -list-models | grep -v "Available"; \
		exit 1; \
	fi
	@echo "🗃️ Creating migration from model: $(model)..."
	@go run ./cmd/migrate -create -from-model $(model)
	@echo "✅ Model migration files created"

# List available models for migration
migrate-list-models:
	@echo "🗃️ Available models for migrations:"
	@go run ./cmd/migrate -list-models

# Roll back the last migration
migrate-down:
	@echo "🗃️ Rolling back the last migration..."
	@go run ./cmd/migrate -down
	@echo "✅ Migration rollback completed"

# Check migration status
migrate-status:
	@echo "🗃️ Checking migration status..."
	@go run ./cmd/migrate -version
	
# Reset all migrations
migrate-reset:
	@echo "⚠️ This will reset all migrations! Are you sure? (y/n)"
	@read -r answer; \
	if [ "$$answer" = "y" ]; then \
		echo "🗃️ Resetting all migrations..."; \
		go run ./cmd/migrate -force 0; \
		echo "✅ All migrations have been reset"; \
	else \
		echo "❌ Operation cancelled"; \
	fi

# Aliases for common commands
s: swagger
su: swagger-ui
r: run
d: dev
t: test
l: lint

# Docker commands
docker-db: ## Start only database containers (MySQL and Redis)
	@echo "🐳 Starting database containers..."
	@docker compose -f docker-compose.yml up -d mysql redis
	@echo "✅ Database containers started"
	@echo "   MySQL: $(call get_env,DB_HOST,localhost):$(call get_env,DB_PORT,3306)"
	@echo "   Redis: $(call get_env,REDIS_HOST,localhost):$(call get_env,REDIS_PORT,6379)"

docker-up: ## Start all containers
	@echo "🐳 Starting all containers..."
	@docker compose -f docker-compose.yml up -d
	@echo "✅ All containers started"
	@echo "   API: http://$(call get_env,API_HOST,localhost):$(call get_env,API_PORT,8080)"
	@echo "   MySQL: $(call get_env,DB_HOST,localhost):$(call get_env,DB_PORT,3306)"
	@echo "   Redis: $(call get_env,REDIS_HOST,localhost):$(call get_env,REDIS_PORT,6379)"

docker-down: ## Stop all containers
	@echo "🐳 Stopping all containers..."
	@docker compose -f docker-compose.yml down
	@echo "✅ All containers stopped"

docker-rebuild: ## Rebuild and restart all containers
	@echo "🐳 Rebuilding all containers..."
	@docker compose -f docker-compose.yml down
	@docker compose -f docker-compose.yml build
	@docker compose -f docker-compose.yml up -d
	@echo "✅ All containers rebuilt and started"
	@echo "   API: http://$(call get_env,API_HOST,localhost):$(call get_env,API_PORT,8080)"
	@echo "   MySQL: $(call get_env,DB_HOST,localhost):$(call get_env,DB_PORT,3306)"
	@echo "   Redis: $(call get_env,REDIS_HOST,localhost):$(call get_env,REDIS_PORT,6379)"

docker-logs: ## View logs from all containers
	@echo "📋 Showing container logs (press Ctrl+C to exit)..."
	@docker compose -f docker-compose.yml logs -f

docker-ps: ## List running containers
	@echo "📋 Running containers:"
	@docker compose -f docker-compose.yml ps

fancy-ps: ## Show fancy container status with colors and details
	@echo ""
	@echo "✨ 🐳 \033[1;35mFancy Container Status\033[0m 🐳 ✨"
	@echo ""
	@echo "\033[1;36m┌───────────────────────────────────────────────────┐\033[0m"
	@echo "\033[1;36m│ 🔍 CONTAINER STATUS                               │\033[0m"
	@echo "\033[1;36m└───────────────────────────────────────────────────┘\033[0m"
	@docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Image}}' 2>/dev/null || echo "   No running containers found"
	@echo ""
	@echo "\033[1;36m┌───────────────────────────────────────────────────┐\033[0m"
	@echo "\033[1;36m│ 📊 RESOURCE USAGE                                 │\033[0m"
	@echo "\033[1;36m└───────────────────────────────────────────────────┘\033[0m"
	@docker stats --no-stream --format 'table {{.Name}}\tCPU: {{.CPUPerc}}\tMEM: {{.MemPerc}} ({{.MemUsage}})' 2>/dev/null || echo "   No stats available"
	@echo ""
	@echo "\033[1;36m┌───────────────────────────────────────────────────┐\033[0m"
	@echo "\033[1;36m│ 🔌 NETWORK INFO                                   │\033[0m"
	@echo "\033[1;36m└───────────────────────────────────────────────────┘\033[0m"
	@docker network ls --format 'table {{.Name}}\t{{.Driver}}\t{{.Scope}}' 2>/dev/null || echo "   No networks found"
	@echo ""
	@echo "\033[1;36m┌───────────────────────────────────────────────────┐\033[0m"
	@echo "\033[1;36m│ 💾 PROJECT VOLUMES                                │\033[0m"
	@echo "\033[1;36m└───────────────────────────────────────────────────┘\033[0m"
	@docker volume ls --filter name=go-api --format 'table {{.Name}}\t{{.Driver}}' 2>/dev/null || echo "   No volumes found"
	@echo ""
	@echo "\033[1;36m┌───────────────────────────────────────────────────┐\033[0m"
	@echo "\033[1;36m│ 🚀 HELPFUL COMMANDS                               │\033[0m"
	@echo "\033[1;36m└───────────────────────────────────────────────────┘\033[0m"
	@echo "   \033[1;33mmake docker-logs\033[0m     → View container logs"
	@echo "   \033[1;33mmake docker-rebuild\033[0m  → Rebuild and restart containers"
	@echo "   \033[1;33mmake docker-clean\033[0m    → Clean up Docker resources"
	@echo ""
	@echo "   \033[1;32mAdd to your shell config:\033[0m"
	@echo ""
	@echo "   \033[1;36m# For bash (.bashrc):\033[0m"
	@echo "   alias fps='make -C $(PWD) fps'"
	@echo ""
	@echo "   \033[1;36m# For zsh (.zshrc):\033[0m"
	@echo "   alias fps='make -C $(PWD) fps'"
	@echo ""

docker-clean: ## Remove all containers, volumes, and images
	@echo "🧹 Cleaning up Docker resources..."
	@docker compose -f docker-compose.yml down -v
	@docker system prune -af --volumes
	@echo "✅ Docker cleanup complete"

# Docker command aliases
ddb: docker-db
dup: docker-up
ddown: docker-down
dps: docker-ps
dlogs: docker-logs
fps: fancy-ps

# Show environment variables
env-info: ## Show environment variables used by the application
	@echo "🔍 Environment variables (from .env file if present):"
	@echo "   API_HOST: $(call get_env,API_HOST,localhost)"
	@echo "   API_PORT: $(call get_env,API_PORT,8080)"
	@echo "   DB_HOST: $(call get_env,DB_HOST,localhost)"
	@echo "   DB_PORT: $(call get_env,DB_PORT,3306)"
	@echo "   DB_USER: $(call get_env,DB_USER,root)"
	@echo "   DB_NAME: $(call get_env,DB_NAME,linkeun_go_api)"
	@echo "   REDIS_HOST: $(call get_env,REDIS_HOST,localhost)"
	@echo "   REDIS_PORT: $(call get_env,REDIS_PORT,6379)"
	@echo "   REDIS_ENABLED: $(call get_env,REDIS_ENABLED,false)"
	@echo "📝 Note: Values shown are actual values from .env or defaults if not defined"

# Other aliases
ei: env-info 