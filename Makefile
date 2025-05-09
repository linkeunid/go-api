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
	@printf "  \033[1mmake dev\033[0m            - 🔄 Run development server with hot reload\n"
	@printf "  \033[1mmake run\033[0m            - 🚀 Run the application\n"
	@printf "  \033[1mmake build\033[0m          - 🔨 Build the application\n"
	@printf "\n"
	@printf "\033[1;36m🐳 Docker Commands\033[0m\n"
	@printf "  \033[1mmake docker-up\033[0m      - 🚀 Start all containers\n"
	@printf "  \033[1mmake docker-down\033[0m    - 🛑 Stop all containers\n"
	@printf "  \033[1mmake docker-db\033[0m      - 🗃️ Start only database containers (MySQL and Redis)\n"
	@printf "  \033[1mmake docker-rebuild\033[0m - 🔄 Rebuild and restart all containers\n"
	@printf "  \033[1mmake fancy-ps\033[0m       - 🌈 Show fancy container status with colors and details\n"
	@printf "  \033[1mmake docker-logs\033[0m    - 📋 View logs from all containers\n"
	@printf "  \033[1mmake docker-ps\033[0m      - 📊 List running containers\n"
	@printf "  \033[1mmake docker-clean\033[0m   - 🧹 Remove all containers, volumes, and images\n"
	@printf "\n"
	@printf "\033[1;36m🗃️ Database Migrations\033[0m\n"
	@printf "  \033[1mmake migrate\033[0m        - 📊 Run database migrations\n"
	@printf "  \033[1mmake migrate-status\033[0m - 📊 Show current migration status\n"
	@printf "  \033[1mmake migrate-down\033[0m   - ⏮️ Roll back the last migration\n"
	@printf "  \033[1mmake migrate-create name=NAME\033[0m - 📝 Create a new migration\n"
	@printf "  \033[1mmake migrate-from-model model=NAME\033[0m - 🔄 Create a migration from a model\n"
	@printf "  \033[1mmake migrate-list-models\033[0m - 📋 List available models for migrations\n"
	@printf "  \033[1mmake migrate-reset\033[0m  - 🔄 Reset all migrations\n"
	@printf "\n"
	@printf "\033[1;36m🌱 Database Seeders\033[0m\n"
	@printf "  \033[1mmake seed\033[0m           - 🌱 Run all database seeders\n"
	@printf "  \033[1mmake seed-count\033[0m     - 🔢 Run seeders with custom count (e.g., make seed-count count=100)\n"
	@printf "  \033[1mmake seed-animal\033[0m    - 🐾 Run only animal seeder\n"
	@printf "  \033[1mmake seed-flower\033[0m    - 🌸 Run only flower seeder\n"
	@printf "\n"
	@printf "\033[1;36m💾 Database Operations\033[0m\n"
	@printf "  \033[1mmake update-model-map\033[0m - 🔄 Update model map for database operations\n"
	@printf "  \033[1mmake sync-model-map\033[0m - 🔄 Sync model map by adding new models and removing deleted ones\n"
	@printf "  \033[1mmake clean-model-map\033[0m - 🧹 Remove models from the map that no longer exist\n"
	@printf "  \033[1mmake truncate model=NAME\033[0m - 🗑️ Truncate specific table with confirmation\n"
	@printf "  \033[1mmake truncate-all\033[0m   - 🧹 Truncate all tables with confirmation\n"
	@printf "\n"
	@printf "\033[1;36m🧪 Testing & Quality\033[0m\n"
	@printf "  \033[1mmake test\033[0m           - 🧪 Run tests\n"
	@printf "  \033[1mmake fmt\033[0m            - ✨ Format code\n"
	@printf "  \033[1mmake lint\033[0m           - 🔍 Lint code\n"
	@printf "  \033[1mmake test-coverage\033[0m  - 📊 Run tests with coverage report\n"
	@printf "  \033[1mmake test-log-rotation\033[0m - 📋 Test log file rotation functionality\n"
	@printf "  \033[1mmake mocks\033[0m          - 🧩 Generate mocks for testing\n"
	@printf "\n"
	@printf "\033[1;36m📚 Documentation\033[0m\n"
	@printf "  \033[1mmake swagger\033[0m        - 📝 Generate Swagger documentation\n"
	@printf "  \033[1mmake swagger-ui\033[0m     - 🌐 Run Swagger UI server\n"
	@printf "\n"
	@printf "\033[1;36m🔐 JWT Authentication\033[0m\n"
	@printf "  \033[1mmake generate-token\033[0m      - 🔑 Generate JWT token with default settings\n"
	@printf "  \033[1mmake generate-token-user id=123\033[0m - 👤 Generate JWT token for specific user ID\n"
	@printf "  \033[1mmake generate-token-admin\033[0m - 👑 Generate JWT token with admin role\n"
	@printf "  \033[1mmake generate-token-force\033[0m - ⚡ Force token generation (works in any environment)\n"
	@printf "\n"
	@printf "\033[1;36m🔧 Project Management\033[0m\n"
	@printf "  \033[1mmake init\033[0m           - 🔧 Initialize the project\n"
	@printf "  \033[1mmake env-info\033[0m       - ℹ️ Show environment variables used by the application\n"
	@printf "  \033[1mmake clean\033[0m          - 🧹 Clean build artifacts (with confirmation)\n"
	@printf "  \033[1mmake clean-all\033[0m      - 🧹 Clean build artifacts (no confirmation, for CI/scripts)\n"
	@printf "  \033[1mmake clean-logs\033[0m     - 🧹 Clean only log files\n"
	@printf "  \033[1mmake flush-redis\033[0m    - 🧹 Explicitly flush Redis cache\n"
	@printf "\n"
	@printf "\033[1;36m🔧 Project Template Setup\033[0m\n"
	@printf "  \033[1mmake setup module=MODULE_NAME\033[0m - 🛠️ Setup project with new module name\n"
	@printf "  \033[1mmake setup-git remote=GIT_URL\033[0m - 🔄 Setup project with new module name and git remote\n"
	@printf "  \033[1mmake setup-full module=MODULE_NAME remote=GIT_URL\033[0m - 🚀 Full setup with new git repo\n"
	@printf "\n"
	@printf "\033[1;36m⚡ Command Aliases\033[0m\n"
	@printf "  \033[1mmake d\033[0m              - ↩️ Alias for 'dev'\n"
	@printf "  \033[1mmake r\033[0m              - ↩️ Alias for 'run'\n"
	@printf "  \033[1mmake s\033[0m              - ↩️ Alias for 'swagger'\n"
	@printf "  \033[1mmake su\033[0m             - ↩️ Alias for 'swagger-ui'\n"
	@printf "  \033[1mmake t\033[0m              - ↩️ Alias for 'test'\n"
	@printf "  \033[1mmake l\033[0m              - ↩️ Alias for 'lint'\n"
	@printf "  \033[1mmake sd\033[0m             - ↩️ Alias for 'seed'\n"
	@printf "  \033[1mmake tr\033[0m             - ↩️ Alias for 'truncate'\n"
	@printf "  \033[1mmake tra\033[0m            - ↩️ Alias for 'truncate-all'\n"
	@printf "  \033[1mmake um\033[0m             - ↩️ Alias for 'update-model-map'\n"
	@printf "  \033[1mmake cm\033[0m             - ↩️ Alias for 'clean-model-map'\n"
	@printf "  \033[1mmake sm\033[0m             - ↩️ Alias for 'sync-model-map'\n"
	@printf "  \033[1mmake fr\033[0m             - ↩️ Alias for 'flush-redis'\n"
	@printf "  \033[1mmake gt\033[0m             - ↩️ Alias for 'generate-token'\n"
	@printf "  \033[1mmake gtu\033[0m            - ↩️ Alias for 'generate-token-user'\n"
	@printf "  \033[1mmake gta\033[0m            - ↩️ Alias for 'generate-token-admin'\n"
	@printf "  \033[1mmake gtf\033[0m            - ↩️ Alias for 'generate-token-force'\n"
	@printf "  \033[1mmake tlr\033[0m            - ↩️ Alias for 'test-log-rotation'\n"
	@printf "  \033[1mmake cl\033[0m             - ↩️ Alias for 'clean-logs'\n"
	@printf "  \033[1mmake ca\033[0m             - ↩️ Alias for 'clean-all'\n"
	@printf "  \033[1mmake setup-s\033[0m        - ↩️ Alias for 'setup'\n"
	@printf "  \033[1mmake setup-g\033[0m        - ↩️ Alias for 'setup-git'\n"
	@printf "  \033[1mmake setup-f\033[0m        - ↩️ Alias for 'setup-full'\n"
	@printf "  \033[1mmake ddb\033[0m            - ↩️ Alias for 'docker-db'\n"
	@printf "  \033[1mmake dup\033[0m            - ↩️ Alias for 'docker-up'\n"
	@printf "  \033[1mmake ddown\033[0m          - ↩️ Alias for 'docker-down'\n"
	@printf "  \033[1mmake dps\033[0m            - ↩️ Alias for 'docker-ps'\n"
	@printf "  \033[1mmake dlogs\033[0m          - ↩️ Alias for 'docker-logs'\n"
	@printf "  \033[1mmake fps\033[0m            - ↩️ Alias for 'fancy-ps'\n"
	@printf "  \033[1mmake ei\033[0m             - ↩️ Alias for 'env-info'\n"

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
um: update-model-map
cm: clean-model-map
sm: sync-model-map
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

docker-clean: ## Remove all containers, volumes, and images
	@if $(call ask_confirmation, This will remove ALL Docker containers\\, volumes\\, and images!, Cleaning Docker resources); then \
		docker compose -f docker-compose.yml down -v; \
		docker system prune -af --volumes; \
		printf "\033[$(GREEN)m✅ Docker cleanup complete\033[0m\n"; \
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
	@echo "   AUTH_ENABLED: $(call get_env,AUTH_ENABLED,false)"
	@echo "   JWT_EXPIRATION: $(call get_env,JWT_EXPIRATION,24h)"
	@echo "   JWT_ALLOWED_ISSUERS: $(call get_env,JWT_ALLOWED_ISSUERS,linkeun-go-api)"
	@echo "📝 Note: Values shown are actual values from .env or defaults if not defined"
	@echo "🔒 Note: JWT_SECRET is not displayed for security reasons"

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
	@if $(call ask_confirmation, This will rename the module\\, reset the Git repository\\, and set a new remote!, Performing full project setup); then \
		go run ./cmd/setup-project -module $(module) -remote $(remote) -reset-git -v; \
		printf "\033[$(GREEN)m✅ Project setup complete\033[0m\n"; \
	fi 