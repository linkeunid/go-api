// Package main provides a command line tool for setting up a new project from the template
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Constants
const (
	currentModuleName = "github.com/linkeunid/go-api"
)

// Command line flags
var (
	newModuleName string
	gitRemoteURL  string
	resetGit      bool
	verbose       bool
)

func init() {
	flag.StringVar(&newModuleName, "module", "", "New module name (e.g., github.com/yourusername/your-project)")
	flag.StringVar(&gitRemoteURL, "remote", "", "Git remote URL (e.g., git@github.com:yourusername/your-project.git)")
	flag.BoolVar(&resetGit, "reset-git", false, "Reset Git repository (remove .git folder and initialize a new one)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.Parse()
}

func main() {
	// Validate flags
	if newModuleName == "" {
		fmt.Println("❌ Error: New module name is required. Use -module flag.")
		fmt.Println("Example: go run ./cmd/setup-project -module github.com/yourusername/your-project")
		os.Exit(1)
	}

	// Confirm with the user before proceeding
	confirmed := confirmAction()
	if !confirmed {
		fmt.Println("❌ Operation cancelled.")
		os.Exit(0)
	}

	// Start the rename process
	fmt.Printf("🔄 Setting up project with new module name: %s\n", newModuleName)

	// Rename module in go.mod
	renameModuleInGoMod()

	// Update import paths in all Go files
	updateImportPaths()

	// Handle Git repository
	handleGitRepository()

	fmt.Println("✅ Project setup completed successfully!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. Review the changes to ensure everything was updated correctly")
	fmt.Println("2. Run 'go mod tidy' to update dependencies")
	fmt.Println("3. Build and test your project to verify everything works")
}

// confirmAction asks the user to confirm the operation
func confirmAction() bool {
	fmt.Println("⚠️ WARNING: This operation will:")
	fmt.Printf("  - Rename module from %s to %s\n", currentModuleName, newModuleName)
	fmt.Println("  - Update all import paths in Go files")

	if resetGit {
		fmt.Println("  - Reset Git repository (remove .git folder and initialize a new one)")
	}

	if gitRemoteURL != "" {
		fmt.Printf("  - Set Git remote origin to %s\n", gitRemoteURL)
	}

	fmt.Println("\nThis operation cannot be undone. Do you want to continue? (y/n)")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes"
}

// renameModuleInGoMod updates the module name in go.mod file
func renameModuleInGoMod() {
	fmt.Println("📝 Updating go.mod file...")

	// Read go.mod file
	goModPath := "go.mod"
	content, err := os.ReadFile(goModPath)
	if err != nil {
		fmt.Printf("❌ Error reading go.mod: %v\n", err)
		os.Exit(1)
	}

	// Replace module name
	newContent := regexp.MustCompile(`module\s+`+regexp.QuoteMeta(currentModuleName)).
		ReplaceAll(content, []byte("module "+newModuleName))

	// Write updated content back to go.mod
	err = os.WriteFile(goModPath, newContent, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing go.mod: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Println("  ✓ go.mod updated")
	}
}

// updateImportPaths updates import paths in all Go files
func updateImportPaths() {
	fmt.Println("📝 Updating import paths in all Go files...")

	// Get all Go files in the project
	var goFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor directory and .git directory
		if info.IsDir() && (info.Name() == "vendor" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		// Process only Go files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error scanning files: %v\n", err)
		os.Exit(1)
	}

	// Process each Go file
	for _, file := range goFiles {
		updateImportsInFile(file)
	}
}

// updateImportsInFile updates import paths in a single Go file
func updateImportsInFile(filePath string) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading %s: %v\n", filePath, err)
		return
	}

	// Replace import paths
	oldImportPattern := regexp.QuoteMeta(currentModuleName)
	newContent := regexp.MustCompile(`"`+oldImportPattern+`(/[^"]*)?"`).
		ReplaceAll(content, []byte(`"`+newModuleName+`$1"`))

	// If content hasn't changed, skip writing
	if bytes.Equal(content, newContent) {
		if verbose {
			fmt.Printf("  - Skipped %s (no changes needed)\n", filePath)
		}
		return
	}

	// Write updated content back to file
	err = os.WriteFile(filePath, newContent, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing %s: %v\n", filePath, err)
		return
	}

	if verbose {
		fmt.Printf("  ✓ Updated %s\n", filePath)
	}
}

// handleGitRepository handles Git repository operations
func handleGitRepository() {
	// Reset Git repository if requested
	if resetGit {
		fmt.Println("🔄 Resetting Git repository...")

		// Remove .git directory
		err := os.RemoveAll(".git")
		if err != nil {
			fmt.Printf("❌ Error removing .git directory: %v\n", err)
			os.Exit(1)
		}

		// Initialize new Git repository
		cmd := exec.Command("git", "init")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error initializing Git repository: %v\n", err)
			fmt.Println(string(output))
			os.Exit(1)
		}

		if verbose {
			fmt.Println("  ✓ Git repository reset")
		}
	}

	// Set Git remote if provided
	if gitRemoteURL != "" {
		fmt.Printf("🔄 Setting Git remote origin to %s\n", gitRemoteURL)

		var cmd *exec.Cmd
		var output []byte
		var err error

		// Check if remote exists
		cmd = exec.Command("git", "remote")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error checking Git remotes: %v\n", err)
			fmt.Println(string(output))
			os.Exit(1)
		}

		remotes := strings.Split(string(output), "\n")
		originExists := false
		for _, remote := range remotes {
			if strings.TrimSpace(remote) == "origin" {
				originExists = true
				break
			}
		}

		// Add or set remote
		if originExists {
			cmd = exec.Command("git", "remote", "set-url", "origin", gitRemoteURL)
		} else {
			cmd = exec.Command("git", "remote", "add", "origin", gitRemoteURL)
		}

		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error setting Git remote: %v\n", err)
			fmt.Println(string(output))
			os.Exit(1)
		}

		if verbose {
			fmt.Println("  ✓ Git remote set to origin")
		}
	}
}
