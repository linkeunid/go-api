package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDailyRotateLogger(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "daily_logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test file path
	logFile := filepath.Join(tempDir, "test.log")

	// Create daily logger
	dailyLogger := NewDailyRotateLogger(logFile, 1, 3, 7, false)
	defer dailyLogger.Close()

	// Test initial write
	testMessage := "Test log message\n"
	n, err := dailyLogger.Write([]byte(testMessage))
	if err != nil {
		t.Fatalf("Failed to write to daily logger: %v", err)
	}
	if n != len(testMessage) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testMessage), n)
	}

	// Verify file was created with today's date
	today := time.Now().Format("2006-01-02")
	expectedFilename := filepath.Join(tempDir, "test-"+today+".log")

	if _, err := os.Stat(expectedFilename); os.IsNotExist(err) {
		t.Errorf("Expected daily log file %s was not created", expectedFilename)
	}

	// Verify content
	content, err := os.ReadFile(expectedFilename)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if string(content) != testMessage {
		t.Errorf("Expected log content %q, got %q", testMessage, string(content))
	}
}

func TestDailyRotateLoggerFilename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "daily_logger_filename_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test different base filenames
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    filepath.Join(tempDir, "app.log"),
			expected: "app-" + time.Now().Format("2006-01-02") + ".log",
		},
		{
			input:    filepath.Join(tempDir, "service"),
			expected: "service-" + time.Now().Format("2006-01-02"),
		},
		{
			input:    filepath.Join(tempDir, "api.txt"),
			expected: "api-" + time.Now().Format("2006-01-02") + ".txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			dailyLogger := NewDailyRotateLogger(tc.input, 1, 3, 7, false)
			defer dailyLogger.Close()

			// Write something to trigger file creation
			_, err := dailyLogger.Write([]byte("test"))
			if err != nil {
				t.Fatalf("Failed to write: %v", err)
			}

			expectedPath := filepath.Join(tempDir, tc.expected)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", expectedPath)
			}
		})
	}
}

func TestDailyRotateLoggerConcurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "daily_logger_concurrency_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "concurrent.log")
	dailyLogger := NewDailyRotateLogger(logFile, 1, 3, 7, false)
	defer dailyLogger.Close()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				message := fmt.Sprintf("Message from goroutine %d, iteration %d\n", id, j)
				_, err := dailyLogger.Write([]byte(message))
				if err != nil {
					t.Errorf("Failed to write from goroutine %d: %v", id, err)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify file exists and has content
	today := time.Now().Format("2006-01-02")
	expectedFilename := filepath.Join(tempDir, "concurrent-"+today+".log")

	info, err := os.Stat(expectedFilename)
	if os.IsNotExist(err) {
		t.Errorf("Expected log file %s was not created", expectedFilename)
	} else if info.Size() == 0 {
		t.Error("Log file is empty after concurrent writes")
	}
}
