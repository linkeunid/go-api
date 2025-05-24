.PHONY: swagger swagger-ui dev test lint fmt help

# Helper function to print colorful messages
# Usage: $(call print_colorful, emoji, text, color, is_bold)
# Example: $(call print_colorful, 🚀, Building application, green, true)
define print_colorful
	printf "\033[$(if $(4),1;)$(if $(3),$(3),32)m$(1) $(2)\033[0m\n"
endef

# Helper function for confirmation prompts
# Usage: $(call ask_confirmation, warning_message, action_message, emoji)
# Example: $(call ask_confirmation, This will remove all files!, Cleaning all files, 🧹)
# Returns "true" if confirmed, "false" if not
define ask_confirmation
	printf "\033[1;$(YELLOW)m⚠️ $(1)\033[0m\n"; \
	printf "\033[1;$(YELLOW)m⚠️ Are you sure you want to proceed? (y/n): \033[0m"; \
	read -r answer; \
	if [ "$$answer" = "y" ] || [ "$$answer" = "Y" ]; then \
		printf "\033[1;$(YELLOW)m$(if $(3),$(3),🔄) $(2)...\033[0m\n"; \
		true; \
	else \
		printf "\033[1;$(RED)m❌ Operation cancelled\033[0m\n"; \
		false; \
	fi
endef

# Helper function for printing help lines with consistent formatting
# Usage: $(call print_help_line, command, description)
# Example: $(call print_help_line, make dev, 🔄 Start development server with hot reload and Redis cache flush)
define print_help_line
	@printf "  %-30s - %s\n" "$(1)" "$(2)"
endef

# Color codes for echo statements
GREEN := 32
YELLOW := 33
BLUE := 34
CYAN := 36
MAGENTA := 35
RED := 31

# Display fancy header for commands that deserve special attention
define fancy_header
	@printf "\n\033[1;$(1)m════════════════════════════════════════════════════════════════\033[0m\n"
	@printf "\033[1;$(1)m 🌟 $(2)\033[0m\n"
	@printf "\033[1;$(1)m════════════════════════════════════════════════════════════════\033[0m\n\n"
endef

