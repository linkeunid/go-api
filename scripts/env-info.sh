#!/bin/bash

# env-info.sh - Environment Configuration Display Script
# This script displays .env file contents and parsed environment variables
# Usage: ./scripts/env-info.sh [--show-all]
#   --show-all    Show all values including sensitive ones (USE WITH CAUTION)

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Default behavior: hide sensitive values
SHOW_ALL=false

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --show-all)
                SHOW_ALL=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Show help information
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Display environment configuration with security-conscious defaults."
    echo ""
    echo "OPTIONS:"
    echo "  --show-all    Show all values including sensitive ones (USE WITH CAUTION)"
    echo "  --help, -h    Show this help message"
    echo ""
    echo "EXAMPLES:"
    echo "  $0                    # Show config with sensitive values hidden (default)"
    echo "  $0 --show-all         # Show all values including sensitive ones"
    echo ""
    echo "SECURITY NOTE:"
    echo "  By default, sensitive values (SECRET, PASSWORD, TOKEN, KEY, USER) are hidden."
    echo "  Use --show-all only in secure environments and avoid logging the output."
}

# Helper function to get environment variable with default
get_env() {
    local var_name="$1"
    local default_value="$2"
    
    if [ -f .env ]; then
        local value=$(grep -E "^${var_name}=" .env 2>/dev/null | cut -d= -f2 || echo "")
        if [ -n "$value" ]; then
            echo "$value"
        else
            echo "$default_value"
        fi
    else
        echo "$default_value"
    fi
}

# Main function
main() {
    echo ""
    echo -e "✨ 🔍 ${BOLD}${PURPLE}Environment Configuration${NC} 🔍 ✨"
    echo ""
    echo -e "${BOLD}${CYAN}┌───────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│ 📄 .ENV FILE CONTENTS                             │${NC}"
    echo -e "${BOLD}${CYAN}└───────────────────────────────────────────────────┘${NC}"
    
    if [ -f .env ]; then
        echo -e "${BOLD}${GREEN}✅ .env file found${NC}"
        echo ""
        
        while IFS= read -r line; do
            if [ -n "$line" ] && [ "${line#\#}" = "$line" ]; then
                # Non-comment, non-empty line
                if [ "$SHOW_ALL" = false ] && echo "$line" | grep -q "SECRET\|PASSWORD\|TOKEN\|KEY\|USER"; then
                    VAR_NAME=$(echo "$line" | cut -d= -f1)
                    echo -e "   ${BOLD}${YELLOW}${VAR_NAME}${NC}=${BOLD}${RED}[HIDDEN]${NC}"
                else
                    echo -e "   ${BOLD}${YELLOW}${line}${NC}"
                fi
            elif [ "${line#\#}" != "$line" ]; then
                # Comment line
                echo -e "   ${GRAY}${line}${NC}"
            fi
        done < .env
    else
        echo -e "${BOLD}${RED}❌ .env file not found${NC}"
        echo "   Create a .env file to configure environment variables"
    fi
    
    echo ""
    echo -e "${BOLD}${CYAN}┌───────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│ 🌐 PARSED ENVIRONMENT VARIABLES                   │${NC}"
    echo -e "${BOLD}${CYAN}└───────────────────────────────────────────────────┘${NC}"
    echo ""
    
    echo -e "${BOLD}${BLUE}🚀 Application Settings${NC}"
    echo -e "   APP_ENV:        ${BOLD}${GREEN}$(get_env "APP_ENV" "development")${NC}"
    echo -e "   API_HOST:       ${BOLD}${GREEN}$(get_env "API_HOST" "localhost")${NC}"
    echo -e "   API_PORT:       ${BOLD}${GREEN}$(get_env "API_PORT" "8080")${NC}"
    echo ""
    
    echo -e "${BOLD}${BLUE}🗃️ Database Configuration${NC}"
    echo -e "   DB_HOST:        ${BOLD}${GREEN}$(get_env "DB_HOST" "localhost")${NC}"
    echo -e "   DB_PORT:        ${BOLD}${GREEN}$(get_env "DB_PORT" "3306")${NC}"
    if [ "$SHOW_ALL" = true ]; then
        echo -e "   DB_USER:        ${BOLD}${GREEN}$(get_env "DB_USER" "root")${NC}"
    else
        echo -e "   DB_USER:        ${BOLD}${RED}[HIDDEN]${NC}"
    fi
    echo -e "   DB_NAME:        ${BOLD}${GREEN}$(get_env "DB_NAME" "linkeun_go_api")${NC}"
    echo ""
    
    echo -e "${BOLD}${BLUE}🔴 Redis Configuration${NC}"
    echo -e "   REDIS_HOST:     ${BOLD}${GREEN}$(get_env "REDIS_HOST" "localhost")${NC}"
    echo -e "   REDIS_PORT:     ${BOLD}${GREEN}$(get_env "REDIS_PORT" "6379")${NC}"
    echo -e "   REDIS_ENABLED:  ${BOLD}${GREEN}$(get_env "REDIS_ENABLED" "false")${NC}"
    echo ""
    
    echo -e "${BOLD}${BLUE}🔐 Authentication Settings${NC}"
    echo -e "   AUTH_ENABLED:   ${BOLD}${GREEN}$(get_env "AUTH_ENABLED" "false")${NC}"
    echo -e "   JWT_EXPIRATION: ${BOLD}${GREEN}$(get_env "JWT_EXPIRATION" "24h")${NC}"
    echo -e "   JWT_ALLOWED_ISSUERS: ${BOLD}${GREEN}$(get_env "JWT_ALLOWED_ISSUERS" "linkeun-go-api")${NC}"
    echo ""
    
    echo -e "${BOLD}${CYAN}┌───────────────────────────────────────────────────┐${NC}"
    echo -e "${BOLD}${CYAN}│ 📝 NOTES                                          │${NC}"
    echo -e "${BOLD}${CYAN}└───────────────────────────────────────────────────┘${NC}"
    echo -e "   ${BOLD}${YELLOW}📋 Values shown are actual values from .env or defaults${NC}"
    if [ "$SHOW_ALL" = true ]; then
        echo -e "   ${BOLD}${RED}⚠️ ALL VALUES ARE VISIBLE (--show-all flag used)${NC}"
        echo -e "   ${BOLD}${RED}🚨 Ensure this output is not logged or shared publicly${NC}"
    else
        echo -e "   ${BOLD}${RED}🔒 Sensitive values (SECRET, PASSWORD, TOKEN, KEY, USER) are hidden${NC}"
        echo -e "   ${BOLD}${YELLOW}💡 Use --show-all flag to reveal sensitive values (USE WITH CAUTION)${NC}"
    fi
    echo -e "   ${BOLD}${GREEN}💡 Use 'make env-info > env-info.txt' to save output to file${NC}"
    echo ""
}

# Run main function
parse_args "$@"
main 