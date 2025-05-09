package bootstrap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/linkeunid/go-api/internal/controller"
	swaggerdocs "github.com/linkeunid/go-api/internal/docs/swaggerdocs"
	"github.com/linkeunid/go-api/pkg/auth"
	"github.com/linkeunid/go-api/pkg/config"
	custommiddleware "github.com/linkeunid/go-api/pkg/middleware"
	"github.com/linkeunid/go-api/pkg/response"
	"github.com/linkeunid/go-api/pkg/util"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// SetupSwagger configures the Swagger documentation
func SetupSwagger(port int, isDevelopment bool) {
	if isDevelopment {
		// Set basic Swagger info
		swaggerdocs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
		swaggerdocs.SwaggerInfo.Title = "Linkeun Go API"
		swaggerdocs.SwaggerInfo.Description = "API for managing various resources including animals"
		swaggerdocs.SwaggerInfo.Version = "1.0"
		swaggerdocs.SwaggerInfo.BasePath = "/api/v1"
		swaggerdocs.SwaggerInfo.Schemes = []string{"http", "https"}

		// We're using annotation-based Swagger docs, so we need to rely on the following
		// annotations in the swagger.go file:
		//
		// @contact.name API Support - Website
		// @contact.url http://linkeunid.com/support
		// @contact.email Send email to API Support
		//
		// @license.name GNU General Public License v2.0
		// @license.url https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html
	}
}

// SetupServer configures and returns an HTTP server with all routes and middleware
func SetupServer(app *App, animalController *controller.Animal) *http.Server {
	logger := app.Logger
	cfg := app.Config

	// Initialize router
	r := chi.NewRouter()

	// Create JWT service
	jwtService := auth.NewJWTService(&cfg.Auth)

	// Create auth middleware
	authMiddleware := custommiddleware.NewAuthMiddleware(jwtService, &cfg.Auth, logger)

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(custommiddleware.ValidationMiddleware) // Add our custom validation middleware

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger documentation - only available in development mode
	if cfg.IsDevelopment() {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // The URL points to API definition
		))
		logger.Info("Swagger UI enabled in development mode")
	}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/public", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				data := map[string]string{
					"message": "This is a public endpoint that doesn't require authentication",
				}
				response.Success(w, r, data, "Public API is working")
			})
		})

		// Protected routes (require authentication)
		r.Route("/protected", func(r chi.Router) {
			// Apply authentication middleware to all routes in this group
			r.Use(authMiddleware.Authenticate)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// Extract user information from context
				userID := r.Context().Value(custommiddleware.KeyUserID)
				username := r.Context().Value(custommiddleware.KeyUsername)
				role := r.Context().Value(custommiddleware.KeyUserRole)
				email := r.Context().Value(custommiddleware.KeyUserEmail)

				data := map[string]interface{}{
					"message": "This endpoint requires authentication",
					"user": map[string]interface{}{
						"id":       userID,
						"username": username,
						"role":     role,
						"email":    email,
					},
				}
				response.Success(w, r, data, "Protected API is working")
			})

			// Admin-only routes
			r.Route("/admin", func(r chi.Router) {
				// Apply role-based middleware
				r.Use(authMiddleware.RequireRole("admin"))

				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					data := map[string]string{
						"message": "This endpoint requires admin role",
					}
					response.Success(w, r, data, "Admin API is working")
				})
			})
		})

		// Animal routes
		animalController.RegisterRoutes(r)
	})

	// Create and return server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return server
}

// LogServerInfo logs server startup information
func LogServerInfo(logger *zap.Logger, port int, isDevelopment bool, config *config.Config) {
	logger.Info("Starting server", zap.Int("port", port))
	logger.Info("Auth configuration", zap.Bool("enabled", config.Auth.Enabled))

	// Log information about file logging if enabled
	if config.Logging.FileOutputPath != "" {
		logger.Info("File logging enabled",
			zap.String("path", config.Logging.FileOutputPath),
			zap.Int("maxSize", config.Logging.FileMaxSize),
			zap.Int("maxBackups", config.Logging.FileMaxBackups),
			zap.Int("maxAge", config.Logging.FileMaxAge),
			zap.Bool("compress", config.Logging.FileCompress))
	}

	if isDevelopment {
		logger.Info("Swagger UI available at", zap.String("url", fmt.Sprintf("http://localhost:%d/swagger/", port)))
	}
}

// GetDataSourceInfo returns masked DSN for logging
func GetDataSourceInfo(dsn string) string {
	return util.MaskDsn(dsn)
}

// MaskSensitiveData masks sensitive data for logging purposes
func MaskSensitiveData(dataType string, value string) string {
	switch dataType {
	case "email":
		return util.MaskEmail(value)
	case "credential", "password":
		return util.MaskCredential(value)
	case "jwt", "token":
		return util.MaskJWT(value)
	case "url":
		return util.MaskURL(value)
	case "dsn":
		return util.MaskDsn(value)
	default:
		return util.MaskSensitive(value, 4, 0) // default mask showing only first 4 chars
	}
}
