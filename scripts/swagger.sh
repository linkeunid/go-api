#!/bin/bash
set -e

# Load environment variables from .env if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Default API settings if not provided
API_HOST=${API_HOST:-0.0.0.0}
API_PORT=${PORT:-8080}

# Use localhost for better accessibility when host is 0.0.0.0
if [ "$API_HOST" = "0.0.0.0" ]; then
  SWAGGER_API_HOST="localhost" 
else
  SWAGGER_API_HOST="$API_HOST"
fi

echo "Generating Swagger docs for API at $SWAGGER_API_HOST:$API_PORT..."

# Create directories if they don't exist
mkdir -p internal/docs/swaggerdocs

# Replace the host annotation in swagger.go
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  sed -i '' "s|// @host .*|// @host $SWAGGER_API_HOST:$API_PORT|" internal/docs/swagger.go
else
  # Linux and others
  sed -i "s|// @host .*|// @host $SWAGGER_API_HOST:$API_PORT|" internal/docs/swagger.go
fi

# Generate Swagger docs directly from the swagger.go file that contains the annotations
swag init -g internal/docs/swagger.go --output ./internal/docs/swaggerdocs

echo "✅ Swagger docs generated at ./internal/docs/swaggerdocs" 