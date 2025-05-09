#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Load environment variables from .env if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
  echo -e "${CYAN}🔍 Found .env file, using environment variables...${NC}"
fi

# Default to port 8090 for Swagger UI
SWAGGER_PORT=${SWAGGER_PORT:-8090}
SWAGGER_HOST=${SWAGGER_HOST:-localhost}

# API settings - where the actual service is running
API_HOST=${API_HOST:-0.0.0.0}
API_PORT=${PORT:-8080}

echo -e "${MAGENTA}🚀 Starting Swagger UI server on ${BOLD}http://${SWAGGER_HOST}:${SWAGGER_PORT}/swagger/${NC}"
echo -e "${BLUE}📘 This is a standalone UI for OpenAPI documentation${NC}"
echo -e "${YELLOW}🔗 Configured to use API at ${BOLD}http://${API_HOST}:${API_PORT}/${NC}"
echo -e "${CYAN}💡 Press Ctrl+C to exit${NC}"
go run ./test/swagger-test.go ${SWAGGER_HOST} ${SWAGGER_PORT} ${API_HOST} ${API_PORT} 