# Default target - show help
help:
	$(call fancy_header,$(MAGENTA),Linkeun Go API - Command Reference)
	@printf "\033[1;36m✨ Build & Run\033[0m\n"
	$(call print_help_line, make dev, 🔄 Start development server with hot reload and Redis cache flush)
	$(call print_help_line, make run, 🚀 Start the API application in production mode)
	$(call print_help_line, make build, 🔨 Compile Go source code into executable binary)
	@printf "\n"
	@printf "\033[1;36m🐳 Docker Commands\033[0m\n"
	$(call print_help_line, make docker-up, 🚀 Start all containers defined in docker-compose.yml)
	$(call print_help_line, make docker-down, 🛑 Stop and remove all running containers)
	$(call print_help_line, make docker-db, 🗃️ Start only database containers (MySQL and Redis))
	$(call print_help_line, make docker-rebuild, 🔄 Rebuild Docker images and restart all containers)
	$(call print_help_line, make fancy-ps, 🌈 Display detailed container status with resource usage and colors)
	$(call print_help_line, make docker-logs, 📋 Stream real-time logs from all running containers)
	$(call print_help_line, make docker-ps, 📊 List currently running Docker containers)
	$(call print_help_line, make docker-clean, 🧹 Remove project containers and volumes with confirmation)
	@printf "\n"
	@printf "\033[1;36m🗃️ Database Migrations\033[0m\n"
	$(call print_help_line, make migrate, 📊 Execute all pending database schema migrations)
	$(call print_help_line, make migrate-status, 📊 Display current migration version and pending migrations)
	$(call print_help_line, make migrate-down, ⏮️ Rollback the most recent migration with confirmation)
	$(call print_help_line, make migrate-create name=NAME, 📝 Generate new empty migration files with timestamp)
	$(call print_help_line, make migrate-from-model model=NAME, 🔄 Auto-generate migration from existing model structure)
	$(call print_help_line, make migrate-all-models, 🚀 Create migrations from all available models (skip existing))
	$(call print_help_line, make migrate-list-models, 📋 Show all models available for migration generation)
	$(call print_help_line, make migrate-reset, 🔄 Rollback all migrations to version 0 with confirmation)
	@printf "\n"
	@printf "\033[1;36m🌱 Database Seeders\033[0m\n"
	$(call print_help_line, make seed, 🌱 Populate database with test data from all available seeders)
	$(call print_help_line, make seed-count, 🔢 Run all seeders with custom record count (e.g., make seed-count count=100))
	$(call print_help_line, make seed-animal, 🐾 Populate database with animal test data only)
	$(call print_help_line, make seed-flower, 🌸 Populate database with flower test data only)
	@printf "\n"
	@printf "\033[1;36m💾 Database Operations\033[0m\n"
	$(call print_help_line, make update-model-map, 🔄 Scan and register new Go models for database operations)
	$(call print_help_line, make sync-model-map, 🔄 Add new models and remove deleted ones from the registry)
	$(call print_help_line, make clean-model-map, 🧹 Remove models from registry that no longer exist in codebase)
	$(call print_help_line, make truncate model=NAME, 🗑️ Empty specific database table after user confirmation)
	$(call print_help_line, make truncate-all, 🧹 Empty all database tables after double confirmation)
	@printf "\n"
	@printf "\033[1;36m🧪 Testing & Quality\033[0m\n"
	$(call print_help_line, make test, 🧪 Execute all unit and integration tests with verbose output)
	$(call print_help_line, make fmt, ✨ Format all Go source code using go fmt)
	$(call print_help_line, make lint, 🔍 Run golangci-lint to check code quality and style)
	$(call print_help_line, make test-coverage, 📊 Run tests and generate detailed coverage report)
	$(call print_help_line, make test-log-rotation, 📋 Test application log file rotation functionality)
	$(call print_help_line, make mocks, 🧩 Generate mock interfaces for unit testing with mockgen)
	@printf "\n"
	@printf "\033[1;36m📚 Documentation\033[0m\n"
	$(call print_help_line, make swagger, 📝 Generate OpenAPI/Swagger documentation from code annotations)
	$(call print_help_line, make swagger-ui, 🌐 Start interactive Swagger UI server on localhost:8081)
	@printf "\n"
	@printf "\033[1;36m🔐 JWT Authentication\033[0m\n"
	$(call print_help_line, make generate-token, 🔑 Generate JWT token with default user settings (dev/test only))
	$(call print_help_line, make generate-token-user id=123, 👤 Generate JWT token for specific user ID (dev/test only))
	$(call print_help_line, make generate-token-admin, 👑 Generate JWT token with administrator privileges (dev/test only))
	$(call print_help_line, make generate-token-force, ⚡ Force JWT token generation bypassing environment restrictions)
	@printf "\n"
	@printf "\033[1;36m🔧 Project Management\033[0m\n"
	$(call print_help_line, make init, 🔧 Initialize project dependencies and generate documentation)
	$(call print_help_line, make env-info, ℹ️ Display all environment variables used by the application)
	$(call print_help_line, make env-info show=all, 🔓 Display environment variables with sensitive values revealed)
	$(call print_help_line, make clean, 🧹 Remove build artifacts and logs with user confirmation)
	$(call print_help_line, make clean-all, 🧹 Remove build artifacts and logs without confirmation (CI/automation))
	$(call print_help_line, make clean-logs, 🧹 Remove application log files only with confirmation)
	$(call print_help_line, make flush-redis, 🧹 Clear Redis cache database using configured credentials)
	@printf "\n"
	@printf "\033[1;36m🔧 Project Template Setup\033[0m\n"
	$(call print_help_line, make setup module=MODULE_NAME, 🛠️ Rename Go module and update all import paths throughout codebase)
	$(call print_help_line, make setup-git remote=GIT_URL, 🔄 Setup project with new module name and configure Git remote repository)
	$(call print_help_line, make setup-full module=MODULE_NAME remote=GIT_URL, 🚀 Complete project setup with module rename and fresh Git repository)
	@printf "\n"
	@printf "\033[1;36m⚡ Command Aliases\033[0m\n"
	$(call print_help_line, make d, ↩️ Alias for 'dev')
	$(call print_help_line, make r, ↩️ Alias for 'run')
	$(call print_help_line, make s, ↩️ Alias for 'swagger')
	$(call print_help_line, make su, ↩️ Alias for 'swagger-ui')
	$(call print_help_line, make t, ↩️ Alias for 'test')
	$(call print_help_line, make l, ↩️ Alias for 'lint')
	$(call print_help_line, make sd, ↩️ Alias for 'seed')
	$(call print_help_line, make tr, ↩️ Alias for 'truncate')
	$(call print_help_line, make tra, ↩️ Alias for 'truncate-all')
	$(call print_help_line, make mam, ↩️ Alias for 'migrate-all-models')
	$(call print_help_line, make um, ↩️ Alias for 'update-model-map')
	$(call print_help_line, make cm, ↩️ Alias for 'clean-model-map')
	$(call print_help_line, make sm, ↩️ Alias for 'sync-model-map')
	$(call print_help_line, make smm, ↩️ Alias for 'sync-model-map')
	$(call print_help_line, make fr, ↩️ Alias for 'flush-redis')
	$(call print_help_line, make gt, ↩️ Alias for 'generate-token')
	$(call print_help_line, make gtu, ↩️ Alias for 'generate-token-user')
	$(call print_help_line, make gta, ↩️ Alias for 'generate-token-admin')
	$(call print_help_line, make gtf, ↩️ Alias for 'generate-token-force')
	$(call print_help_line, make tlr, ↩️ Alias for 'test-log-rotation')
	$(call print_help_line, make cl, ↩️ Alias for 'clean-logs')
	$(call print_help_line, make ca, ↩️ Alias for 'clean-all')
	$(call print_help_line, make setup-s, ↩️ Alias for 'setup')
	$(call print_help_line, make setup-g, ↩️ Alias for 'setup-git')
	$(call print_help_line, make setup-f, ↩️ Alias for 'setup-full')
	$(call print_help_line, make ddb, ↩️ Alias for 'docker-db')
	$(call print_help_line, make dup, ↩️ Alias for 'docker-up')
	$(call print_help_line, make ddown, ↩️ Alias for 'docker-down')
	$(call print_help_line, make dps, ↩️ Alias for 'docker-ps')
	$(call print_help_line, make dlogs, ↩️ Alias for 'docker-logs')
	$(call print_help_line, make fps, ↩️ Alias for 'fancy-ps')
	$(call print_help_line, make ei, ↩️ Alias for 'env-info')

