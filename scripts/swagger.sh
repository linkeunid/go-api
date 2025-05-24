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
  # Export only valid environment variables, stripping inline comments and empty lines
  set -a # Enable auto-export of variables
  source <(grep -v '^#' .env | grep -E '^[A-Z_][A-Z0-9_]*=' | sed 's/#.*//' | sed 's/^[[:space:]]*//' | sed '/^$/d')
  set +a # Disable auto-export
  echo -e "${CYAN}🔍 Found .env file, using environment variables...${NC}"
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

echo -e "${MAGENTA}📚 Generating Swagger docs for API at ${BOLD}$SWAGGER_API_HOST:$API_PORT${NC}..."

# Create directories if they don't exist
mkdir -p internal/docs/swaggerdocs
echo -e "${BLUE}📁 Created docs directories${NC}"

# Replace the host annotation in swagger.go
echo -e "${YELLOW}🔄 Updating API host in swagger annotations...${NC}"
if [[ "$OSTYPE" == "darwin"* ]]; then
  # Mac OSX
  sed -i '' "s|// @host .*|// @host $SWAGGER_API_HOST:$API_PORT|" internal/docs/swagger.go
else
  # Linux and others
  sed -i "s|// @host .*|// @host $SWAGGER_API_HOST:$API_PORT|" internal/docs/swagger.go
fi

# Generate Swagger docs directly from the swagger.go file that contains the annotations
echo -e "${BLUE}⚙️ Running swag init...${NC}"
swag init -g internal/docs/swagger.go --output ./internal/docs/swaggerdocs

echo -e "${GREEN}✅ Swagger docs generated at ${BOLD}./internal/docs/swaggerdocs${NC}"
echo -e "${CYAN}💡 Run ${BOLD}./scripts/swagger-ui.sh${NC} ${CYAN}to start Swagger UI server${NC}" 