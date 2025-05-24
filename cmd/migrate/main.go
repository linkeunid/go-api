package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqldriver "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/linkeunid/go-api/internal/model"
	"github.com/linkeunid/go-api/pkg/config"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const migrationsPath = "migrations"

// ModelRegistry contains all models that can be used for migrations
var ModelRegistry = map[string]interface{}{
	"animal": &model.Animal{},
	"flower": &model.Flower{},
}

// MigrationGenerator handles the generation of migrations using GORM
type MigrationGenerator struct {
	db *gorm.DB
}

// NewMigrationGenerator creates a new migration generator instance
func NewMigrationGenerator() (*MigrationGenerator, error) {
	cfg := config.LoadConfig()

	db, err := gorm.Open(gormMysql.Open(cfg.Database.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		DryRun:                                   true, // Enable dry run mode for DDL generation
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &MigrationGenerator{db: db}, nil
}

// GenerateModelMigration generates SQL for creating and dropping a table based on a GORM model
func (mg *MigrationGenerator) GenerateModelMigration(model interface{}) (upSQL, downSQL string, err error) {
	// Parse the model to get table information
	stmt := &gorm.Statement{DB: mg.db}
	if err := stmt.Parse(model); err != nil {
		return "", "", fmt.Errorf("failed to parse model: %w", err)
	}

	tableName := stmt.Schema.Table

	// Create a temporary database connection for actual DDL operations
	cfg := config.LoadConfig()
	actualDB, err := gorm.Open(gormMysql.Open(cfg.Database.DSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to database for DDL: %w", err)
	}

	actualMigrator := actualDB.Migrator()

	// Create the table temporarily to get the CREATE TABLE statement
	if err := actualMigrator.CreateTable(model); err != nil {
		return "", "", fmt.Errorf("failed to create table for DDL generation: %w", err)
	}

	// Get the underlying SQL database connection
	sqlDB, err := actualDB.DB()
	if err != nil {
		return "", "", fmt.Errorf("failed to get SQL database connection: %w", err)
	}

	// Get the CREATE TABLE statement
	var showTableName, createTableSQL string
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	if err := sqlDB.QueryRow(query).Scan(&showTableName, &createTableSQL); err != nil {
		return "", "", fmt.Errorf("failed to get CREATE TABLE SQL: %w", err)
	}

	// Clean up the temporary table
	if err := actualMigrator.DropTable(model); err != nil {
		log.Printf("Warning: Failed to clean up temporary table %s: %v", tableName, err)
	}

	// Format the migration SQL
	upSQL = mg.formatUpSQL(createTableSQL)
	downSQL = mg.formatDownSQL(tableName)

	return upSQL, downSQL, nil
}

// formatUpSQL formats the CREATE TABLE statement for the up migration
func (mg *MigrationGenerator) formatUpSQL(createTableSQL string) string {
	// Replace CREATE TABLE with CREATE TABLE IF NOT EXISTS
	createTableSQL = strings.Replace(createTableSQL, "CREATE TABLE", "CREATE TABLE IF NOT EXISTS", 1)

	return fmt.Sprintf(`-- Migration Up
-- SQL in section 'Up' is executed when this migration is applied

%s;`, createTableSQL)
}

// formatDownSQL creates the DROP TABLE statement for the down migration
func (mg *MigrationGenerator) formatDownSQL(tableName string) string {
	return fmt.Sprintf(`-- Migration Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS `+"`%s`"+`;`, tableName)
}

// MigrationManager handles all migration operations
type MigrationManager struct {
	migrator  *migrate.Migrate
	generator *MigrationGenerator
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager() (*MigrationManager, error) {
	generator, err := NewMigrationGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create migration generator: %w", err)
	}

	migrator := getMigrator()

	return &MigrationManager{
		migrator:  migrator,
		generator: generator,
	}, nil
}

func main() {
	// Parse command-line flags
	var (
		createCmd  = flag.Bool("create", false, "Create a new migration")
		upCmd      = flag.Bool("up", false, "Run all migrations or up to a specific version")
		downCmd    = flag.Bool("down", false, "Roll back all migrations or to a specific version")
		versionCmd = flag.Bool("version", false, "Show the current migration version")
		forceTo    = flag.Int("force", -1, "Force migration to a specific version")
		steps      = flag.Int("steps", 0, "Number of migrations to apply (use with -up or -down)")
		dryRun     = flag.Bool("dry-run", false, "Show what would be done without actually running migrations")
		fromModel  = flag.String("from-model", "", "Create migration from a model (e.g., animal)")
		listModels = flag.Bool("list-models", false, "List available models for migrations")
	)
	flag.Parse()

	// Get migration name from remaining arguments
	var migrationName string
	if flag.NArg() > 0 {
		migrationName = strings.Join(flag.Args(), "_")
	}

	// Handle commands
	switch {
	case *listModels:
		listAvailableModels()
	case *createCmd:
		handleCreateCommand(*fromModel, migrationName)
	case *upCmd:
		handleMigrationCommand("up", *steps, *dryRun)
	case *downCmd:
		handleMigrationCommand("down", *steps, *dryRun)
	case *versionCmd:
		showVersion()
	case *forceTo >= 0:
		forceMigration(*forceTo)
	default:
		showHelp()
	}
}

// handleCreateCommand handles the create migration command
func handleCreateCommand(fromModel, migrationName string) {
	manager, err := NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to initialize migration manager: %v", err)
	}

	if fromModel != "" {
		if migrationName == "" {
			migrationName = fmt.Sprintf("create_%s_table", fromModel)
		}
		createModelMigration(manager, fromModel, migrationName)
	} else {
		if migrationName == "" {
			log.Fatal("Migration name is required for create command")
		}
		createEmptyMigration(migrationName)
	}
}

// handleMigrationCommand handles up/down migration commands
func handleMigrationCommand(direction string, steps int, dryRun bool) {
	manager, err := NewMigrationManager()
	if err != nil {
		log.Fatalf("Failed to initialize migration manager: %v", err)
	}

	runMigrations(manager.migrator, direction, steps, dryRun)
}

// listAvailableModels lists all available models for migrations
func listAvailableModels() {
	fmt.Println("Available models for migration:")
	for name := range ModelRegistry {
		fmt.Printf("  - %s\n", name)
	}
}

// createModelMigration creates a migration based on a GORM model
func createModelMigration(manager *MigrationManager, modelName, migrationName string) {
	model, exists := ModelRegistry[strings.ToLower(modelName)]
	if !exists {
		fmt.Printf("Error: Model '%s' not found. Available models:\n", modelName)
		listAvailableModels()
		os.Exit(1)
	}

	// Generate migration files
	timestamp := time.Now().Unix()
	version := strconv.FormatInt(timestamp, 10)
	safeVersion := fmt.Sprintf("%s_%s", version, migrationName)

	// Ensure migrations directory exists
	if err := os.MkdirAll(migrationsPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Generate SQL using GORM migrator
	upSQL, downSQL, err := manager.generator.GenerateModelMigration(model)
	if err != nil {
		log.Fatalf("Failed to generate migration SQL: %v", err)
	}

	// Create migration files
	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.up.sql", safeVersion))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.down.sql", safeVersion))

	if err := os.WriteFile(upFile, []byte(upSQL), 0644); err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	if err := os.WriteFile(downFile, []byte(downSQL), 0644); err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created model-based migration files:\n  %s\n  %s\n", upFile, downFile)
}

// createEmptyMigration creates empty migration files
func createEmptyMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.up.sql", filename))
	downFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.down.sql", filename))

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	upTemplate := `-- Migration Up
-- SQL in section 'Up' is executed when this migration is applied

`

	downTemplate := `-- Migration Down
-- SQL in section 'Down' is executed when this migration is rolled back

`

	if err := os.WriteFile(upFile, []byte(upTemplate), 0644); err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	if err := os.WriteFile(downFile, []byte(downTemplate), 0644); err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
}