# Helper function to get env variable with default value
# Usage: $(call get_env,VARIABLE_NAME,DEFAULT_VALUE)
define get_env
$(shell if [ -f .env ]; then grep -E "^$(1)=" .env | cut -d= -f2 || echo "$(2)"; else echo "$(2)"; fi)
endef

# Helper function to flush Redis cache
# Usage: $(call flush_redis_cache)
define flush_redis_cache
	@printf "\033[1;$(YELLOW)m🧹 Attempting to flush Redis cache...\033[0m\n"
	@if command -v redis-cli > /dev/null; then \
		REDIS_HOST=$$(grep -E "^REDIS_HOST=" .env 2>/dev/null | cut -d= -f2 || echo "localhost"); \
		REDIS_PORT=$$(grep -E "^REDIS_PORT=" .env 2>/dev/null | cut -d= -f2 || echo "6380"); \
		REDIS_PASSWORD=$$(grep -E "^REDIS_PASSWORD=" .env 2>/dev/null | cut -d= -f2 || echo ""); \
		if [ -z "$$REDIS_HOST" ]; then \
			REDIS_HOST="localhost"; \
			printf "\033[$(YELLOW)m⚠️ REDIS_HOST is empty, using default: localhost\033[0m\n"; \
		fi; \
		if [ -z "$$REDIS_PORT" ]; then \
			REDIS_PORT="6380"; \
			printf "\033[$(YELLOW)m⚠️ REDIS_PORT is empty, using default: 6380\033[0m\n"; \
		fi; \
		printf "\033[$(BLUE)m🔗 Connecting to Redis at $$REDIS_HOST:$$REDIS_PORT\033[0m\n"; \
		if [ -n "$$REDIS_PASSWORD" ]; then \
			printf "\033[$(CYAN)m🔑 Using password from .env file\033[0m\n"; \
			AUTH_RESULT=$$(redis-cli -h $$REDIS_HOST -p $$REDIS_PORT AUTH "$$REDIS_PASSWORD" 2>&1); \
			if echo "$$AUTH_RESULT" | grep -q "OK"; then \
				printf "\033[$(GREEN)m✅ Authentication successful\033[0m\n"; \
				FLUSH_RESULT=$$(redis-cli -h $$REDIS_HOST -p $$REDIS_PORT -a "$$REDIS_PASSWORD" FLUSHDB 2>&1); \
				if echo "$$FLUSH_RESULT" | grep -q "OK"; then \
					printf "\033[1;$(GREEN)m✅ Redis cache flushed successfully\033[0m\n"; \
				else \
					printf "\033[1;$(RED)m❌ Flush failed: $$FLUSH_RESULT\033[0m\n"; \
				fi; \
			else \
				printf "\033[1;$(YELLOW)m⚠️ Authentication failed: $$AUTH_RESULT\033[0m\n"; \
				printf "\033[$(YELLOW)m🔑 Trying with default fallback method...\033[0m\n"; \
				FLUSH_RESULT=$$(redis-cli -h $$REDIS_HOST -p $$REDIS_PORT -a "redis" FLUSHDB 2>&1); \
				if echo "$$FLUSH_RESULT" | grep -q "OK"; then \
					printf "\033[1;$(GREEN)m✅ Redis cache flushed successfully with fallback\033[0m\n"; \
				else \
					printf "\033[1;$(RED)m❌ All authentication methods failed\033[0m\n"; \
					printf "\033[$(RED)m   Please verify Redis is running and credentials are correct\033[0m\n"; \
					printf "\033[$(RED)m   Redis error: $$FLUSH_RESULT\033[0m\n"; \
				fi; \
			fi; \
		else \
			printf "\033[$(YELLOW)m🔑 No password found in .env file, trying without authentication...\033[0m\n"; \
			FLUSH_RESULT=$$(redis-cli -h $$REDIS_HOST -p $$REDIS_PORT FLUSHDB 2>&1); \
			if echo "$$FLUSH_RESULT" | grep -q "OK"; then \
				printf "\033[1;$(GREEN)m✅ Redis cache flushed successfully\033[0m\n"; \
			elif echo "$$FLUSH_RESULT" | grep -q "NOAUTH"; then \
				printf "\033[$(YELLOW)m🔑 Authentication required, trying with default 'redis' password...\033[0m\n"; \
				FALLBACK_RESULT=$$(redis-cli -h $$REDIS_HOST -p $$REDIS_PORT -a "redis" FLUSHDB 2>&1); \
				if echo "$$FALLBACK_RESULT" | grep -q "OK"; then \
					printf "\033[1;$(GREEN)m✅ Redis cache flushed successfully with default password\033[0m\n"; \
				else \
					printf "\033[1;$(RED)m❌ Could not authenticate to Redis\033[0m\n"; \
					printf "\033[$(RED)m   Please add the correct password to your .env file\033[0m\n"; \
					printf "\033[$(RED)m   Redis error: $$FALLBACK_RESULT\033[0m\n"; \
				fi; \
			else \
				printf "\033[1;$(YELLOW)m⚠️ Could not connect to Redis: $$FLUSH_RESULT\033[0m\n"; \
				printf "\033[1;$(YELLOW)m⚠️ Continuing without flushing Redis cache\033[0m\n"; \
			fi; \
		fi; \
	else \
		printf "\033[1;$(YELLOW)m⚠️ redis-cli not found. Skipping cache flush.\033[0m\n"; \
	fi
