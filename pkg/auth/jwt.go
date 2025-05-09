package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/linkeunid/go-api/pkg/config"
)

// Common JWT errors
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenNotProvided = errors.New("token not provided")
	ErrInvalidIssuer    = errors.New("token has invalid issuer")
	ErrEmptySecret      = errors.New("JWT secret is empty")
)

// Claims represents the JWT claims with standard and custom claims
type Claims struct {
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// JWTService provides JWT operations
type JWTService struct {
	config *config.AuthConfig
}

// NewJWTService creates a new JWT service with the provided configuration
func NewJWTService(cfg *config.AuthConfig) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

// GenerateToken generates a new JWT token with the provided claims
func (s *JWTService) GenerateToken(userID uint64, username, role, email string) (string, error) {
	if s.config.JWTSecret == "" {
		return "", ErrEmptySecret
	}

	// Set the token expiration time
	expirationTime := time.Now().Add(s.config.JWTExpiration)

	// Create the JWT claims with standard and custom fields
	claims := &Claims{
		Username: username,
		Role:     role,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set standard claims
			Issuer:    "linkeun-go-api",
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create and sign the token with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the provided token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrTokenNotProvided
	}

	if s.config.JWTSecret == "" {
		return nil, ErrEmptySecret
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	// Handle parsing errors
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if issuer is allowed (if configured)
		if len(s.config.AllowedIssuers) > 0 {
			issuerAllowed := false
			for _, allowedIssuer := range s.config.AllowedIssuers {
				if claims.Issuer == allowedIssuer {
					issuerAllowed = true
					break
				}
			}
			if !issuerAllowed {
				return nil, ErrInvalidIssuer
			}
		}
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ExtractTokenFromBearer extracts token from "Bearer <token>" format
func ExtractTokenFromBearer(authHeader string) string {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return ""
}
