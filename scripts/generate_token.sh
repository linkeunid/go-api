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

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Compile the token generator if not already compiled
TOKEN_GENERATOR_BIN="$ROOT_DIR/bin/token-generator"
if [ ! -f "$TOKEN_GENERATOR_BIN" ] || [ "$ROOT_DIR/cmd/token-generator/main.go" -nt "$TOKEN_GENERATOR_BIN" ]; then
    echo -e "${BLUE}🔨 Compiling token generator...${NC}"
    mkdir -p "$ROOT_DIR/bin"
    go build -o "$TOKEN_GENERATOR_BIN" "$ROOT_DIR/cmd/token-generator"
    echo -e "${GREEN}✅ Token generator compiled successfully!${NC}"
fi

# Run the token generator with all arguments passed to this script
echo -e "${MAGENTA}🔑 Generating authentication token...${NC}"
"$TOKEN_GENERATOR_BIN" "$@"
echo -e "${GREEN}✨ Token generation completed!${NC}" 