endef

# Build the application
build:
	@printf "\033[1;$(BLUE)m🔨 Building application...\033[0m\n"
	@go build -o bin/api ./cmd/api
	@printf "\033[$(GREEN)m✅ Build complete: ./bin/api\033[0m\n"

# Run the application
run:
	@printf "\033[1;$(GREEN)m🚀 Starting application...\033[0m\n"
	@go run ./cmd/api

# Development mode with hot reload (alias)
dev:
	@printf "\033[1;$(CYAN)m🔄 Starting development server with hot reload...\033[0m\n"
	$(call flush_redis_cache)
	@if command -v air > /dev/null; then \
		air -c .air.toml; \
	else \
		printf "\033[1;$(YELLOW)m⚠️ Air not found. Installing...\033[0m\n"; \
		go install github.com/cosmtrek/air@latest; \
		air -c .air.toml; \
	fi

# Generate Swagger documentation
swagger:
	@printf "\033[1;$(BLUE)m📝 Generating Swagger documentation...\033[0m\n"
	@./scripts/swagger.sh
	@printf "\033[$(GREEN)m✅ Swagger documentation generated\033[0m\n"
	@printf "\033[$(YELLOW)mℹ️ Note: You may see warnings about 'no Go files in root directory' - this is normal and can be ignored\033[0m\n"

