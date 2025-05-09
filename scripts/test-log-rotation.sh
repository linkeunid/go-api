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
PORT=8080

# API process ID
API_PID=""

# Check if port is already in use
check_port() {
    local pid=$(lsof -t -i:$PORT 2>/dev/null)
    if [ -n "$pid" ]; then
        local process_info=$(ps -p $pid -o comm= 2>/dev/null)
        echo -e "${YELLOW}⚠️ Port $PORT is already in use by process $pid ($process_info)${NC}"
        echo -e "${YELLOW}Options:${NC}"
        echo -e "  ${YELLOW}1) Kill the process and continue${NC}"
        echo -e "  ${YELLOW}2) Exit${NC}"
        read -p "Enter your choice (1/2): " choice

        case "$choice" in
            1)
                echo -e "${YELLOW}Killing process $pid...${NC}"
                kill -9 $pid
                sleep 1
                if lsof -t -i:$PORT &>/dev/null; then
                    echo -e "${RED}❌ Failed to kill process. Please stop it manually and try again.${NC}"
                    cleanup
                else
                    echo -e "${GREEN}✅ Process killed successfully${NC}"
                fi
                ;;
            *)
                echo -e "${YELLOW}Exiting without running the test.${NC}"
                exit 0
                ;;
        esac
    else
        echo -e "${GREEN}✅ Port $PORT is available${NC}"
    fi
}

# Cleanup function to ensure .env is restored and processes are killed
cleanup() {
    echo -e "\n${YELLOW}Caught termination signal. Cleaning up...${NC}"
    
    # Kill API process if running
    if [ -n "$API_PID" ] && ps -p $API_PID > /dev/null; then
        echo -e "${YELLOW}Stopping API process...${NC}"
        kill -15 $API_PID 2>/dev/null

        # Wait for process to terminate gracefully
        for i in {1..3}; do
            if ! ps -p $API_PID > /dev/null; then
                break
            fi
            sleep 1
        done

        # Force kill if still running
        if ps -p $API_PID > /dev/null; then
            echo -e "${YELLOW}⚠️ API process did not terminate gracefully, forcing...${NC}"
            kill -9 $API_PID 2>/dev/null
        fi

        wait $API_PID 2>/dev/null || true
    fi

    # Double-check no processes are using our port
    local pid=$(lsof -t -i:$PORT 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}⚠️ Port $PORT is still in use by process $pid. Attempting to kill...${NC}"
        kill -9 $pid 2>/dev/null
        sleep 1
    fi
    
    # Restore original .env file
    if [ -f .env.backup ]; then
        echo -e "${YELLOW}Restoring original .env file...${NC}"
        cp .env.backup .env
        rm .env.backup
        echo -e "${GREEN}Restored original .env file${NC}"
    fi
    
    # Clean up test files
    if [ -f .env.test ]; then
        rm .env.test
    fi
    
    echo -e "${YELLOW}Cleanup complete. Exiting.${NC}"
    exit 1
}

# Set up trap to handle interruption signals
trap cleanup SIGINT SIGTERM SIGHUP

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
PORT=$PORT
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
    if [ -f .env.test ]; then
        rm .env.test
    fi
}

# Run the API with test settings
run_test() {
    echo -e "${BLUE}Starting API with log rotation settings:${NC}"
    echo -e "${YELLOW}  Log file: $LOG_FILE${NC}"
    echo -e "${YELLOW}  Max size: $MAX_SIZE_MB MB${NC}"
    echo -e "${YELLOW}  Backups: $BACKUP_COUNT${NC}"
    echo -e "${YELLOW}  Max age: $MAX_AGE_DAYS days${NC}"
    echo -e "${YELLOW}  Test duration: $TEST_DURATION seconds${NC}"
    echo -e "${YELLOW}  (Press Ctrl+C to stop the test at any time)${NC}"
    
    # Start the API in background
    go run ./cmd/api &
    API_PID=$!
    
    # Wait a moment for API to start
    echo -e "${YELLOW}Waiting for API to start...${NC}"
    sleep 2

    # Check if API started successfully
    if ! ps -p $API_PID > /dev/null; then
        echo -e "${RED}❌ API failed to start. Check logs for details:${NC}"
        tail -n 10 $LOG_FILE
        cleanup
    fi

    # Generate lots of log entries
    echo -e "${BLUE}Generating log entries to trigger rotation...${NC}"
    end_time=$((SECONDS + TEST_DURATION))
    
    while [ $SECONDS -lt $end_time ]; do
        # Send a request to the API to generate logs
        curl -s "http://localhost:$PORT/health" > /dev/null
        curl -s "http://localhost:$PORT/api/v1/animals" > /dev/null
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
    if [ -n "$API_PID" ] && ps -p $API_PID > /dev/null; then
        kill -15 $API_PID 2>/dev/null

        # Wait for process to terminate gracefully
        for i in {1..5}; do
            if ! ps -p $API_PID > /dev/null; then
                break
            fi
            sleep 1
        done

        # Force kill if still running
        if ps -p $API_PID > /dev/null; then
            echo -e "${YELLOW}⚠️ API process did not terminate gracefully, forcing...${NC}"
            kill -9 $API_PID 2>/dev/null
        fi

        wait $API_PID 2>/dev/null || true
        echo -e "${GREEN}✅ API process terminated${NC}"
    else
        echo -e "${YELLOW}⚠️ API process was not running${NC}"
    fi
    API_PID=""
    
    # Double-check no processes are using our port
    local pid=$(lsof -t -i:$PORT 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}⚠️ Port $PORT is still in use by process $pid. Attempting to kill...${NC}"
        kill -9 $pid 2>/dev/null
        sleep 1
    fi

    # Show final log files
    echo -e "${GREEN}Final log files in $LOG_DIR:${NC}"
    ls -lh "$LOG_DIR"
}

# Main function
main() {
    echo -e "${BLUE}=== Log Rotation Test Script ===${NC}"
    
    # Create log directory if it doesn't exist
    mkdir -p "$LOG_DIR"
    
    # Check if port is available
    check_port

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