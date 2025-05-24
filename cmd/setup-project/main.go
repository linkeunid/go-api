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
	skipConfirm   bool
)

func init() {
	flag.StringVar(&newModuleName, "module", "", "New module name (e.g., github.com/yourusername/your-project)")
	flag.StringVar(&gitRemoteURL, "remote", "", "Git remote URL (e.g., git@github.com:yourusername/your-project.git)")
	flag.BoolVar(&resetGit, "reset-git", false, "Reset Git repository (remove .git folder and initialize a new one)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.BoolVar(&skipConfirm, "y", false, "Skip confirmation prompt (use with caution)")
	flag.Parse()
}

func main() {
	// Validate flags
	if newModuleName == "" {
		fmt.Println("❌ Error: New module name is required. Use -module flag.")
		fmt.Println("Example: go run ./cmd/setup-project -module github.com/yourusername/your-project")
		os.Exit(1)
	}

	// Confirm with the user before proceeding (if not skipped)
	if !skipConfirm {
		confirmed := confirmAction()
		if !confirmed {
			fmt.Println("❌ Operation cancelled.")
			os.Exit(0)
		}
	}

	// Start the rename process
	fmt.Printf("🔄 Setting up project with new module name: %s\n", newModuleName)

	// Rename module in go.mod
	renameModuleInGoMod()

	// Update import paths in all Go files
	updateImportPaths()

	// Update docker-compose.yml with new service and container names
	updateDockerCompose()

	// Handle Git repository
	handleGitRepository()

	fmt.Println("✅ Project setup completed successfully!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. Review the changes to ensure everything was updated correctly")
	fmt.Println("2. Run 'go mod tidy' to update dependencies")
	fmt.Println("3. Build and test your project to verify everything works")
	fmt.Println("4. Update your .env file if needed to match the new service names")
}

// extractProjectName extracts the project name from the module path
// e.g., "github.com/yourusername/your-project" -> "your-project"
func extractProjectName(modulePath string) string {
	parts := strings.Split(modulePath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "go-api" // fallback
}

// confirmAction asks the user to confirm the operation
func confirmAction() bool {
	projectName := extractProjectName(newModuleName)

	fmt.Println("⚠️ WARNING: This operation will:")
	fmt.Printf("  - Rename module from %s to %s\n", currentModuleName, newModuleName)
	fmt.Println("  - Update all import paths in Go files")
	fmt.Printf("  - Update docker-compose.yml service names (api -> %s-api, mysql -> %s-mysql, redis -> %s-redis)\n", projectName, projectName, projectName)
	fmt.Printf("  - Update docker-compose.yml container names accordingly\n")

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

// updateDockerCompose updates service names and container names in docker-compose.yml
func updateDockerCompose() {
	fmt.Println("🐳 Updating docker-compose.yml with new service and container names...")

	dockerComposePath := "docker-compose.yml"

	// Check if docker-compose.yml exists
	if _, err := os.Stat(dockerComposePath); os.IsNotExist(err) {
		if verbose {
			fmt.Println("  - Skipped docker-compose.yml (file not found)")
		}
		return
	}

	// Read docker-compose.yml file
	content, err := os.ReadFile(dockerComposePath)
	if err != nil {
		fmt.Printf("❌ Error reading docker-compose.yml: %v\n", err)
		return
	}

	projectName := extractProjectName(newModuleName)
	originalContent := content

	// Define service mapping: old name -> new name
	serviceMap := map[string]string{
		"api":   projectName + "-api",
		"mysql": projectName + "-mysql",
		"redis": projectName + "-redis",
	}

	// Define container mapping: old name -> new name
	containerMap := map[string]string{
		"go-api":        projectName + "-api",
		"go-mysql":      projectName + "-mysql",
		"linkeun-redis": projectName + "-redis",
	}

	// Update service names (as top-level services)
	for oldService, newService := range serviceMap {
		// Match service definitions (e.g., "  api:" or "services:\n  api:")
		servicePattern := regexp.MustCompile(`(\s+)` + regexp.QuoteMeta(oldService) + `:`)
		content = servicePattern.ReplaceAll(content, []byte(`$1`+newService+`:`))
	}

	// Update container names
	for oldContainer, newContainer := range containerMap {
		// Match container_name lines (e.g., "container_name: go-api")
		containerPattern := regexp.MustCompile(`container_name:\s*` + regexp.QuoteMeta(oldContainer))
		content = containerPattern.ReplaceAll(content, []byte(`container_name: `+newContainer))
	}

	// Update service references in depends_on and other cross-references
	for oldService, newService := range serviceMap {
		// Update depends_on references (e.g., "mysql:" under depends_on)
		dependsPattern := regexp.MustCompile(`(\s+` + regexp.QuoteMeta(oldService) + `:)(\s+condition:)`)
		content = dependsPattern.ReplaceAll(content, []byte(`      `+newService+`:$2`))

		// Update service references in environment variables (e.g., REDIS_HOST=redis)
		envPattern := regexp.MustCompile(`(=)` + regexp.QuoteMeta(oldService) + `(\s|$)`)
		content = envPattern.ReplaceAll(content, []byte(`$1`+newService+`$2`))
	}

	// Update network name
	oldNetworkName := "linkeun-network"
	newNetworkName := projectName + "-network"

	// Update network definition
	networkDefPattern := regexp.MustCompile(`(\s+)` + regexp.QuoteMeta(oldNetworkName) + `:`)
	content = networkDefPattern.ReplaceAll(content, []byte(`$1`+newNetworkName+`:`))

	// Update network references
	networkRefPattern := regexp.MustCompile(`- ` + regexp.QuoteMeta(oldNetworkName))
	content = networkRefPattern.ReplaceAll(content, []byte(`- `+newNetworkName))

	// Update volume names to include project name
	volumeMap := map[string]string{
		"mysql_data": projectName + "_mysql_data",
		"redis_data": projectName + "_redis_data",
	}

	for oldVolume, newVolume := range volumeMap {
		// Update volume definitions
		volumeDefPattern := regexp.MustCompile(`(\s+)` + regexp.QuoteMeta(oldVolume) + `:`)
		content = volumeDefPattern.ReplaceAll(content, []byte(`$1`+newVolume+`:`))

		// Update volume references
		volumeRefPattern := regexp.MustCompile(`- ` + regexp.QuoteMeta(oldVolume) + `:`)
		content = volumeRefPattern.ReplaceAll(content, []byte(`- `+newVolume+`:`))
	}

	// If content hasn't changed, skip writing
	if bytes.Equal(originalContent, content) {
		if verbose {
			fmt.Println("  - Skipped docker-compose.yml (no changes needed)")
		}
		return
	}

	// Write updated content back to docker-compose.yml
	err = os.WriteFile(dockerComposePath, content, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing docker-compose.yml: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("  ✓ Updated docker-compose.yml with project name: %s\n", projectName)
		fmt.Printf("    - Services: api -> %s-api, mysql -> %s-mysql, redis -> %s-redis\n", projectName, projectName, projectName)
		fmt.Printf("    - Containers: go-api -> %s-api, go-mysql -> %s-mysql, linkeun-redis -> %s-redis\n", projectName, projectName, projectName)
		fmt.Printf("    - Network: linkeun-network -> %s-network\n", projectName)
		fmt.Printf("    - Volumes: mysql_data -> %s_mysql_data, redis_data -> %s_redis_data\n", projectName, projectName)
	} else {
		fmt.Printf("  ✓ Updated docker-compose.yml with project name: %s\n", projectName)
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