# Run Swagger UI
swagger-ui: swagger
	@printf "\033[1;$(MAGENTA)m🌐 Starting Swagger UI server...\033[0m\n"
	@./scripts/swagger-ui.sh

# Install Swagger tools
swagger-tools:
	@printf "\033[1;$(BLUE)m⚙️ Installing Swagger tools...\033[0m\n"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@printf "\033[$(GREEN)m✅ Swagger tools installed\033[0m\n"

# Initialize the project (download dependencies, generate swagger, etc.)
init: swagger-tools swagger update-model-map
	@printf "\033[1;$(BLUE)m🔧 Initializing project...\033[0m\n"
	@printf "\033[1;$(GREEN)m✅ Project initialized successfully\033[0m\n"

# Clean the project
clean:
	@if $(call ask_confirmation, This will remove all generated files including logs, builds, and docs!, Cleaning project, 🧹); then \
		rm -rf bin/; \
		rm -rf internal/docs/swaggerdocs/; \
		rm -rf coverage/; \
		rm -rf logs/; \
		find . -name "*.test" -type f -delete; \
		find . -name ".env.test" -type f -delete; \
		find . -name ".env.backup" -type f -delete; \
		printf "\033[$(GREEN)m✅ Project cleaned\033[0m\n"; \
	fi

# Clean only log files
clean-logs:
	@if [ ! -d logs ]; then \
		printf "\033[$(GREEN)m✅ No log files to clean\033[0m\n"; \
	elif $(call ask_confirmation, This will remove all log files!, Cleaning log files, 📄); then \
		rm -rf logs/; \
		printf "\033[$(GREEN)m✅ Log files cleaned\033[0m\n"; \
	fi

# Clean the project without confirmation (for use in scripts/CI)
clean-all:
	@printf "\033[1;$(YELLOW)m🧹 Cleaning all project files without confirmation...\033[0m\n"
	@rm -rf bin/
	@rm -rf internal/docs/swaggerdocs/
	@rm -rf coverage/
	@rm -rf logs/
	@find . -name "*.test" -type f -delete
	@find . -name ".env.test" -type f -delete
	@find . -name ".env.backup" -type f -delete
	@printf "\033[$(GREEN)m✅ Project cleaned\033[0m\n"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@./scripts/test.sh --verbose

# Run only API-related tests
test-api:
	@echo "🧪 Running API tests..."
	@./scripts/test.sh --package "./internal/controller/..." --verbose

# Run only service-related tests
test-service:
	@echo "🧪 Running service tests..."
	@./scripts/test.sh --package "./internal/service/..." --verbose

# Run only repository-related tests
test-repository:
	@echo "🧪 Running repository tests..."
	@./scripts/test.sh --package "./internal/repository/..." --verbose

