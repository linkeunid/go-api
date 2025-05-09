package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Init loads .env file only in development mode
func init() {
	// Check if APP_ENV is set and not 'development'
	env := os.Getenv("APP_ENV")
	if env != "" && env != "development" {
		// Don't load .env in non-development environments
		return
	}

	// Only load .env file in development mode
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in development
		// Just log to stdout for debugging
		if os.Getenv("APP_ENV") == "development" {
			fmt.Println("Warning: .env file not found, using environment variables")
		}
	} else if os.Getenv("APP_ENV") == "development" {
		fmt.Println("Successfully loaded .env file")
	}
}

// Config represents application configuration
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Logging     LoggingConfig
	Auth        AuthConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Enabled      bool
	Host         string
	Port         int
	Password     string
	DB           int
	CacheTTL     time.Duration
	PaginatedTTL time.Duration
	QueryCache   bool
	KeyPrefix    string
	PoolSize     int
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level          string
	Format         string
	OutputPath     string
	FileOutputPath string // Path to log file if file logging is enabled
	FileMaxSize    int    // Maximum size of log files in megabytes before rotation
	FileMaxBackups int    // Maximum number of old log files to retain
	FileMaxAge     int    // Maximum number of days to retain old log files
	FileCompress   bool   // Whether to compress rotated log files
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled        bool          // Whether authentication is enabled
	JWTSecret      string        // Secret key for JWT signing
	JWTExpiration  time.Duration // JWT expiration time
	AllowedIssuers []string      // Allowed JWT issuers
}

// LoadConfig loads application configuration from environment variables
func LoadConfig() *Config {
	// Get current environment
	env := getEnv("APP_ENV", "development")

	// Prepare DSN if not explicitly provided
	dsn := getEnv("DSN", "")
	if dsn == "" {
		// Build DSN from individual DB_* environment variables
		dbUser := getEnv("DB_USER", "root")
		dbPassword := getEnv("DB_PASSWORD", "root")
		dbHost := getEnv("DB_HOST", "localhost")

		// Use environment-specific default port if not specified
		var defaultDBPort int
		if env == "production" {
			// In production, prefer 3306 as default port
			defaultDBPort = 3306
			// For production Docker environments, prefer mysql hostname
			if dbHost == "localhost" && os.Getenv("DOCKER_ENV") == "true" {
				dbHost = "mysql"
			}
		} else {
			// In development, prefer 3307 as default port
			defaultDBPort = 3307
		}

		dbPort := getEnvAsInt("DB_PORT", defaultDBPort)
		dbName := getEnv("DB_NAME", "linkeun_go_api")
		dbParams := getEnv("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local")

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			dbUser, dbPassword, dbHost, dbPort, dbName, dbParams)
	}

	return &Config{
		Environment: env,
		Server: ServerConfig{
			Port:            getEnvAsInt("PORT", 8080),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             dsn,
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Enabled:      getEnvAsBool("REDIS_ENABLED", false),
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getRedisPort(env),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			CacheTTL:     getEnvAsDuration("REDIS_CACHE_TTL", 15*time.Minute),
			PaginatedTTL: getEnvAsDuration("REDIS_PAGINATED_TTL", 5*time.Minute),
			QueryCache:   getEnvAsBool("REDIS_QUERY_CACHING", true),
			KeyPrefix:    getEnv("REDIS_KEY_PREFIX", "linkeun_api:"),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
		},
		Logging: LoggingConfig{
			Level:          getLogLevel(env),
			Format:         getEnv("LOG_FORMAT", "json"),
			OutputPath:     getEnv("LOG_OUTPUT_PATH", "stdout"),
			FileOutputPath: getEnv("LOG_FILE_PATH", ""),
			FileMaxSize:    getEnvAsInt("LOG_FILE_MAX_SIZE", 100),
			FileMaxBackups: getEnvAsInt("LOG_FILE_MAX_BACKUPS", 3),
			FileMaxAge:     getEnvAsInt("LOG_FILE_MAX_AGE", 28),
			FileCompress:   getEnvAsBool("LOG_FILE_COMPRESS", true),
		},
		Auth: AuthConfig{
			Enabled:        getEnvAsBool("AUTH_ENABLED", false),
			JWTSecret:      getEnv("JWT_SECRET", ""),
			JWTExpiration:  getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
			AllowedIssuers: getEnvAsSlice("JWT_ALLOWED_ISSUERS", []string{}, ","),
		},
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsTest returns true if the environment is test
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}

	// If the duration string doesn't have a unit, assume seconds
	if i, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(i) * time.Second
	}

	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value := strings.Split(valueStr, separator)

	// Remove empty strings
	var result []string
	for _, v := range value {
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

// getRedisPort returns the default Redis port based on environment
func getRedisPort(env string) int {
	defaultPort := 6379 // Default for production
	if env != "production" {
		defaultPort = 6380 // Default for development
	}
	return getEnvAsInt("REDIS_PORT", defaultPort)
}

// getLogLevel returns the default log level based on environment
func getLogLevel(env string) string {
	defaultLevel := "info" // Default for development
	if env == "production" {
		defaultLevel = "error" // Default for production
	}
	return getEnv("LOG_LEVEL", defaultLevel)
}
