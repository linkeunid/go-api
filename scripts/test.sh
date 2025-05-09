#!/bin/bash

# Set colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default settings
COVERAGE_DIR="./coverage"
COVERAGE_PROFILE="coverage.out"
COVERAGE_HTML="coverage.html"
TEST_ARGS=""
PACKAGE_PATTERN="./..."
SHOW_COVERAGE=true
GENERATE_HTML=true

# Function to print usage
function print_usage {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -p, --package PATTERN   Package pattern to test (default: ./...)"
    echo "  -c, --coverage DIR      Directory to store coverage reports (default: ./coverage)"
    echo "  --no-coverage           Don't display coverage"
    echo "  --no-html               Don't generate HTML coverage report"
    echo "  --race                  Enable race detection"
    echo "  --verbose               Enable verbose output"
    echo "  -h, --help              Display this help message"
    echo ""
    echo "Example: $0 --package './internal/...' --race"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -p|--package)
            PACKAGE_PATTERN="$2"
            shift
            shift
            ;;
        -c|--coverage)
            COVERAGE_DIR="$2"
            shift
            shift
            ;;
        --no-coverage)
            SHOW_COVERAGE=false
            shift
            ;;
        --no-html)
            GENERATE_HTML=false
            shift
            ;;
        --race)
            TEST_ARGS="$TEST_ARGS -race"
            shift
            ;;
        --verbose)
            TEST_ARGS="$TEST_ARGS -v"
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown option $key${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# Create coverage directory if it doesn't exist
mkdir -p $COVERAGE_DIR

# Full path to coverage files
COVERAGE_PROFILE_PATH="$COVERAGE_DIR/$COVERAGE_PROFILE"
COVERAGE_HTML_PATH="$COVERAGE_DIR/$COVERAGE_HTML"

echo -e "${BLUE}Running tests for packages: $PACKAGE_PATTERN${NC}"
echo -e "${BLUE}Test arguments: $TEST_ARGS${NC}"

# Run tests with coverage
echo -e "${YELLOW}Running tests with coverage...${NC}"
go test $TEST_ARGS -coverprofile=$COVERAGE_PROFILE_PATH $PACKAGE_PATTERN

# Check if tests passed
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Tests passed successfully!${NC}"
else
    echo -e "${RED}Tests failed!${NC}"
    exit 1
fi

# Show coverage statistics if enabled
if [ "$SHOW_COVERAGE" = true ]; then
    echo -e "${YELLOW}Coverage statistics:${NC}"
    go tool cover -func=$COVERAGE_PROFILE_PATH
fi

# Generate HTML coverage report if enabled
if [ "$GENERATE_HTML" = true ]; then
    echo -e "${YELLOW}Generating HTML coverage report...${NC}"
    go tool cover -html=$COVERAGE_PROFILE_PATH -o=$COVERAGE_HTML_PATH
    echo -e "${GREEN}HTML coverage report generated at $COVERAGE_HTML_PATH${NC}"
fi

# Show total coverage percentage
if [ "$SHOW_COVERAGE" = true ]; then
    TOTAL_COVERAGE=$(go tool cover -func=$COVERAGE_PROFILE_PATH | grep total | awk '{print $3}')
    echo -e "${GREEN}Total coverage: $TOTAL_COVERAGE${NC}"
fi

echo -e "${GREEN}All done!${NC}" 