# Run tests with race detection
test-race:
	@echo "🧪 Running tests with race detection..."
	@./scripts/test.sh --race --verbose

# Test log file rotation
test-log-rotation:
	@echo "📋 Testing log file rotation..."
	@./scripts/test-log-rotation.sh

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
	@echo "🔄 Updating model map..."
	@go run ./cmd/model-mapper -sync
	@echo "✅ Model map updated"

# List available models for migration
migrate-list-models:
	@echo "🗃️ Available models for migrations:"
	@go run ./cmd/migrate -list-models

# Roll back the last migration
migrate-down:
	@if $(call ask_confirmation, This will roll back the last migration!, Rolling back migration, 🔙); then \
		go run ./cmd/migrate -down; \
		printf "\033[$(GREEN)m✅ Migration rollback completed\033[0m\n"; \
	fi

# Check migration status
migrate-status:
	@echo "🗃️ Checking migration status..."
	@go run ./cmd/migrate -version
	
# Reset all migrations
migrate-reset:
	@if $(call ask_confirmation, This will reset all migrations!, Resetting all migrations, 🔄); then \
		go run ./cmd/migrate -force 0; \
		printf "\033[$(GREEN)m✅ All migrations have been reset\033[0m\n"; \
	fi

# Create migrations from all available models (skip if table already exists)
migrate-all-models:
	@echo "🗃️ Creating migrations from all available models..."
	@go run ./cmd/migrate -all-models

# Run all seeders
seed:
	@echo "🌱 Running all database seeders..."
	@go run ./cmd/seed -all
	@echo "✅ Database seeding completed"

# Run only animal seeder
seed-animal:
	@echo "🐾 Running animal seeder..."
	@go run ./cmd/seed -seeder=animal
	@echo "✅ Animal seeding completed"

# Run only flower seeder
seed-flower:
	@echo "🌸 Running flower seeder..."
	@go run ./cmd/seed -seeder=flower
	@echo "✅ Flower seeding completed"

# Run seeders with custom count
seed-count:
	@echo "🌱 Running seeders with count: $(count)..."
	@if [ -z "$(count)" ]; then \
		echo "❌ Count is required. Usage: make seed-count count=100"; \
		exit 1; \
	fi
	@go run ./cmd/seed -all -count=$(count)
	@echo "✅ Database seeding completed with count: $(count)"

# Truncate table(s) based on model name
truncate:
	@if [ -z "$(model)" ]; then \
		echo "❌ Model name is required. Usage: make truncate model=animal"; \
		echo "Available models:"; \
		go run ./cmd/db -help | grep -A100 "Available models:" | grep "^  - " | sed 's/^  - //'; \
		exit 1; \
	fi
	@if $(call ask_confirmation, This will permanently delete ALL data from the $(model) table!, Truncating $(model) table, 🗑️); then \
		go run ./cmd/db -truncate $(model); \
		printf "\033[$(GREEN)m✅ Table truncated successfully\033[0m\n"; \
	fi

# Truncate all tables with confirmation
truncate-all:
	@if $(call ask_confirmation, DANGER: This will permanently delete ALL DATA from ALL TABLES!, Truncating all tables, ⚠️); then \
		go run ./cmd/db -truncate-all; \
		printf "\033[$(GREEN)m✅ All tables truncated successfully\033[0m\n"; \
	fi

# Update model map for database operations
update-model-map:
	@echo "🔄 Updating model map for database operations..."
	@go run ./cmd/model-mapper

# Clean model map (remove models that no longer exist)
clean-model-map:
	@echo "🧹 Cleaning model map (removing non-existent models)..."
	@go run ./cmd/model-mapper -clean-only

# Sync model map (add new models and remove deleted ones)
sync-model-map:
	@echo "🔄 Syncing model map (adding new models and removing deleted ones)..."
	@go run ./cmd/model-mapper -sync

# Add a new target to explicitly flush the Redis cache
flush-redis:
	$(call flush_redis_cache)

