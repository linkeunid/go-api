// Package main provides a command line tool for updating the seeder registry in cmd/seed/main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Command line flags
var (
	cleanOnly bool
	syncMode  bool
	verbose   bool
)

func init() {
	flag.BoolVar(&cleanOnly, "clean-only", false, "Only remove seeders that no longer exist without adding new ones")
	flag.BoolVar(&syncMode, "sync", false, "Both add new seeders and remove seeders that no longer exist")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()
}

func main() {
	// Scan for seeders in the filesystem
	fsSeeders, err := scanSeeders("pkg/seeder")
	if err != nil {
		fmt.Printf("Error scanning seeders: %v\n", err)
		os.Exit(1)
	}

	// Get the current seeders in the seeder registry
	currentSeeders, err := getCurrentSeeders("cmd/seed/main.go")
	if err != nil {
		fmt.Printf("Error getting current seeders: %v\n", err)
		os.Exit(1)
	}

	// Determine which seeders to use, add, or remove
	var seedersToUse []string
	var added, removed []string

	if cleanOnly || syncMode {
		// For clean or sync mode, determine which seeders to keep
		for _, seeder := range currentSeeders {
			// Keep seeder if it still exists in filesystem
			if contains(fsSeeders, seeder) {
				seedersToUse = append(seedersToUse, seeder)
			} else {
				removed = append(removed, seeder)
			}
		}

		// For sync mode, also add new seeders
		if syncMode {
			for _, seeder := range fsSeeders {
				if !contains(currentSeeders, seeder) {
					seedersToUse = append(seedersToUse, seeder)
					added = append(added, seeder)
				}
			}
		}
	} else {
		// For regular update mode, use all seeders found in filesystem
		seedersToUse = fsSeeders

		// Calculate added seeders for reporting
		for _, seeder := range fsSeeders {
			if !contains(currentSeeders, seeder) {
				added = append(added, seeder)
			}
		}
	}

	// Update the seeder registry in cmd/seed/main.go
	if err := updateSeederRegistry("cmd/seed/main.go", seedersToUse); err != nil {
		fmt.Printf("Error updating seeder registry in cmd/seed/main.go: %v\n", err)
		os.Exit(1)
	}

	// Print summary based on the mode
	if cleanOnly {
		if len(removed) > 0 {
			fmt.Printf("✅ Cleaned seeder registry: removed %d seeders\n", len(removed))
			for _, seeder := range removed {
				fmt.Printf("  - Removed: %s\n", seeder)
			}
		} else {
			fmt.Println("✅ No seeders needed to be removed")
		}
	} else if syncMode {
		fmt.Printf("✅ Synced seeder registry with %d seeders (added: %d, removed: %d)\n",
			len(seedersToUse), len(added), len(removed))

		if len(added) > 0 {
			fmt.Println("  Added seeders:")
			for _, seeder := range added {
				fmt.Printf("    - %s\n", seeder)
			}
		}

		if len(removed) > 0 {
			fmt.Println("  Removed seeders:")
			for _, seeder := range removed {
				fmt.Printf("    - %s\n", seeder)
			}
		}
	} else {
		fmt.Printf("✅ Updated seeder registry with %d seeders\n", len(seedersToUse))
		if len(added) > 0 {
			fmt.Println("  Newly added seeders:")
			for _, seeder := range added {
				fmt.Printf("    - %s\n", seeder)
			}
		}

		fmt.Println("  All current seeders:")
		for _, seeder := range seedersToUse {
			fmt.Printf("    - %s\n", seeder)
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// getCurrentSeeders extracts the current seeder names from the seeder registry in the given file
func getCurrentSeeders(filePath string) ([]string, error) {
	var seeders []string

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract seeder names from the registerSeeders function
	// Look for patterns like: seeder.NewAnimalSeeder(db, logger, count),
	pattern := regexp.MustCompile(`seeder\.New([A-Za-z0-9_]+)Seeder\(`)
	matches := pattern.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			// Return the seeder name (e.g., "Animal" from "NewAnimalSeeder")
			seeders = append(seeders, match[1])
		}
	}

	return seeders, nil
}

// scanSeeders scans the seeder directory and returns a list of seeder struct names
func scanSeeders(seederDir string) ([]string, error) {
	var seeders []string

	// Get all Go files in the seeder directory
	files, err := filepath.Glob(filepath.Join(seederDir, "*_seeder.go"))
	if err != nil {
		return nil, fmt.Errorf("failed to list seeder files: %w", err)
	}

	// Parse each file and extract struct declarations that implement Seeder interface
	fset := token.NewFileSet()
	for _, file := range files {
		// Skip test files
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		// Parse the Go file
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}

		// Find struct declarations
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if it's a struct
				if _, ok := typeSpec.Type.(*ast.StructType); ok {
					structName := typeSpec.Name.Name

					// Check if this struct has Seed and GetName methods (Seeder interface)
					if hasSeederMethods(file, structName) {
						// Extract the seeder name by removing "Seeder" suffix
						if strings.HasSuffix(structName, "Seeder") {
							seederName := strings.TrimSuffix(structName, "Seeder")
							seeders = append(seeders, seederName)
						}
					}
				}
			}
		}
	}

	return seeders, nil
}

// hasSeederMethods checks if a struct has Seed and GetName methods (implements Seeder interface)
func hasSeederMethods(filePath, structName string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	// Look for function declarations like:
	// func (s *StructName) Seed(ctx context.Context) error
	// func (s *StructName) GetName() string
	seedPattern := fmt.Sprintf(`func\s*\(\s*[a-zA-Z0-9_]+\s+\*?%s\s*\)\s*Seed\s*\(\s*ctx\s+context\.Context\s*\)\s*error`,
		structName)
	getNamePattern := fmt.Sprintf(`func\s*\(\s*[a-zA-Z0-9_]+\s+\*?%s\s*\)\s*GetName\s*\(\s*\)\s*string`,
		structName)

	hasSeed, _ := regexp.MatchString(seedPattern, string(content))
	hasGetName, _ := regexp.MatchString(getNamePattern, string(content))

	return hasSeed && hasGetName
}

// updateSeederRegistry updates the registerSeeders function in the specified Go file
func updateSeederRegistry(filePath string, seeders []string) error {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Define a pattern to match the registerSeeders function
	// Match from "func registerSeeders" up to the final closing brace
	pattern := regexp.MustCompile(`func\s+registerSeeders\([^)]*\)\s+\[\]Seeder\s+\{[\s\S]*?\n\}`)

	// Generate the new registerSeeders function with proper formatting
	var newRegistry bytes.Buffer
	newRegistry.WriteString("func registerSeeders(db database.Database, logger *zap.Logger, count int) []Seeder {\n")
	newRegistry.WriteString("\treturn []Seeder{\n")

	for _, seeder := range seeders {
		newRegistry.WriteString(fmt.Sprintf("\t\tseeder.New%sSeeder(db, logger, count),\n", seeder))
	}

	newRegistry.WriteString("\t\t// Add more seeders here as they are implemented\n")
	newRegistry.WriteString("\t}\n}")

	// Replace the registerSeeders function
	newContent := pattern.ReplaceAllString(string(content), newRegistry.String())

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
