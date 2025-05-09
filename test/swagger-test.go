package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/linkeunid/go-api/internal/docs/swaggerdocs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// Get parameters from command line arguments or use defaults
	host := "localhost"
	port := "8090"
	apiHost := os.Getenv("API_HOST") // Actual API host to use in Swagger UI
	apiPort := os.Getenv("API_PORT") // Actual API port to use in Swagger UI

	if apiHost == "" {
		apiHost = "0.0.0.0"
	}

	if apiPort == "" {
		apiPort = "8080"
	}

	// Override from command line args
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	if len(os.Args) > 3 {
		apiHost = os.Args[3]
	}

	if len(os.Args) > 4 {
		apiPort = os.Args[4]
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// For local development, if API host is 0.0.0.0, change to localhost for browser access
	if apiHost == "0.0.0.0" {
		apiHost = "localhost"
	}

	// Set the Swagger host to point to the real API
	swaggerdocs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", apiHost, apiPort)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Add a header middleware to set specific configuration for the Swagger UI
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers specifically for swagger requests
			if r.URL.Path == "/swagger/doc.json" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			next.ServeHTTP(w, r)
		})
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL points to API definition
	))

	// Simple route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Visit /swagger/ for API documentation")
	})

	// Add a test route to verify CORS is working
	r.Get("/cors-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "CORS test successful"}`))
	})

	// Start server
	fmt.Printf("Starting Swagger UI server on %s\n", addr)
	fmt.Printf("Swagger UI available at http://%s/swagger/\n", addr)
	fmt.Printf("Configured to use API at http://%s:%s/\n", apiHost, apiPort)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-c
	fmt.Println("\nShutting down server...")
	os.Exit(0)
}
