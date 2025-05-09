package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
	Error string `json:"error"`
}

// Validator is a global validator instance
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
}

// New creates a new validator instance
func New() *Validator {
	v := &Validator{}
	v.init()
	return v
}

// init initializes the validator with custom validation rules
func (v *Validator) init() {
	v.once.Do(func() {
		// Create a new validator instance
		validate := validator.New()

		// Register function to get json tag name instead of struct field name
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom validation functions
		v.registerCustomValidations(validate)

		// Set up the translator
		english := en.New()
		uni := ut.New(english, english)
		trans, _ := uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(validate, trans)

		// Register custom error messages
		v.registerCustomTranslations(validate, trans)

		v.validate = validate
		v.trans = trans
	})
}

// registerCustomValidations registers custom validation functions
func (v *Validator) registerCustomValidations(validate *validator.Validate) {
	// Example: Register a custom validation for animal names
	validate.RegisterValidation("animalname", func(fl validator.FieldLevel) bool {
		// Animal names should not contain numbers or special characters
		name := fl.Field().String()
		reg := regexp.MustCompile(`^[a-zA-Z\s-]+$`)
		return reg.MatchString(name)
	})

	// Add more custom validations as needed
}

// registerCustomTranslations registers custom error messages for validations
func (v *Validator) registerCustomTranslations(validate *validator.Validate, trans ut.Translator) {
	// Register custom error message for animalname validation
	validate.RegisterTranslation("animalname", trans, func(ut ut.Translator) error {
		return ut.Add("animalname", "{0} must contain only letters, spaces, and hyphens", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("animalname", fe.Field())
		return t
	})

	// Customize the required error message
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// Add more custom translations as needed
}

// ValidateStruct validates a struct and returns a list of validation errors
func (v *Validator) ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	// Initialize the validator if not already done
	v.init()

	// Validate the struct
	err := v.validate.Struct(s)
	if err != nil {
		// Convert validation errors to our custom format
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = fmt.Sprintf("%v", err.Value())
			element.Error = err.Translate(v.trans)
			errors = append(errors, element)
		}
	}

	return errors
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	// Initialize the validator if not already done
	v.init()

	return v.validate.Var(field, tag)
}

// Global instance for convenience
var validate = New()

// Validate is a convenience function that uses the global validator instance
func Validate(s interface{}) []ValidationError {
	return validate.ValidateStruct(s)
}

// ValidateVar is a convenience function that uses the global validator instance
func ValidateVar(field interface{}, tag string) error {
	return validate.ValidateVar(field, tag)
}
