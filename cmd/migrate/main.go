package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqldriver "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/pkg/config"
)

const (
	migrationsPath = "migrations"
)

// modelRegistry contains all models that can be used for migrations
var modelRegistry = map[string]interface{}{
	"animal": model.Animal{},
}

func main() {
	// Initialize environment
	loadEnv()

	// Parse command-line arguments
	createCmd := flag.Bool("create", false, "Create a new migration")
	upCmd := flag.Bool("up", false, "Run all migrations or up to a specific version")
	downCmd := flag.Bool("down", false, "Roll back all migrations or to a specific version")
	versionCmd := flag.Bool("version", false, "Show the current migration version")
	forceTo := flag.Int("force", -1, "Force migration to a specific version")
	steps := flag.Int("steps", 0, "Number of migrations to apply (use with -up or -down)")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without actually running migrations")
	fromModel := flag.String("from-model", "", "Create migration from a model (e.g., animal)")
	listModels := flag.Bool("list-models", false, "List available models for migrations")
	flag.Parse()

	// Get migration name from remaining arguments
	var migrationName string
	if flag.NArg() > 0 {
		migrationName = strings.Join(flag.Args(), "_")
	}

	// List available models
	if *listModels {
		listAvailableModels()
		return
	}

	// Execute the appropriate command
	if *createCmd {
		if *fromModel != "" {
			// Create migration from model
			if migrationName == "" {
				migrationName = fmt.Sprintf("create_%s_table", *fromModel)
			}
			createModelMigration(*fromModel, migrationName)
		} else {
			// Create empty migration
			if migrationName == "" {
				log.Fatal("Migration name is required for create command")
			}
			createEmptyMigration(migrationName)
		}
	} else if *upCmd {
		runMigrations("up", *steps, *dryRun)
	} else if *downCmd {
		runMigrations("down", *steps, *dryRun)
	} else if *versionCmd {
		showVersion()
	} else if *forceTo >= 0 {
		forceMigration(*forceTo)
	} else {
		// Display help if no command is provided
		fmt.Println("Migration tool for LinkeunID Go API")
		fmt.Println("\nUsage:")
		fmt.Println("  migrate -create NAME                Create a new empty migration")
		fmt.Println("  migrate -create -from-model MODEL   Create a migration from a model")
		fmt.Println("  migrate -up                         Run all pending migrations")
		fmt.Println("  migrate -up -steps N                Run N up migrations")
		fmt.Println("  migrate -down                       Roll back the last migration")
		fmt.Println("  migrate -down -steps N              Roll back N migrations")
		fmt.Println("  migrate -version                    Show the current migration version")
		fmt.Println("  migrate -force VERSION              Force migration to a specific version")
		fmt.Println("  migrate -dry-run -up|-down          Show migrations that would be applied without running them")
		fmt.Println("  migrate -list-models                List available models for migrations")
		fmt.Println("\nExamples:")
		fmt.Println("  migrate -create add_users_table")
		fmt.Println("  migrate -create -from-model animal")
		fmt.Println("  migrate -up")
		fmt.Println("  migrate -down -steps 1")
		fmt.Println("  migrate -force 0               (Reset all migrations)")
	}
}

// listAvailableModels lists all available models that can be used for migrations
func listAvailableModels() {
	fmt.Println("Available models for migration:")
	for name := range modelRegistry {
		fmt.Printf("  - %s\n", name)
	}
}

// createModelMigration creates a migration based on a model
func createModelMigration(modelName string, migrationName string) {
	model, exists := modelRegistry[strings.ToLower(modelName)]
	if !exists {
		fmt.Printf("Error: Model '%s' not found. Available models:\n", modelName)
		listAvailableModels()
		os.Exit(1)
	}

	// Create timestamp and filenames
	timestamp := time.Now().Unix()
	version := strconv.FormatInt(timestamp, 10)
	safeVersion := fmt.Sprintf("%s_%s", version, migrationName)

	// Ensure migrations directory exists
	err := os.MkdirAll(migrationsPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create migration files
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.up.sql", safeVersion))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.down.sql", safeVersion))

	// Generate SQL from model
	upSQL, downSQL := generateSQLFromModel(model, modelName)

	// Write up migration
	err = os.WriteFile(upFile, []byte(upSQL), 0644)
	if err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	// Write down migration
	err = os.WriteFile(downFile, []byte(downSQL), 0644)
	if err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created model-based migration files:\n  %s\n  %s\n", upFile, downFile)
}

