#!/bin/bash

# Set colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default settings
LOG_DIR="./logs"
LOG_FILE="$LOG_DIR/app.log"
TEST_DURATION=60  # seconds
LOG_INTERVAL=0.1  # seconds
MAX_SIZE_MB=1     # small size for testing
BACKUP_COUNT=3
MAX_AGE_DAYS=1

# Create temporary .env file with test settings
create_test_env() {
    echo "Creating temporary .env file for testing..."
    # Backup existing .env file if it exists
    if [ -f .env ]; then
        cp .env .env.backup
        echo -e "${YELLOW}Backed up existing .env file to .env.backup${NC}"
    fi

    # Create test .env file
    cat > .env.test << EOF
# Test environment configuration for log rotation
APP_ENV=development
PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT_PATH=stdout
LOG_FILE_PATH=$LOG_FILE
LOG_FILE_MAX_SIZE=$MAX_SIZE_MB
LOG_FILE_MAX_BACKUPS=$BACKUP_COUNT
LOG_FILE_MAX_AGE=$MAX_AGE_DAYS
LOG_FILE_COMPRESS=true
EOF

    # Use the test .env file
    cp .env.test .env
    echo -e "${GREEN}Created test .env file with log rotation settings${NC}"
}

# Restore original .env file
restore_env() {
    if [ -f .env.backup ]; then
        cp .env.backup .env
        rm .env.backup
        echo -e "${GREEN}Restored original .env file${NC}"
    else
        rm .env
        echo -e "${GREEN}Removed test .env file${NC}"
    fi
    rm .env.test
}

# Run the API with test settings
run_test() {
    echo -e "${BLUE}Starting API with log rotation settings:${NC}"
    echo -e "${YELLOW}  Log file: $LOG_FILE${NC}"
    echo -e "${YELLOW}  Max size: $MAX_SIZE_MB MB${NC}"
    echo -e "${YELLOW}  Backups: $BACKUP_COUNT${NC}"
    echo -e "${YELLOW}  Max age: $MAX_AGE_DAYS days${NC}"
    echo -e "${YELLOW}  Test duration: $TEST_DURATION seconds${NC}"
    
    # Start the API in background
    go run ./cmd/api &
    API_PID=$!
    
    # Generate lots of log entries
    echo -e "${BLUE}Generating log entries to trigger rotation...${NC}"
    end_time=$((SECONDS + TEST_DURATION))
    
    while [ $SECONDS -lt $end_time ]; do
        # Send a request to the API to generate logs
        curl -s "http://localhost:8080/health" > /dev/null
        curl -s "http://localhost:8080/api/v1/animals" > /dev/null
        sleep $LOG_INTERVAL
        
        # Show progress
        remaining=$((end_time - SECONDS))
        if [ $((remaining % 5)) -eq 0 ]; then
            echo -e "${YELLOW}Remaining time: $remaining seconds${NC}"
            # Check if rotation has occurred
            backup_count=$(ls -1 "$LOG_DIR" | grep -c "$LOG_FILE")
            if [ $backup_count -gt 1 ]; then
                echo -e "${GREEN}Log rotation occurred! Files in $LOG_DIR:${NC}"
                ls -lh "$LOG_DIR"
                break
            fi
        fi
    done
    
    # Stop the API
    echo -e "${BLUE}Stopping API...${NC}"
    kill $API_PID
    wait $API_PID 2>/dev/null
    
    # Show final log files
    echo -e "${GREEN}Final log files in $LOG_DIR:${NC}"
    ls -lh "$LOG_DIR"
}

# Main function
main() {
    echo -e "${BLUE}=== Log Rotation Test Script ===${NC}"
    
    # Create log directory if it doesn't exist
    mkdir -p "$LOG_DIR"
    
    # Setup test environment
    create_test_env
    
    # Run the test
    run_test
    
    # Cleanup
    restore_env
    
    echo -e "${GREEN}Test completed!${NC}"
    echo -e "${YELLOW}Log files were retained in $LOG_DIR for inspection.${NC}"
    echo -e "${YELLOW}You can delete them manually when no longer needed.${NC}"
}

# Run the main function
main 