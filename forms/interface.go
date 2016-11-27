package forms

import (
	"html/template"
	"strings"
)

type Field interface {
	// Render generates an HTML code which can be displayed in the form.
	Render() template.HTML

	// GetName returns the value of the name attribute.
	GetName() string

	// GetLabel returns the text which should be displayed as label.
	GetLabel() string

	// GetHelpText returns an instruction for the user.
	GetHelpText() string

	// AddError adds an error to this field.
	AddError(err string)

	// Validate returns false if any validators defined on this field
	// return a non-empty list of errors. This method has to be called
	// before calling Errors(), IsValid() and GetValue(). Validate runs the
	// value through the cleaners defined on the field before validating it
	// using the validators. After calling this function GetValue() will
	// return the cleaned value, Errors() will return the list of errors
	// returned by the validators and IsValid() will return the same value
	// that was returned by this function.
	Validate(value string) bool

	// Errors can be called after calling Validate.
	Errors() ValidationErrors

	// IsValid can be called after calling Validate.
	IsValid() bool

	// GetValue can be called after calling Validate.
	GetValue() string
}

// Cleaner accepts a field value and standarizes it.
type Cleaner func(value string) string

// Validator accepts a field value and returns a list of errors.
type Validator func(value string) []string

// ValidationErrors provide methods useful for displaying errors in the
// templates.
type ValidationErrors []string

// AsText returns the errors separated with dots and ending with a dot.
func (f ValidationErrors) AsText() string {
	if len(f) == 0 {
		return ""
	}
	return strings.Join(f, ". ") + "."
}
