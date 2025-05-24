package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/linkeunid/go-api/pkg/auth"
	"github.com/linkeunid/go-api/pkg/config"
)

func main() {
	// Load configuration first
	cfg := config.LoadConfig()

	// Define command-line flags
	var (
		userID   uint64
		username string
		role     string
		email    string
		secret   string
		expire   time.Duration
		force    bool
	)

	// Parse command-line arguments
	flag.Uint64Var(&userID, "id", 1, "User ID")
	flag.StringVar(&username, "username", "testuser", "Username")
	flag.StringVar(&role, "role", "user", "User role (user, admin, etc.)")
	flag.StringVar(&email, "email", "test@example.com", "User email")
	flag.StringVar(&secret, "secret", cfg.Auth.JWTSecret, "JWT secret key (defaults to JWT_SECRET env var)")
	flag.DurationVar(&expire, "expire", cfg.Auth.JWTExpiration, "Token expiration duration (e.g., 24h, 30m)")
	flag.BoolVar(&force, "force", false, "Force token generation even in production (use with caution)")
	flag.Parse()

	// Check if the environment is development or test
	env := cfg.Environment
	if env != "development" && env != "test" && env != "" && !force {
		fmt.Println("❌ Error: Token generation is only available in development and test environments.")
		fmt.Printf("Current environment: %s\n", env)
		fmt.Println("Set APP_ENV to 'development' or 'test' to use this command.")
		fmt.Println("Or use the --force flag to override this restriction (use with caution).")
		os.Exit(1)
	}

	// Allow force override for special cases, but with a warning
	if env == "production" && force {
		fmt.Println("⚠️  WARNING: Forcing token generation in production environment!")
		fmt.Println("This should only be done in exceptional circumstances.")
		fmt.Println()
	}

	// Validate required parameters
	if secret == "" {
		fmt.Println("Error: JWT secret is required. Set JWT_SECRET env var or use --secret flag.")
		os.Exit(1)
	}

	// Create auth config with custom values if provided
	authConfig := &config.AuthConfig{
		Enabled:       true,
		JWTSecret:     secret,
		JWTExpiration: expire,
	}

	// Create JWT service
	jwtService := auth.NewJWTService(authConfig)

	// Generate token
	token, err := jwtService.GenerateToken(userID, username, role, email)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		os.Exit(1)
	}

	// Print the token with usage examples
	fmt.Println("\n🔑 JWT Token Generated 🔑")
	if env == "production" {
		fmt.Println("⚠️  WARNING: Token generated in production environment!")
	}
	fmt.Println("-------------------------------")
	fmt.Printf("Token: %s\n", token)
	fmt.Println("-------------------------------")
	fmt.Println("\nClaims:")
	fmt.Printf("  User ID (sub): %d\n", userID)
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Role: %s\n", role)
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  Expires: %s\n", time.Now().Add(expire).Format(time.RFC1123))
	fmt.Printf("  Environment: %s\n", env)
	fmt.Println("\nUsage Examples:")
	fmt.Println("  cURL:")
	fmt.Printf("    curl -H \"Authorization: Bearer %s\" http://localhost:%d/api/v1/protected\n", token, cfg.Server.Port)
	fmt.Println("\n  JavaScript Fetch:")
	fmt.Printf("    fetch('http://localhost:%d/api/v1/protected', {\n      headers: {\n        'Authorization': 'Bearer %s'\n      }\n    })\n", cfg.Server.Port, token)

	// Add information about available protected endpoints
	fmt.Println("\nAvailable Protected Endpoints:")
	fmt.Println("  1. General User Endpoint:")
	fmt.Println("     - URL: /api/v1/protected")
	fmt.Println("     - Method: GET")
	fmt.Println("     - Access: Any authenticated user")
	fmt.Println("     - Returns: User information (ID, username, role, email)")

	fmt.Println("\n  2. Admin-Only Endpoint:")
	fmt.Println("     - URL: /api/v1/protected/admin")
	fmt.Println("     - Method: GET")
	fmt.Println("     - Access: Only users with admin role")
	fmt.Println("     - Returns: Admin-specific information")

	fmt.Println("\n  3. Animal Resource Endpoints:")
	fmt.Println("     - GET    /api/v1/animals            (List all animals)")
	fmt.Println("     - POST   /api/v1/animals            (Create a new animal)")
	fmt.Println("     - GET    /api/v1/animals/{animalID} (Get a specific animal)")
	fmt.Println("     - PUT    /api/v1/animals/{animalID} (Update an animal)")
	fmt.Println("     - DELETE /api/v1/animals/{animalID} (Delete an animal)")
}
