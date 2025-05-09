package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/linkeunid/go-api/pkg/response"
	"github.com/linkeunid/go-api/pkg/validator"
)

// ValidationMiddleware is a middleware that validates the request body against a model
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate on POST, PUT, PATCH
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			// Check if content type is application/json
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				handleValidationError(w, r, []validator.ValidationError{
					{
						Field: "Content-Type",
						Tag:   "required",
						Error: "Content-Type must be application/json",
					},
				})
				return
			}

			// Store the validation data in the request context
			r = r.WithContext(r.Context())

			// Call the next handler, which will validate the model specifically
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateModel validates a model from the request body
func ValidateModel(model interface{}, r *http.Request) []validator.ValidationError {
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		return []validator.ValidationError{
			{
				Field: "body",
				Tag:   "json",
				Error: "Invalid JSON format: " + err.Error(),
			},
		}
	}

	// If the model has a custom Validate method, use it
	if v, ok := model.(interface {
		Validate() []validator.ValidationError
	}); ok {
		return v.Validate()
	}

	// Otherwise use the standard validator
	return validator.Validate(model)
}

// handleValidationError responds with validation errors
func handleValidationError(w http.ResponseWriter, r *http.Request, errors []validator.ValidationError) {
	response.ValidationError(w, r, errors)
}

// HandleValidateRequest validates a model and returns appropriate response
func HandleValidateRequest(w http.ResponseWriter, r *http.Request, model interface{}) bool {
	errors := ValidateModel(model, r)
	if len(errors) > 0 {
		handleValidationError(w, r, errors)
		return false
	}
	return true
}
