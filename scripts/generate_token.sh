#!/bin/bash
set -e

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Compile the token generator if not already compiled
TOKEN_GENERATOR_BIN="$ROOT_DIR/bin/token-generator"
if [ ! -f "$TOKEN_GENERATOR_BIN" ] || [ "$ROOT_DIR/cmd/token-generator/main.go" -nt "$TOKEN_GENERATOR_BIN" ]; then
    echo "🔨 Compiling token generator..."
    mkdir -p "$ROOT_DIR/bin"
    go build -o "$TOKEN_GENERATOR_BIN" "$ROOT_DIR/cmd/token-generator"
fi

# Run the token generator with all arguments passed to this script
"$TOKEN_GENERATOR_BIN" "$@" 