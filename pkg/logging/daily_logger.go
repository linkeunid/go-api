package logging

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// DailyRotateLogger provides daily log rotation functionality
// It wraps lumberjack.Logger and creates a new log file each day
type DailyRotateLogger struct {
	baseFilename string
	currentDate  string
	logger       *lumberjack.Logger
	maxSize      int
	maxBackups   int
	maxAge       int
	compress     bool
	mu           sync.Mutex
}

// NewDailyRotateLogger creates a new daily rotating logger
func NewDailyRotateLogger(baseFilename string, maxSize, maxBackups, maxAge int, compress bool) *DailyRotateLogger {
	dl := &DailyRotateLogger{
		baseFilename: baseFilename,
		currentDate:  time.Now().Format("2006-01-02"),
		maxSize:      maxSize,
		maxBackups:   maxBackups,
		maxAge:       maxAge,
		compress:     compress,
	}

	// Initialize the logger with today's filename
	dl.logger = &lumberjack.Logger{
		Filename:   dl.generateDailyFilename(),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return dl
}

// generateDailyFilename creates a filename with the current date
func (d *DailyRotateLogger) generateDailyFilename() string {
	dir := filepath.Dir(d.baseFilename)
	ext := filepath.Ext(d.baseFilename)
	name := d.baseFilename[:len(d.baseFilename)-len(ext)]

	today := time.Now().Format("2006-01-02")
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", filepath.Base(name), today, ext))
}

// Write implements io.Writer interface
func (d *DailyRotateLogger) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	today := time.Now().Format("2006-01-02")

	// Check if date has changed
	if today != d.currentDate {
		// Close current logger
		if d.logger != nil {
			d.logger.Close()
		}

		// Update current date
		d.currentDate = today

		// Create new logger with today's filename
		d.logger = &lumberjack.Logger{
			Filename:   d.generateDailyFilename(),
			MaxSize:    d.maxSize,
			MaxBackups: d.maxBackups,
			MaxAge:     d.maxAge,
			Compress:   d.compress,
		}
	}

	return d.logger.Write(p)
}

// Close closes the underlying logger
func (d *DailyRotateLogger) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.logger != nil {
		return d.logger.Close()
	}
	return nil
}

// Rotate triggers a rotation of the current log file
func (d *DailyRotateLogger) Rotate() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.logger != nil {
		return d.logger.Rotate()
	}
	return nil
}
