package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Define a new validator struct which contains a map of validation error messages
type Validator struct {
	Field_of_Errors map[string]string
}

// Valid() return true if the Field_of_Errors map doesn´t contain any entries
func (v *Validator) Valid() bool {
	return len(v.Field_of_Errors) == 0
}

// AddFieldError() adds an error message to the Field_of_Errors map (so long as no entry exists for the given key)
func (v *Validator) AddFieldError(key, msg string) {
	//Note: We need to initialize the map first, if it isn´t already initialized.
	if v.Field_of_Errors == nil {
		v.Field_of_Errors = make(map[string]string)
	}

	if _, exists := v.Field_of_Errors[key]; !exists {
		v.Field_of_Errors[key] = msg
	}

}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not ´ok´
func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

// NotBlank() returns true if a value is not empty string.
func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""

}

// MaxChars() returns true if a value contains no more than n characters
func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

// PermittedValue() returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
