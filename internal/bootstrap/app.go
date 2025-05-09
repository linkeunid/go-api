package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/linkeunid/go-api/internal/controller"
	"github.com/linkeunid/go-api/internal/repository"
	"github.com/linkeunid/go-api/internal/service"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// App represents the application dependencies
type App struct {
	Logger           *zap.Logger
	DB               database.Database
	Config           *config.Config
	AnimalController *controller.Animal
}

// InitializeApp initializes the application dependencies
func InitializeApp() (*App, error) {
	// Log environment info in development mode
	if env := os.Getenv("APP_ENV"); env == "" || env == "development" {
		dbPort := os.Getenv("DB_PORT")
		dbHost := os.Getenv("DB_HOST")
		fmt.Printf("Environment variables - DB_HOST: %s, DB_PORT: %s\n", dbHost, dbPort)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger, err := initializeLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
	dbWrapper, err := initializeDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories and services
	animalRepo := repository.NewAnimalRepository(dbWrapper, logger)
	animalService := service.NewAnimalService(cfg, logger, animalRepo)

	// Initialize controllers
	animalController := controller.NewAnimal(logger, animalService)

	// Configure Swagger
	SetupSwagger(cfg.Server.Port, cfg.IsDevelopment())

	// Return the app with all dependencies
	return &App{
		Logger:           logger,
		DB:               dbWrapper,
		Config:           cfg,
		AnimalController: animalController,
	}, nil
}

// initializeLogger creates and configures the logger
func initializeLogger(cfg *config.Config) (*zap.Logger, error) {
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

		// Setup file output with log rotation
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Logging.FileOutputPath,
			MaxSize:    cfg.Logging.FileMaxSize,    // megabytes
			MaxBackups: cfg.Logging.FileMaxBackups, // number of backups
			MaxAge:     cfg.Logging.FileMaxAge,     // days
			Compress:   cfg.Logging.FileCompress,   // compress rotated files
		})

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

// initializeDatabase sets up the database connection
func initializeDatabase(cfg *config.Config, logger *zap.Logger) (database.Database, error) {
	// Configure GORM logger
	logLevel := gormlogger.Info
	if cfg.IsProduction() {
		logLevel = gormlogger.Error
	}

	gormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  !cfg.IsProduction(),
		},
	)

	// Log DSN with password masked for debugging
	dsnForLog := GetDataSourceInfo(cfg.Database.DSN)
	logger.Info("Connecting to database", zap.String("dsn", dsnForLog))

	// Connect to database
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Successfully connected to database")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Create cache manager (nil for now)
	var cacheManager database.CacheManager

	// Initialize Redis cache manager if Redis is enabled
	if cfg.Redis.Enabled {
		logger.Info("Redis caching is enabled",
			zap.String("host", cfg.Redis.Host),
			zap.Int("port", cfg.Redis.Port),
			zap.Duration("cacheTTL", cfg.Redis.CacheTTL))

		// Initialize Redis cache manager
		redisManager, err := database.NewRedisCacheManager(cfg, logger)
		if err != nil {
			logger.Warn("Failed to initialize Redis cache manager, continuing without caching", zap.Error(err))
		} else {
			cacheManager = redisManager
			logger.Info("Redis cache manager initialized successfully")
		}
	} else {
		logger.Info("Redis caching is disabled")
	}

	// Create database wrapper
	dbWrapper := database.NewDatabase(cfg, logger, db, cacheManager)

	return dbWrapper, nil
}
