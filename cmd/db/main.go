package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/linkeunid/go-api/internal/bootstrap"
	"github.com/linkeunid/go-api/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Command line flags
var (
	truncateModel string
	truncateAll   bool
	help          bool
	verbose       bool
)

// Register command line flags
func init() {
	flag.StringVar(&truncateModel, "truncate", "", "Truncate a specific table based on model name")
	flag.BoolVar(&truncateAll, "truncate-all", false, "Truncate all tables")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
}

// Model map to get the table name for a model
var modelMap = map[string]interface{}{
	"animal": &model.Animal{},
	"flower": &model.Flower{},
	// Add more models here as they are implemented
}

func main() {
	flag.Parse()

	// Show help
	if help {
		showHelp()
		os.Exit(0)
	}

	// Validate flags
	if truncateModel == "" && !truncateAll {
		fmt.Println("❌ Error: You must specify -truncate MODEL or -truncate-all")
		showHelp()
		os.Exit(1)
	}

	// Initialize the app
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Logger.Sync()

	// Get dependencies
	logger := app.Logger
	db := app.DB.GetDB()

	// Process command
	if truncateModel != "" {
		truncateModelTable(logger, db, truncateModel)
	} else if truncateAll {
		truncateAllTables(logger, db)
	}
}

// truncateModelTable truncates a single table based on the model name
func truncateModelTable(logger *zap.Logger, db *gorm.DB, modelName string) {
	modelName = strings.ToLower(modelName)

	// Check if model exists
	model, exists := modelMap[modelName]
	if !exists {
		logger.Error("Model not found", zap.String("model", modelName))
		fmt.Printf("❌ Model '%s' not found\n", modelName)
		showAvailableModels()
		os.Exit(1)
	}

	// Get the table name using GORM's TableName method if available
	tableName := ""

	// Check if the model implements TableName() method
	if tableNamer, ok := model.(interface{ TableName() string }); ok {
		tableName = tableNamer.TableName()
	} else {
		// Default table name (lowercase model name with 's' appended)
		tableName = modelName + "s"
	}

	// Execute TRUNCATE statement
	logger.Info("Truncating table", zap.String("table", tableName))
	result := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))

	// Check for errors
	if result.Error != nil {
		logger.Error("Failed to truncate table",
			zap.String("table", tableName),
			zap.Error(result.Error),
		)
		fmt.Printf("❌ Failed to truncate table '%s': %v\n",
			tableName, result.Error)
		os.Exit(1)
	}

	logger.Info("Table truncated successfully", zap.String("table", tableName))
}

// truncateAllTables truncates all tables in the database
func truncateAllTables(logger *zap.Logger, db *gorm.DB) {
	logger.Info("Truncating all tables")

	// Execute SET FOREIGN_KEY_CHECKS=0 to temporarily disable foreign key constraints
	db.Exec("SET FOREIGN_KEY_CHECKS=0")

	// Truncate each table
	for modelName, model := range modelMap {
		// Get table name
		tableName := ""
		if tableNamer, ok := model.(interface{ TableName() string }); ok {
			tableName = tableNamer.TableName()
		} else {
			tableName = modelName + "s"
		}

		// Execute TRUNCATE statement
		if verbose {
			logger.Info("Truncating table", zap.String("table", tableName))
		}

		result := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))

		// Check for errors but continue with other tables
		if result.Error != nil {
			logger.Error("Failed to truncate table",
				zap.String("table", tableName),
				zap.Error(result.Error),
			)
			if verbose {
				fmt.Printf("❌ Failed to truncate table '%s': %v\n",
					tableName, result.Error)
			}
		} else if verbose {
			fmt.Printf("✅ Truncated table '%s'\n", tableName)
		}
	}

	// Execute SET FOREIGN_KEY_CHECKS=1 to re-enable foreign key constraints
	db.Exec("SET FOREIGN_KEY_CHECKS=1")

	logger.Info("All tables truncated successfully")
}

// showHelp displays help information
func showHelp() {
	fmt.Println("🗄️ Database Operations Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/db [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -truncate MODEL  Truncate a specific table based on model name")
	fmt.Println("  -truncate-all    Truncate all tables")
	fmt.Println("  -v               Verbose output")
	fmt.Println("  -help, -h        Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/db -truncate animal    # Truncate only the animals table")
	fmt.Println("  go run ./cmd/db -truncate-all       # Truncate all tables")
	fmt.Println("")
	showAvailableModels()
	fmt.Println("")
	fmt.Println("Note: To update the model list after adding new models, run:")
	fmt.Println("  make update-model-map")
	fmt.Println("  (or make um for short)")
}

// showAvailableModels displays a list of available models
func showAvailableModels() {
	fmt.Println("Available models:")
	for modelName := range modelMap {
		fmt.Printf("  - %s\n", modelName)
	}
}
