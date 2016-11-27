package forms

import (
	"fmt"
	"regexp"
)

// MaxLength returns an error if the value is too long.
func MaxLength(maxLength int) Validator {
	return func(v string) []string {
		var errors []string
		if len(v) > maxLength {
			err := fmt.Sprintf("Max length of this field is %d", maxLength)
			errors = append(errors, err)
		}
		return errors
	}
}

// MinLength returns an error if the value is too short.
func MinLength(minLength int) Validator {
	return func(v string) []string {
		var errors []string
		if len(v) < minLength {
			err := fmt.Sprintf("Min length of this field is %d", minLength)
			errors = append(errors, err)
		}
		return errors
	}
}

// Regexp returns an error if the value doesn't match the regex.
func Regexp(regExp string) Validator {
	r := regexp.MustCompile(regExp)
	return func(v string) []string {
		var errors []string
		if !r.MatchString(v) {
			err := fmt.Sprintf("Allowed characters %s", regExp)
			errors = append(errors, err)
		}
		return errors
	}
}
