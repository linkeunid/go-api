#!/bin/bash
set -e

# Load environment variables from .env if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Default to port 8090 for Swagger UI
SWAGGER_PORT=${SWAGGER_PORT:-8090}
SWAGGER_HOST=${SWAGGER_HOST:-localhost}

# API settings - where the actual service is running
API_HOST=${API_HOST:-0.0.0.0}
API_PORT=${PORT:-8080}

echo "Starting Swagger UI server on http://${SWAGGER_HOST}:${SWAGGER_PORT}/swagger/"
echo "This is a standalone UI for OpenAPI documentation"
echo "Configured to use API at http://${API_HOST}:${API_PORT}/"
go run ./test/swagger-test.go ${SWAGGER_HOST} ${SWAGGER_PORT} ${API_HOST} ${API_PORT} 