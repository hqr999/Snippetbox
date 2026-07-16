package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile() function to parse a regular expression pattern
// for sanity checking the format of an email address. This returns a pointer to
// a 'compiled' regexp.Regexp type, or panics in the event of an error. Parsing
// this pattern once at startup and storing the compiled *regexp.Regexp in a
// variable is more performant than re-parsing the pattern each time we need it.
var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

// MinChars() returns true if a value contains n characters or more 
func MinChars(val string, n int) bool {
	return utf8.RuneCountInString(val) >= n
}

// MaxBytes() returns true if a value contains n bytes or less 
func MaxBytes(val string,n int) bool  {
	return len(val) <= n 
}

// Matches() returns true if a value matches a provided compiled regular 
// expression pattern.
func Matches(val string, reg_ex *regexp.Regexp) bool  {
	return reg_ex.MatchString(val)
}