# Generate JWT token with default settings
generate-token:
	@echo "🔑 Generating JWT token with default settings..."
	@APP_ENV=$$(grep -E "^APP_ENV=" .env 2>/dev/null | cut -d= -f2 || echo "development"); \
	if [ "$$APP_ENV" != "development" ] && [ "$$APP_ENV" != "test" ] && [ -n "$$APP_ENV" ]; then \
		echo "❌ Error: Token generation is only available in development and test environments."; \
		echo "Current environment: $$APP_ENV"; \
		echo "Set APP_ENV to 'development' or 'test' to use this command."; \
		exit 1; \
	fi; \
	if [ ! -f ./bin/token-generator ] || [ ./cmd/token-generator/main.go -nt ./bin/token-generator ]; then \
		echo "🔨 Compiling token generator..."; \
		mkdir -p ./bin; \
		go build -o ./bin/token-generator ./cmd/token-generator; \
	fi; \
	JWT_SECRET=$$(grep -E "^JWT_SECRET=" .env 2>/dev/null | cut -d= -f2); \
	if [ -z "$$JWT_SECRET" ]; then JWT_SECRET="default-dev-secret"; fi; \
	./bin/token-generator --secret="$$JWT_SECRET"
	@echo "✅ JWT token generation complete"

# Generate JWT token with custom user ID
generate-token-user:
	@if [ -z "$(id)" ]; then \
		echo "❌ User ID is required. Usage: make generate-token-user id=123"; \
		exit 1; \
	fi
	@echo "🔑 Generating JWT token for user ID: $(id)..."
	@APP_ENV=$$(grep -E "^APP_ENV=" .env 2>/dev/null | cut -d= -f2 || echo "development"); \
	if [ "$$APP_ENV" != "development" ] && [ "$$APP_ENV" != "test" ] && [ -n "$$APP_ENV" ]; then \
		echo "❌ Error: Token generation is only available in development and test environments."; \
		echo "Current environment: $$APP_ENV"; \
		echo "Set APP_ENV to 'development' or 'test' to use this command."; \
		exit 1; \
	fi; \
	if [ ! -f ./bin/token-generator ] || [ ./cmd/token-generator/main.go -nt ./bin/token-generator ]; then \
		echo "🔨 Compiling token generator..."; \
		mkdir -p ./bin; \
		go build -o ./bin/token-generator ./cmd/token-generator; \
	fi; \
	JWT_SECRET=$$(grep -E "^JWT_SECRET=" .env 2>/dev/null | cut -d= -f2); \
	if [ -z "$$JWT_SECRET" ]; then JWT_SECRET="default-dev-secret"; fi; \
	./bin/token-generator --secret="$$JWT_SECRET" --id=$(id)
	@echo "✅ JWT token generation complete"

# Generate JWT token with admin role
generate-token-admin:
	@echo "🔑 Generating JWT token with admin role..."
	@APP_ENV=$$(grep -E "^APP_ENV=" .env 2>/dev/null | cut -d= -f2 || echo "development"); \
	if [ "$$APP_ENV" != "development" ] && [ "$$APP_ENV" != "test" ] && [ -n "$$APP_ENV" ]; then \
		echo "❌ Error: Token generation is only available in development and test environments."; \
		echo "Current environment: $$APP_ENV"; \
		echo "Set APP_ENV to 'development' or 'test' to use this command."; \
		exit 1; \
	fi; \
	if [ ! -f ./bin/token-generator ] || [ ./cmd/token-generator/main.go -nt ./bin/token-generator ]; then \
		echo "🔨 Compiling token generator..."; \
		mkdir -p ./bin; \
		go build -o ./bin/token-generator ./cmd/token-generator; \
	fi; \
	JWT_SECRET=$$(grep -E "^JWT_SECRET=" .env 2>/dev/null | cut -d= -f2); \
	if [ -z "$$JWT_SECRET" ]; then JWT_SECRET="default-dev-secret"; fi; \
	./bin/token-generator --secret="$$JWT_SECRET" --role=admin
	@echo "✅ JWT token generation complete"

