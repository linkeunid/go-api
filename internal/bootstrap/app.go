package bootstrap

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/linkeunid/go-api/internal/controller"
	"github.com/linkeunid/go-api/internal/repository"
	"github.com/linkeunid/go-api/internal/service"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/database"
	"github.com/linkeunid/go-api/pkg/logging"
	"go.uber.org/zap"
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
	logger, err := logging.InitializeLogger(cfg)
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
