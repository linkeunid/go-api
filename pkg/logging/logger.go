package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkeunid/go-api/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogRotationType defines the type of log rotation
type LogRotationType string

const (
	// RotationTypeSize rotates logs based on file size
	RotationTypeSize LogRotationType = "size"
	// RotationTypeDaily rotates logs daily
	RotationTypeDaily LogRotationType = "daily"
)

// InitializeLogger creates and configures the logger based on configuration
func InitializeLogger(cfg *config.Config) (*zap.Logger, error) {
	// Determine rotation type from configuration
	rotationType := RotationTypeDaily // default
	if cfg.Logging.RotationType == "size" {
		rotationType = RotationTypeSize
	}
	return InitializeLoggerWithRotation(cfg, rotationType)
}

// InitializeLoggerWithRotation creates a logger with specified rotation type
func InitializeLoggerWithRotation(cfg *config.Config, rotationType LogRotationType) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	if cfg.IsDevelopment() {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Configure log level
	switch cfg.Logging.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	// Check if file output is enabled
	if cfg.Logging.FileOutputPath != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(cfg.Logging.FileOutputPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
		}

		// Configure encoder based on format
		var encoder zapcore.Encoder
		if cfg.Logging.Format == "json" {
			encoder = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
		}

		// Setup file output with rotation based on type
		var fileWriter zapcore.WriteSyncer
		switch rotationType {
		case RotationTypeDaily:
			dailyLogger := NewDailyRotateLogger(
				cfg.Logging.FileOutputPath,
				cfg.Logging.FileMaxSize,
				cfg.Logging.FileMaxBackups,
				cfg.Logging.FileMaxAge,
				cfg.Logging.FileCompress,
			)
			fileWriter = zapcore.AddSync(dailyLogger)
		case RotationTypeSize:
			fallthrough
		default:
			fileWriter = zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.Logging.FileOutputPath,
				MaxSize:    cfg.Logging.FileMaxSize,
				MaxBackups: cfg.Logging.FileMaxBackups,
				MaxAge:     cfg.Logging.FileMaxAge,
				Compress:   cfg.Logging.FileCompress,
			})
		}

		// Create a core that writes to the file
		fileCore := zapcore.NewCore(encoder, fileWriter, zapConfig.Level)

		// If standard output is also requested, create a multi-writer core
		if cfg.Logging.OutputPath == "stdout" || cfg.Logging.OutputPath == "stderr" {
			var stdWriter zapcore.WriteSyncer
			if cfg.Logging.OutputPath == "stderr" {
				stdWriter = zapcore.AddSync(os.Stderr)
			} else {
				stdWriter = zapcore.AddSync(os.Stdout)
			}

			// Create a core that writes to standard output
			stdCore := zapcore.NewCore(encoder, stdWriter, zapConfig.Level)

			// Use a tee to write to both outputs
			return zap.New(zapcore.NewTee(fileCore, stdCore), zap.AddCaller()), nil
		}

		// Only write to file
		return zap.New(fileCore, zap.AddCaller()), nil
	}

	// Standard zap logger with no file output
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// InitializeSizeBasedLogger creates a logger with size-based rotation (legacy)
func InitializeSizeBasedLogger(cfg *config.Config) (*zap.Logger, error) {
	return InitializeLoggerWithRotation(cfg, RotationTypeSize)
}

// InitializeDailyLogger creates a logger with daily rotation
func InitializeDailyLogger(cfg *config.Config) (*zap.Logger, error) {
	return InitializeLoggerWithRotation(cfg, RotationTypeDaily)
}
