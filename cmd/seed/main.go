package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/linkeunid/go-api/internal/bootstrap"
	"github.com/linkeunid/go-api/pkg/database"
	"github.com/linkeunid/go-api/pkg/seeder"
	"go.uber.org/zap"
)

// Command line flags
var (
	all        bool
	seederName string
	count      int
	help       bool
)

// Register command line flags
func init() {
	flag.BoolVar(&all, "all", false, "Run all seeders")
	flag.StringVar(&seederName, "seeder", "", "Run a specific seeder by name")
	flag.IntVar(&count, "count", 100, "Number of records to generate")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")
}

// Seeder interface for all seeders
type Seeder interface {
	Seed(ctx context.Context) error
	GetName() string
}

func main() {
	flag.Parse()

	// Show help
	if help {
		showHelp()
		os.Exit(0)
	}

	// Validate flags
	if !all && seederName == "" {
		fmt.Println("❌ Error: You must specify -all or -seeder=NAME")
		showHelp()
		os.Exit(1)
	}

	// Validate count
	if count <= 0 {
		fmt.Println("❌ Error: Count must be greater than 0")
		os.Exit(1)
	}

	// Initialize the app
	app, err := bootstrap.InitializeApp()
	if err != nil {
		fmt.Printf("❌ Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	// Get dependencies
	logger := app.Logger
	db := app.DB

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create seeders registry
	seeders := registerSeeders(db, logger, count)

	// Run seeders
	if all {
		runAllSeeders(ctx, seeders, logger)
	} else {
		runNamedSeeder(ctx, seeders, seederName, logger)
	}
}

// Register all available seeders
func registerSeeders(db database.Database, logger *zap.Logger, count int) []Seeder {
	return []Seeder{
		seeder.NewAnimalSeeder(db, logger, count),
		seeder.NewFlowerSeeder(db, logger, count),
		// Add more seeders here as they are implemented
	}
}

// Run all registered seeders
func runAllSeeders(ctx context.Context, seeders []Seeder, logger *zap.Logger) {
	logger.Info("Running all seeders", zap.Int("count", len(seeders)))

	for _, s := range seeders {
		runSeeder(ctx, s, logger)
	}

	logger.Info("All seeders completed successfully")
}

// Run a specific seeder by name
func runNamedSeeder(ctx context.Context, seeders []Seeder, name string, logger *zap.Logger) {
	name = strings.ToLower(name)
	logger.Info("Looking for seeder", zap.String("name", name))

	for _, s := range seeders {
		if strings.ToLower(s.GetName()) == name {
			runSeeder(ctx, s, logger)
			return
		}
	}

	logger.Error("Seeder not found", zap.String("name", name))
	fmt.Printf("❌ Seeder '%s' not found\n", name)
	showAvailableSeeders(seeders)
	os.Exit(1)
}

// Run a single seeder
func runSeeder(ctx context.Context, s Seeder, logger *zap.Logger) {
	name := s.GetName()
	logger.Info("Running seeder", zap.String("name", name))

	fmt.Printf("🌱 Running seeder: %s\n", name)
	startTime := time.Now()

	if err := s.Seed(ctx); err != nil {
		logger.Error("Failed to run seeder",
			zap.String("name", name),
			zap.Error(err),
		)
		fmt.Printf("❌ Seeder '%s' failed: %v\n", name, err)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	logger.Info("Seeder completed successfully",
		zap.String("name", name),
		zap.Duration("duration", duration),
	)
	fmt.Printf("✅ Seeder '%s' completed in %v\n", name, duration)
}

// Show help message
func showHelp() {
	fmt.Println("🌱 Database Seeder Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/seed [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -all           Run all seeders")
	fmt.Println("  -seeder=NAME   Run a specific seeder by name")
	fmt.Println("  -count=N       Number of records to generate (default: 100)")
	fmt.Println("  -help, -h      Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run ./cmd/seed -all                # Run all seeders with default count")
	fmt.Println("  go run ./cmd/seed -all -count=500     # Run all seeders with 500 records each")
	fmt.Println("  go run ./cmd/seed -seeder=animal      # Run only the animal seeder")
	fmt.Println("  go run ./cmd/seed -seeder=animal -count=50  # Run animal seeder with 50 records")
}

// Show available seeders
func showAvailableSeeders(seeders []Seeder) {
	fmt.Println("\nAvailable seeders:")
	for _, s := range seeders {
		fmt.Printf("  - %s\n", s.GetName())
	}
}