// generateSQLFromModel converts a Go struct to SQL CREATE TABLE and DROP TABLE statements
func generateSQLFromModel(model interface{}, modelName string) (string, string) {
	modelType := reflect.TypeOf(model)
	tableName := getTableName(model)

	// Start building up SQL
	upSQL := fmt.Sprintf("-- Migration Up\n-- SQL in section 'Up' is executed when this migration is applied\n\nCREATE TABLE IF NOT EXISTS `%s` (\n", tableName)

	// Track fields and indexes
	var fields []string
	var indexes []string
	var primaryKey string

	// Analyze struct fields
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get the field name and tags
		fieldName := field.Name
		gormTag := field.Tag.Get("gorm")
		jsonTag := field.Tag.Get("json")

		// Extract column name from json tag or use field name
		columnName := fieldName
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				columnName = parts[0]
			}
		}
		columnName = toSnakeCase(columnName)

		// Parse GORM tags
		columnType := getSQLType(field.Type)
		columnDef := fmt.Sprintf("  `%s` %s", columnName, columnType)

		// Handle common GORM attributes
		if strings.Contains(gormTag, "primaryKey") {
			primaryKey = columnName
			if strings.Contains(gormTag, "autoIncrement") {
				columnDef += " NOT NULL AUTO_INCREMENT"
			}
		}

		if strings.Contains(gormTag, "not null") {
			columnDef += " NOT NULL"
		}

		if strings.Contains(gormTag, "default:") {
			re := regexp.MustCompile(`default:([^;]+)`)
			matches := re.FindStringSubmatch(gormTag)
			if len(matches) > 1 {
				columnDef += fmt.Sprintf(" DEFAULT %s", matches[1])
			}
		}

		// Handle auto timestamps
		if strings.Contains(gormTag, "autoCreateTime") {
			columnDef += " DEFAULT CURRENT_TIMESTAMP(3)"
		}
		if strings.Contains(gormTag, "autoUpdateTime") {
			columnDef += " DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"
		}

		fields = append(fields, columnDef)

		// Handle indexes
		if strings.Contains(gormTag, "index:") {
			re := regexp.MustCompile(`index:([^;,]+)`)
			matches := re.FindStringSubmatch(gormTag)
			if len(matches) > 1 {
				indexName := matches[1]
				indexes = append(indexes, fmt.Sprintf("  KEY `%s` (`%s`)", indexName, columnName))
			} else {
				// Default index name
				indexes = append(indexes, fmt.Sprintf("  KEY `idx_%s_%s` (`%s`)", tableName, columnName, columnName))
			}
		}
	}

	// Add primary key
	if primaryKey != "" {
		fields = append(fields, fmt.Sprintf("  PRIMARY KEY (`%s`)", primaryKey))
	}

	// Add all fields and indexes
	allLines := append(fields, indexes...)
	upSQL += strings.Join(allLines, ",\n")
	upSQL += "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	// Create down SQL (drop table)
	downSQL := fmt.Sprintf("-- Migration Down\n-- SQL in section 'Down' is executed when this migration is rolled back\n\nDROP TABLE IF EXISTS `%s`;", tableName)

	return upSQL, downSQL
}

// getTableName gets the table name from a model, calling its TableName() method if available
func getTableName(model interface{}) string {
	// Try to call TableName() method
	if tableNamer, ok := model.(interface{ TableName() string }); ok {
		return tableNamer.TableName()
	}

	// Fallback to type name
	t := reflect.TypeOf(model)
	return toSnakeCase(t.Name())
}

// toSnakeCase converts a string from CamelCase to snake_case
func toSnakeCase(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(unicode.ToLower(r))
	}
	return result
}

// getSQLType converts a Go type to an SQL type
func getSQLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "bigint"
	case reflect.Int8, reflect.Int16:
		return "int"
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return "bigint unsigned"
	case reflect.Uint8, reflect.Uint16:
		return "int unsigned"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "varchar(255)"
	case reflect.Struct:
		if t.Name() == "Time" {
			return "datetime(3)"
		}
	}
	return "text"
}