// runMigrations executes migration operations
func runMigrations(m *migrate.Migrate, direction string, steps int, dryRun bool) {
	if dryRun {
		showDryRunInfo(m, direction, steps)
		return
	}

	var err error
	switch {
	case steps > 0:
		fmt.Printf("Running %d %s migrations...\n", steps, direction)
		if direction == "up" {
			err = m.Steps(steps)
		} else {
			err = m.Steps(-steps)
		}
	case direction == "up":
		fmt.Println("Running all pending migrations...")
		err = m.Up()
	default:
		fmt.Println("Rolling back one migration...")
		err = m.Steps(-1)
	}

	handleMigrationResult(err)
}

// showDryRunInfo shows what migrations would be executed
func showDryRunInfo(m *migrate.Migrate, direction string, steps int) {
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get migration version: %v", err)
	}

	fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	switch {
	case steps > 0:
		fmt.Printf("Would %s %d migrations\n", getActionWord(direction), steps)
	case direction == "up":
		fmt.Println("Would apply all pending migrations")
	default:
		fmt.Println("Would roll back one migration")
	}
}

// getActionWord returns the appropriate action word for the direction
func getActionWord(direction string) string {
	if direction == "up" {
		return "apply"
	}
	return "roll back"
}

// handleMigrationResult handles the result of migration operations
func handleMigrationResult(err error) {
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
	if err := m.Force(version); err != nil {
		log.Fatalf("Failed to force migration: %v", err)
	}
	fmt.Printf("Successfully forced migration to version %d\n", version)
}

// getMigrator creates and returns a migrator instance
func getMigrator() *migrate.Migrate {
	cfg := config.LoadConfig()
	dsn := prepareDSNForMigration(cfg.Database.DSN)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	driver, err := mysqldriver.WithInstance(db, &mysqldriver.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    extractDatabaseName(dsn),
	})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "mysql", driver)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	return m
}

// prepareDSNForMigration prepares the DSN for migration use
func prepareDSNForMigration(dsn string) string {
	if strings.Contains(dsn, "?") {
		baseDSN := strings.Split(dsn, "?")[0]
		return fmt.Sprintf("%s?multiStatements=true", baseDSN)
	}
	return dsn
}

// extractDatabaseName extracts the database name from a DSN
func extractDatabaseName(dsn string) string {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		// Fallback for simple DSNs
		parts := strings.Split(dsn, "/")
		if len(parts) > 1 {
			dbPart := parts[len(parts)-1]
			return strings.Split(dbPart, "?")[0]
		}
		return ""
	}
	return cfg.DBName
}

// showHelp displays the help information
func showHelp() {
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
	fmt.Println("  migrate -dry-run -up|-down          Show migrations that would be applied")
	fmt.Println("  migrate -list-models                List available models for migrations")
	fmt.Println("\nExamples:")
	fmt.Println("  migrate -create add_users_table")
	fmt.Println("  migrate -create -from-model animal")
	fmt.Println("  migrate -up")
	fmt.Println("  migrate -down -steps 1")
	fmt.Println("  migrate -force 0               (Reset all migrations)")
}
