# Seeder Mapper Tool

This tool automatically manages the seeder registry in `cmd/seed/main.go` by scanning the `pkg/seeder` directory for seeder implementations and updating the `registerSeeders` function accordingly.

## Overview

The seeder mapper is designed to automatically:
- Scan the `pkg/seeder` directory for files matching the pattern `*_seeder.go`
- Detect struct types that implement the `Seeder` interface
- Update the `registerSeeders` function in `cmd/seed/main.go` with the discovered seeders

## Usage

### Commands

#### Basic Update (Default)
```bash
go run ./cmd/seeder-mapper
# or
make update-seeder-registry
# or
make usr  # alias
```
Scans for all seeders and adds any new ones to the registry. Does not remove existing entries.

#### Sync Mode
```bash
go run ./cmd/seeder-mapper -sync
# or
make sync-seeder
# or
make ss  # alias
```
Both adds new seeders and removes seeders that no longer exist in the filesystem.

#### Clean Only Mode
```bash
go run ./cmd/seeder-mapper -clean-only
# or
make clean-seeder-registry
# or
make csr  # alias
```
Only removes seeders that no longer exist in the filesystem. Does not add new ones.

#### Verbose Mode
```bash
go run ./cmd/seeder-mapper -v
```
Enables verbose output for debugging purposes.

## Seeder Requirements

For a struct to be automatically detected and registered, it must:

1. Be located in the `pkg/seeder` directory
2. Be in a file with the naming pattern `*_seeder.go`
3. Have a struct name ending with "Seeder" (e.g., `AnimalSeeder`, `UserSeeder`)
4. Implement the `Seeder` interface with these methods:
   - `Seed(ctx context.Context) error`
   - `GetName() string`
5. Have a corresponding constructor function named `New{Name}Seeder` (e.g., `NewAnimalSeeder`)

## Example Seeder Structure

```go
package seeder

import (
    "context"
    "github.com/linkeunid/go-api/pkg/database"
    "go.uber.org/zap"
)

// ProductSeeder seeds product data
type ProductSeeder struct {
    db     database.Database
    logger *zap.Logger
    count  int
}

// NewProductSeeder creates a new product seeder
func NewProductSeeder(db database.Database, logger *zap.Logger, count int) *ProductSeeder {
    return &ProductSeeder{
        db:     db,
        logger: logger,
        count:  count,
    }
}

// GetName returns the name of the seeder
func (s *ProductSeeder) GetName() string {
    return "product"
}

// Seed seeds product data
func (s *ProductSeeder) Seed(ctx context.Context) error {
    // Implementation here
    return nil
}
```

## How It Works

1. **File Discovery**: Scans `pkg/seeder/*.go` files (excluding test files)
2. **AST Parsing**: Uses Go's AST parser to analyze struct declarations
3. **Interface Validation**: Checks if structs implement the required `Seeder` interface methods
4. **Registry Update**: Modifies the `registerSeeders` function in `cmd/seed/main.go`

## Generated Code

The tool generates/updates the `registerSeeders` function in this format:

```go
func registerSeeders(db database.Database, logger *zap.Logger, count int) []Seeder {
    return []Seeder{
        seeder.NewAnimalSeeder(db, logger, count),
        seeder.NewFlowerSeeder(db, logger, count),
        seeder.NewProductSeeder(db, logger, count),
        // Add more seeders here as they are implemented
    }
}
```

## Benefits

- **Automatic Registration**: No need to manually update the registry when adding new seeders
- **Safe Cleanup**: Automatically removes references to deleted seeders
- **Error Prevention**: Ensures only valid seeders that implement the interface are registered
- **Consistent Naming**: Enforces naming conventions for seeder files and constructors
- **CI/CD Integration**: Can be run in build pipelines to ensure seeders are always up-to-date

## Integration with Makefile

The tool is integrated into the project's Makefile with three commands:

- `make update-seeder-registry` (alias: `usr`) - Add new seeders only
- `make sync-seeder` (alias: `ss`) - Add new and remove deleted seeders
- `make clean-seeder-registry` (alias: `csr`) - Remove deleted seeders only

It's recommended to run `make sync-seeder` regularly during development to keep the registry in sync with your seeder files. 