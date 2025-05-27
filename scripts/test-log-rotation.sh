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
ROTATION_TYPE="size"  # Default to size for quick testing

# API process ID
API_PID=""

# Function to show usage
show_usage() {
    echo -e "${BLUE}Usage: $0 [rotation_type]${NC}"
    echo -e "${YELLOW}  rotation_type: 'size' (default) or 'daily'${NC}"
    echo -e "${YELLOW}Examples:${NC}"
    echo -e "  $0           # Test size-based rotation (quick)"
    echo -e "  $0 size      # Test size-based rotation (quick)"
    echo -e "  $0 daily     # Test daily rotation (requires date change simulation)"
    exit 1
}

# Parse command line arguments
if [ $# -gt 1 ]; then
    show_usage
elif [ $# -eq 1 ]; then
    case "$1" in
        "size"|"daily")
            ROTATION_TYPE="$1"
            ;;
        "-h"|"--help"|"help")
            show_usage
            ;;
        *)
            echo -e "${RED}❌ Invalid rotation type: $1${NC}"
            show_usage
            ;;
    esac
fi

# Adjust test parameters based on rotation type
if [ "$ROTATION_TYPE" = "daily" ]; then
    echo -e "${YELLOW}⚠️ Daily rotation testing requires date change simulation.${NC}"
    echo -e "${YELLOW}This test will demonstrate daily rotation setup but may not trigger actual rotation.${NC}"
    echo -e "${YELLOW}For full daily rotation testing, run the application across midnight.${NC}"
    TEST_DURATION=30  # Shorter for daily test
fi

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
    echo -e "${BLUE}Testing ${ROTATION_TYPE^^} rotation${NC}"
    
    # Backup existing .env file if it exists
    if [ -f .env ]; then
        cp .env .env.backup
        echo -e "${YELLOW}Backed up existing .env file to .env.backup${NC}"
    fi

    # Create test .env file with rotation type
    cat > .env.test << EOF
# Test environment configuration for log rotation
APP_ENV=development
PORT=$PORT
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT_PATH=stdout
LOG_ROTATION_TYPE=$ROTATION_TYPE
LOG_FILE_PATH=$LOG_FILE
LOG_FILE_MAX_SIZE=$MAX_SIZE_MB
LOG_FILE_MAX_BACKUPS=$BACKUP_COUNT
LOG_FILE_MAX_AGE=$MAX_AGE_DAYS
LOG_FILE_COMPRESS=true
EOF

    # Use the test .env file
    cp .env.test .env
    echo -e "${GREEN}Created test .env file with ${ROTATION_TYPE} rotation settings${NC}"
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
    echo -e "${YELLOW}  Rotation type: $ROTATION_TYPE${NC}"
    echo -e "${YELLOW}  Log file: $LOG_FILE${NC}"
    if [ "$ROTATION_TYPE" = "size" ]; then
        echo -e "${YELLOW}  Max size: $MAX_SIZE_MB MB${NC}"
    else
        echo -e "${YELLOW}  Daily rotation: new file created each day${NC}"
    fi
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
        if [ -f "$LOG_FILE" ]; then
            tail -n 10 $LOG_FILE
        else
            echo -e "${RED}No log file found at $LOG_FILE${NC}"
        fi
        cleanup
    fi

    echo -e "${GREEN}✅ API started successfully (PID: $API_PID)${NC}"

    # Test rotation based on type
    if [ "$ROTATION_TYPE" = "size" ]; then
        test_size_rotation
    else
        test_daily_rotation
    fi
}

# Test size-based rotation
test_size_rotation() {
    echo -e "${BLUE}Generating log entries to trigger size-based rotation...${NC}"
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
            # Check if rotation has occurred by looking for backup files
            backup_count=$(find "$LOG_DIR" -name "app.log.*" 2>/dev/null | wc -l)
            if [ $backup_count -gt 0 ]; then
                echo -e "${GREEN}✅ Size-based rotation occurred! Files in $LOG_DIR:${NC}"
                ls -lh "$LOG_DIR"
                break
            fi
        fi
    done
}

# Test daily rotation
test_daily_rotation() {
    echo -e "${BLUE}Testing daily rotation setup...${NC}"
    echo -e "${YELLOW}Daily rotation creates files with format: app-YYYY-MM-DD.log${NC}"
    
    # Wait a bit and check the log file pattern
    sleep 5
    
    # Look for daily log files
    today=$(date +%Y-%m-%d)
    expected_file="$LOG_DIR/app-$today.log"
    
    if [ -f "$expected_file" ]; then
        echo -e "${GREEN}✅ Daily log file created: $expected_file${NC}"
    else
        echo -e "${YELLOW}⚠️ Expected daily log file not found: $expected_file${NC}"
        echo -e "${YELLOW}Available log files:${NC}"
        ls -la "$LOG_DIR"
    fi
    
    # Generate some log entries
    echo -e "${BLUE}Generating log entries for daily log...${NC}"
    end_time=$((SECONDS + TEST_DURATION))
    
    while [ $SECONDS -lt $end_time ]; do
        curl -s "http://localhost:$PORT/health" > /dev/null
        curl -s "http://localhost:$PORT/api/v1/animals" > /dev/null
        sleep $LOG_INTERVAL
        
        remaining=$((end_time - SECONDS))
        if [ $((remaining % 10)) -eq 0 ]; then
            echo -e "${YELLOW}Remaining time: $remaining seconds${NC}"
            if [ -f "$expected_file" ]; then
                size=$(du -h "$expected_file" | cut -f1)
                echo -e "${GREEN}Current daily log file size: $size${NC}"
            fi
        fi
    done
    
    echo -e "${YELLOW}Note: Daily rotation occurs at midnight. To test actual rotation:${NC}"
    echo -e "${YELLOW}  1. Run the application before midnight${NC}"
    echo -e "${YELLOW}  2. Let it run past midnight${NC}"
    echo -e "${YELLOW}  3. Check for new daily log file creation${NC}"
}

# Stop the API gracefully
stop_api() {
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
    
    # Stop the API
    stop_api
    
    # Cleanup
    restore_env
    
    echo -e "${GREEN}Test completed!${NC}"
    echo -e "${YELLOW}Log files were retained in $LOG_DIR for inspection.${NC}"
    echo -e "${YELLOW}You can delete them manually when no longer needed.${NC}"
}

# Run the main function
main 