// createEmptyMigration creates a new empty migration file
func createEmptyMigration(name string) {
	timestamp := time.Now().Unix()
	version := strconv.FormatInt(timestamp, 10)
	safeVersion := fmt.Sprintf("%s_%s", version, name)

	// Ensure migrations directory exists
	err := os.MkdirAll(migrationsPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create migration files
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.up.sql", safeVersion))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.down.sql", safeVersion))

	// Create up migration
	err = os.WriteFile(upFile, []byte("-- Migration Up\n-- SQL in section 'Up' is executed when this migration is applied\n\n"), 0644)
	if err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	// Create down migration
	err = os.WriteFile(downFile, []byte("-- Migration Down\n-- SQL in section 'Down' is executed when this migration is rolled back\n\n"), 0644)
	if err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" || env == "development" {
		err := godotenv.Load()
		if err != nil {
			// It's okay if .env doesn't exist
			fmt.Println("Warning: .env file not found, using environment variables")
		} else {
			fmt.Println("Successfully loaded .env file")
		}
	}
}

// runMigrations runs or rolls back migrations
func runMigrations(direction string, steps int, dryRun bool) {
	m := getMigrator()

	if dryRun {
		// Get current version
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get migration version: %v", err)
		}

		// Show what would be done
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		if steps > 0 {
			if direction == "up" {
				fmt.Printf("Would apply %d migrations up\n", steps)
			} else {
				fmt.Printf("Would roll back %d migrations\n", steps)
			}
		} else {
			if direction == "up" {
				fmt.Println("Would apply all pending migrations")
			} else {
				fmt.Println("Would roll back all migrations")
			}
		}
		return
	}

	// Execute migrations
	var err error
	if steps > 0 {
		fmt.Printf("Running %d %s migrations...\n", steps, direction)
		if direction == "up" {
			err = m.Steps(steps)
		} else {
			err = m.Steps(-steps)
		}
	} else {
		if direction == "up" {
			fmt.Println("Running all pending migrations...")
			err = m.Up()
		} else {
			fmt.Println("Rolling back one migration...")
			err = m.Steps(-1) // By default, roll back one step
		}
	}

	// Handle migration errors
	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migration changes to apply")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("Migration completed successfully!")
	}
}

// showVersion displays the current migration version
func showVersion() {
	m := getMigrator()
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations have been applied yet")
		} else {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		return
	}
	fmt.Printf("Current migration version: %d\n", version)
	if dirty {
		fmt.Println("Warning: Database is in a dirty state, last migration failed")
	}
}

// forceMigration forces a migration to a specific version
func forceMigration(version int) {
	m := getMigrator()
	err := m.Force(version)
	if err != nil {
		log.Fatalf("Failed to force migration: %v", err)
	}
	fmt.Printf("Successfully forced migration to version %d\n", version)
}

// getMigrator returns a new migrator instance
func getMigrator() *migrate.Migrate {
	// Get database configuration
	cfg := config.LoadConfig()

	// Extract the DSN without params for migrate (it adds its own params)
	dsn := cfg.Database.DSN

	// If using a complex DSN, parse it to ensure compatibility
	// Some params might interfere with migrations
	if strings.Contains(dsn, "?") {
		baseDSN := strings.Split(dsn, "?")[0]
		// Add only the essential parameters
		dsn = fmt.Sprintf("%s?multiStatements=true", baseDSN)
	}

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Create migration driver instance
	driver, err := mysqldriver.WithInstance(db, &mysqldriver.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    extractDatabaseName(dsn),
	})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Create migrator
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "mysql", driver)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	return m
}

// extractDatabaseName extracts the database name from a DSN
func extractDatabaseName(dsn string) string {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		// If we can't parse, try a simple extraction (for basic DSNs)
		parts := strings.Split(dsn, "/")
		if len(parts) > 1 {
			dbPart := parts[len(parts)-1]
			return strings.Split(dbPart, "?")[0]
		}
		return ""
	}
	return cfg.DBName
}