# Force generate JWT token (works in any environment, for emergencies only)
generate-token-force:
	@if $(call ask_confirmation, Forcing token generation regardless of environment! This should only be used in emergencies., Generating emergency token); then \
		if [ ! -f ./bin/token-generator ] || [ ./cmd/token-generator/main.go -nt ./bin/token-generator ]; then \
			echo "🔨 Compiling token generator..."; \
			mkdir -p ./bin; \
			go build -o ./bin/token-generator ./cmd/token-generator; \
		fi; \
		JWT_SECRET=$$(grep -E "^JWT_SECRET=" .env 2>/dev/null | cut -d= -f2); \
		if [ -z "$$JWT_SECRET" ]; then JWT_SECRET="default-dev-secret"; fi; \
		./bin/token-generator --secret="$$JWT_SECRET" --force; \
		printf "\033[$(GREEN)m✅ Emergency token generation complete\033[0m\n"; \
	fi

# Aliases for common commands
s: swagger
su: swagger-ui
r: run
d: dev
t: test
l: lint
sd: seed
tr: truncate
tra: truncate-all
mam: migrate-all-models
um: update-model-map
cm: clean-model-map
sm: sync-model-map
smm: sync-model-map
fr: flush-redis
gt: generate-token
gtu: generate-token-user
gta: generate-token-admin
gtf: generate-token-force
tlr: test-log-rotation
cl: clean-logs
ca: clean-all

# Project Template Setup aliases
setup-s: setup
setup-g: setup-git
setup-f: setup-full

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

docker-clean: ## Remove containers and volumes defined in docker-compose.yml only (keep images)
	@if $(call ask_confirmation, This will remove containers and volumes from docker-compose.yml!, Cleaning project Docker resources); then \
		docker compose -f docker-compose.yml stop || true; \
		docker compose -f docker-compose.yml down -v --remove-orphans; \
		printf "\033[$(GREEN)m✅ Project Docker cleanup complete (images preserved)\033[0m\n"; \
	fi

# Docker command aliases
ddb: docker-db
dup: docker-up
ddown: docker-down
dps: docker-ps
dlogs: docker-logs
fps: fancy-ps

# Show environment variables
env-info: ## Show environment variables used by the application
	@if [ "$(show)" = "all" ]; then \
		./scripts/env-info.sh --show-all; \
	else \
		./scripts/env-info.sh; \
	fi

# Other aliases
ei: env-info

# Setup project as template
setup:
	@if [ -z "$(module)" ]; then \
		echo "❌ Module name is required. Usage: make setup module=github.com/yourusername/your-project"; \
		exit 1; \
	fi
	@echo "🔄 Setting up project with new module name: $(module)..."
	@go run ./cmd/setup-project -module $(module)

# Setup project with Git remote
setup-git:
	@if [ -z "$(module)" ]; then \
		echo "❌ Module name is required."; \
		exit 1; \
	fi
	@if [ -z "$(remote)" ]; then \
		echo "❌ Git remote URL is required."; \
		exit 1; \
	fi
	@echo "🔄 Setting up project with new module name and Git remote..."
	@go run ./cmd/setup-project -module $(module) -remote $(remote)

# Full setup with module name, Git remote, and new Git repository
setup-full:
	@if [ -z "$(module)" ]; then \
		echo "❌ Module name is required."; \
		exit 1; \
	fi
	@if [ -z "$(remote)" ]; then \
		echo "❌ Git remote URL is required."; \
		exit 1; \
	fi
	@echo "⚠️ This will perform a full project setup with:"
	@echo "  - Rename module to $(module)"
	@echo "  - Reset Git repository (remove .git folder and create new one)"
	@echo "  - Set Git remote to $(remote)"
	@echo ""
	@if $(call ask_confirmation, This will perform a complete project setup. All current Git history will be lost!, Performing full project setup, 🚀); then \
		echo "🔄 Running setup-project tool..."; \
		echo "y" | go run ./cmd/setup-project -module $(module) -remote $(remote) -reset-git -v; \
		printf "\033[$(GREEN)m✅ Project setup complete\033[0m\n"; \
	fi 