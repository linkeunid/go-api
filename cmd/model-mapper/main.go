// Package main provides a command line tool for updating the model map in cmd/db/main.go
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
	flag.BoolVar(&cleanOnly, "clean-only", false, "Only remove models that no longer exist without adding new ones")
	flag.BoolVar(&syncMode, "sync", false, "Both add new models and remove models that no longer exist")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()
}

func main() {
	// Scan for models in the filesystem
	fsModels, err := scanModels("internal/model")
	if err != nil {
		fmt.Printf("Error scanning models: %v\n", err)
		os.Exit(1)
	}

	// Get the current models in the model map
	currentModels, err := getCurrentModels("cmd/db/main.go")
	if err != nil {
		fmt.Printf("Error getting current models: %v\n", err)
		os.Exit(1)
	}

	// Determine which models to keep, add, or remove
	var modelsToUse []string
	var added, removed []string

	if cleanOnly || syncMode {
		// For clean or sync mode, determine which models to keep
		for _, model := range currentModels {
			// Keep model if it still exists in filesystem
			if contains(fsModels, model) {
				modelsToUse = append(modelsToUse, model)
			} else {
				removed = append(removed, model)
			}
		}

		// For sync mode, also add new models
		if syncMode {
			for _, model := range fsModels {
				if !contains(currentModels, model) {
					modelsToUse = append(modelsToUse, model)
					added = append(added, model)
				}
			}
		}
	} else {
		// For regular update mode, use all models found in filesystem
		modelsToUse = fsModels

		// Calculate added models for reporting
		for _, model := range fsModels {
			if !contains(currentModels, model) {
				added = append(added, model)
			}
		}
	}

	// Update the model map in cmd/db/main.go
	if err := updateModelMap("cmd/db/main.go", modelsToUse); err != nil {
		fmt.Printf("Error updating model map: %v\n", err)
		os.Exit(1)
	}

	// Print summary based on the mode
	if cleanOnly {
		if len(removed) > 0 {
			fmt.Printf("✅ Cleaned model map: removed %d models\n", len(removed))
			for _, model := range removed {
				fmt.Printf("  - Removed: %s\n", model)
			}
		} else {
			fmt.Println("✅ No models needed to be removed")
		}
	} else if syncMode {
		fmt.Printf("✅ Synced model map with %d models (added: %d, removed: %d)\n",
			len(modelsToUse), len(added), len(removed))

		if len(added) > 0 {
			fmt.Println("  Added models:")
			for _, model := range added {
				fmt.Printf("    - %s\n", model)
			}
		}

		if len(removed) > 0 {
			fmt.Println("  Removed models:")
			for _, model := range removed {
				fmt.Printf("    - %s\n", model)
			}
		}
	} else {
		fmt.Printf("✅ Updated model map with %d models\n", len(modelsToUse))
		if len(added) > 0 {
			fmt.Println("  Newly added models:")
			for _, model := range added {
				fmt.Printf("    - %s\n", model)
			}
		}

		fmt.Println("  All current models:")
		for _, model := range modelsToUse {
			fmt.Printf("    - %s\n", model)
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

// getCurrentModels extracts the current model names from the model map in the given file
func getCurrentModels(filePath string) ([]string, error) {
	var models []string

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract model names from the model map
	// Look for patterns like: "model_name": &model.ModelName{},
	pattern := regexp.MustCompile(`"([a-z0-9_]+)":\s*&model\.([A-Za-z0-9_]+)\{\}`)
	matches := pattern.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 2 {
			// Return the struct name (PascalCase), not the key (snake_case)
			models = append(models, match[2])
		}
	}

	return models, nil
}

// scanModels scans the model directory and returns a list of model struct names
func scanModels(modelDir string) ([]string, error) {
	var models []string

	// Get all Go files in the model directory
	files, err := filepath.Glob(filepath.Join(modelDir, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("failed to list model files: %w", err)
	}

	// Parse each file and extract struct declarations that implement TableName()
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

					// Check if this struct has a TableName method in the file
					if hasTableNameMethod(file, structName) {
						models = append(models, structName)
					}
				}
			}
		}
	}

	return models, nil
}

// hasTableNameMethod checks if a struct has a TableName method in the file
func hasTableNameMethod(filePath, structName string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	// Look for function declarations like: func (ModelName) TableName() string
	pattern := fmt.Sprintf(`func\s*\(\s*(%s|[a-zA-Z0-9_]+\s+%s)\s*\)\s*TableName\s*\(\s*\)\s*string`,
		structName, structName)
	matched, _ := regexp.MatchString(pattern, string(content))
	return matched
}

// toSnakeCase converts a string from PascalCase to snake_case
func toSnakeCase(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		if r >= 'A' && r <= 'Z' {
			result += string(r - 'A' + 'a')
		} else {
			result += string(r)
		}
	}
	return result
}

// updateModelMap updates the modelMap in the specified Go file
func updateModelMap(filePath string, models []string) error {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Define a more precise pattern to match the model map declaration
	// Match from "var modelMap" up to the final closing brace
	pattern := regexp.MustCompile(`var\s+modelMap\s*=\s*map\[string\]interface\{\}\s*\{[\s\S]*?\n\}`)

	// Generate the new model map with proper formatting
	var newMap bytes.Buffer
	newMap.WriteString("var modelMap = map[string]interface{}{\n")
	for _, model := range models {
		// Convert PascalCase to snake_case for the key
		modelKey := toSnakeCase(model)
		// Use the original PascalCase struct name
		newMap.WriteString(fmt.Sprintf("\t\"%s\": &model.%s{},\n", modelKey, model))
	}
	newMap.WriteString("\t// Add more models here as they are implemented\n}")

	// Replace the model map declaration
	newContent := pattern.ReplaceAllString(string(content), newMap.String())

	// Write the updated content back to the file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
