#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Load environment variables from .env if it exists
if [ -f .env ]; then
  # Export only valid environment variables, stripping inline comments and empty lines
  set -a # Enable auto-export of variables
  source <(grep -v '^#' .env | grep -E '^[A-Z_][A-Z0-9_]*=' | sed 's/#.*//' | sed 's/^[[:space:]]*//' | sed '/^$/d')
  set +a # Disable auto-export
  echo -e "${CYAN}🔍 Found .env file, using environment variables...${NC}"
fi

# Default to port 8090 for Swagger UI
SWAGGER_PORT=${SWAGGER_PORT:-8090}
SWAGGER_HOST=${SWAGGER_HOST:-localhost}

# API settings - where the actual service is running
API_HOST=${API_HOST:-localhost}
API_PORT=${PORT:-8080}

# Check if API is running
if ! command -v lsof > /dev/null; then
  echo -e "${YELLOW}⚠️ Warning: lsof command not found, skipping API check${NC}"
else
  if ! lsof -i:${API_PORT} > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️ Warning: API service doesn't appear to be running on ${API_HOST}:${API_PORT}${NC}"
    echo -e "${YELLOW}⚠️ Some Swagger features may not work correctly${NC}"
    echo -e "${YELLOW}⚠️ Consider running 'make dev' in another terminal first${NC}"
    echo ""
  else
    echo -e "${GREEN}✅ API service detected on ${API_HOST}:${API_PORT}${NC}"
  fi
fi

echo -e "${MAGENTA}🚀 Starting Swagger UI server on ${BOLD}http://${SWAGGER_HOST}:${SWAGGER_PORT}/swagger/${NC}"
echo -e "${BLUE}📘 This is a standalone UI for OpenAPI documentation${NC}"
echo -e "${YELLOW}🔗 Configured to use API at ${BOLD}http://${API_HOST}:${API_PORT}/${NC}"
echo -e "${CYAN}💡 Press Ctrl+C to exit${NC}"
go run ./test/swagger-test.go ${SWAGGER_HOST} ${SWAGGER_PORT} ${API_HOST} ${API_PORT} 