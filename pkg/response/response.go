package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/linkeunid/go-api/pkg/pagination"
)

// APIResponse represents a standardized API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// sendResponse sends a JSON response with the provided status code and data
func sendResponse(w http.ResponseWriter, r *http.Request, statusCode int, resp APIResponse) {
	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Add timestamp if not set
	if resp.Timestamp.IsZero() {
		resp.Timestamp = time.Now()
	}

	// Encode response to JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// If encoding fails, send a plain text error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Success sends a successful response with data
func Success(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	sendResponse(w, r, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a response for successful resource creation
func Created(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	sendResponse(w, r, http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a response with no content
func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(w http.ResponseWriter, r *http.Request, items interface{}, params pagination.Params, message string) {
	paginatedData := pagination.PagedData{
		Items: items,
		Meta:  params,
	}

	sendResponse(w, r, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    paginatedData,
	})
}

// BadRequest sends a bad request error response
func BadRequest(w http.ResponseWriter, r *http.Request, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	sendResponse(w, r, http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// NotFound sends a not found error response
func NotFound(w http.ResponseWriter, r *http.Request, message string) {
	sendResponse(w, r, http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	sendResponse(w, r, http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: "An internal server error occurred",
		Error:   errorMsg,
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	sendResponse(w, r, http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
	})
}

// Forbidden sends a forbidden error response
func Forbidden(w http.ResponseWriter, r *http.Request, message string) {
	sendResponse(w, r, http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
	})
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, r *http.Request, errors interface{}) {
	sendResponse(w, r, http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "Validation failed",
		Error:   "The provided data is invalid",
		Data:    errors,
	})
}
