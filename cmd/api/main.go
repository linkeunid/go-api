package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/linkeunid/go-api/internal/controller"
	swaggerdocs "github.com/linkeunid/go-api/internal/docs/swaggerdocs" // Import swagger docs with named import
	"github.com/linkeunid/go-api/internal/repository"
	"github.com/linkeunid/go-api/internal/service"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/database"
	custommiddleware "github.com/linkeunid/go-api/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// maskPassword replaces password in DSN with asterisks
func maskPassword(dsn string) string {
	// Match password pattern in MySQL DSN: user:password@tcp(...)
	re := regexp.MustCompile(`([^:]+):([^@]+)@`)
	return re.ReplaceAllString(dsn, "$1:******@")
}

func main() {
	// Debug environment variables only in development mode
	if env := os.Getenv("APP_ENV"); env == "" || env == "development" {
		dbPort := os.Getenv("DB_PORT")
		dbHost := os.Getenv("DB_HOST")
		fmt.Printf("Environment variables - DB_HOST: %s, DB_PORT: %s\n", dbHost, dbPort)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Set Swagger Info only in development mode
	if cfg.IsDevelopment() {
		swaggerdocs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Server.Port)
		swaggerdocs.SwaggerInfo.Title = "Linkeun Go API"
		swaggerdocs.SwaggerInfo.Description = "API for managing various resources including animals"
		swaggerdocs.SwaggerInfo.Version = "1.0"
		swaggerdocs.SwaggerInfo.BasePath = "/api/v1"
		swaggerdocs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

	// Initialize logger
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

	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database
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
	dsnForLog := maskPassword(cfg.Database.DSN)
	logger.Info("Connecting to database", zap.String("dsn", dsnForLog))

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err), zap.String("dsn", dsnForLog))
	}

	logger.Info("Successfully connected to database")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database connection", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Create database wrapper
	dbWrapper := database.NewDatabase(cfg, logger, db, nil)

	// Initialize repositories
	animalRepo := repository.NewAnimalRepository(dbWrapper, logger)

	// Initialize services
	animalService := service.NewAnimalService(cfg, logger, animalRepo)

	// Initialize controllers
	animalController := controller.NewAnimal(logger, animalService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(custommiddleware.ValidationMiddleware) // Add our custom validation middleware

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Register routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger documentation - only available in development mode
	if cfg.IsDevelopment() {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The URL points to API definition
		))
		logger.Info("Swagger UI enabled in development mode")
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Animal routes
		animalController.RegisterRoutes(r)
	})

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("Starting server", zap.Int("port", cfg.Server.Port))
		if cfg.IsDevelopment() {
			logger.Info("Swagger UI available at", zap.String("url", fmt.Sprintf("http://localhost:%d/swagger/", cfg.Server.Port)))
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
