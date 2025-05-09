package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/linkeunid/go-api/pkg/auth"
	"github.com/linkeunid/go-api/pkg/config"
	"github.com/linkeunid/go-api/pkg/response"
	"go.uber.org/zap"
)

// ContextKey type for context keys to avoid collisions
type ContextKey string

// Context keys for user information
const (
	// KeyUserID is the context key for user ID
	KeyUserID ContextKey = "user_id"
	// KeyUsername is the context key for username
	KeyUsername ContextKey = "username"
	// KeyUserRole is the context key for user role
	KeyUserRole ContextKey = "user_role"
	// KeyUserEmail is the context key for user email
	KeyUserEmail ContextKey = "user_email"
)

// AuthMiddleware provides JWT authentication
type AuthMiddleware struct {
	jwtService *auth.JWTService
	config     *config.AuthConfig
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *auth.JWTService, config *config.AuthConfig, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		config:     config,
		logger:     logger,
	}
}

// Authenticate middleware for JWT authentication
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication if disabled in config
		if !am.config.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, r, "Authorization header is required")
			return
		}

		// Extract the token from the Authorization header
		tokenString := auth.ExtractTokenFromBearer(authHeader)
		if tokenString == "" {
			response.Unauthorized(w, r, "Invalid token format, expected 'Bearer <token>'")
			return
		}

		// Validate the token
		claims, err := am.jwtService.ValidateToken(tokenString)
		if err != nil {
			// Log the error
			am.logger.Debug("JWT validation failed", zap.Error(err), zap.String("token", tokenString))

			// Map common JWT errors to appropriate responses
			switch err {
			case auth.ErrTokenExpired:
				response.Unauthorized(w, r, "Token has expired")
			case auth.ErrTokenInvalid:
				response.Unauthorized(w, r, "Invalid token")
			case auth.ErrInvalidIssuer:
				response.Unauthorized(w, r, "Invalid token issuer")
			default:
				response.Unauthorized(w, r, "Authentication failed")
			}
			return
		}

		// Convert subject to userID
		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			am.logger.Error("Failed to parse subject as userID", zap.Error(err), zap.String("subject", claims.Subject))
			response.Unauthorized(w, r, "Invalid token subject")
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), KeyUserID, userID)
		ctx = context.WithValue(ctx, KeyUsername, claims.Username)
		ctx = context.WithValue(ctx, KeyUserRole, claims.Role)
		ctx = context.WithValue(ctx, KeyUserEmail, claims.Email)

		// Pass the request with user information in context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware for role-based access control
func (am *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip role check if auth is disabled
			if !am.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Get user role from context
			role, ok := r.Context().Value(KeyUserRole).(string)
			if !ok {
				response.Forbidden(w, r, "User role not found in context")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, requiredRole := range roles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Forbidden(w, r, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
