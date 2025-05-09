package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/linkeunid/go-api/pkg/response"
)

// ValidateRequest validates a struct against validation tags
func ValidateRequest(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	// Get the type of the data
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	// If it's a pointer, get the element
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	// Ensure it's a struct
	if t.Kind() != reflect.Struct {
		response.BadRequest(w, r, "Invalid request data", fmt.Errorf("expected struct, got %v", t.Kind()))
		return false
	}

	// Simple validation errors
	validationErrors := make(map[string]string)

	// Iterate over fields and validate
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Get validation tag
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Split validation rules
		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			// Skip if the rule is empty
			if rule == "" {
				continue
			}

			// Check if the field is required
			if rule == "required" && isZeroValue(fieldValue) {
				validationErrors[field.Name] = "Field is required"
				continue
			}

			// Add more validation rules as needed
			// This is a simplified validation - in a real app,
			// consider using a library like go-playground/validator
		}
	}

	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		response.BadRequest(w, r, "Validation failed", fmt.Errorf("%v", validationErrors))
		return false
	}

	return true
}

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	return !v.IsValid() || v.Interface() == reflect.Zero(v.Type()).Interface()
}
