package forms

// GetFormValue is used to load parameters into the form. Usually this will
// simply be an http.Request.FormValue.
type GetFormValue func(string) string

// FormValidator can validate multiple form fields or generate errors which
// are not related to any fields. FormValidator can attach errors to a form
// by calling Form.AddError or to a specific field by calling Field.AddError.
// If a validator validates a single field it is better to attach a validator
// directly to that field.
type FormValidator func(f *Form)

type Form struct {
	Fields     []Field
	Validators []FormValidator
	Errors     ValidationErrors
}

// AddField is a shortcut for append.
func (f *Form) AddField(field Field) {
	f.Fields = append(f.Fields, field)
}

// AddValidator is a shortcut for append.
func (f *Form) AddValidator(validator FormValidator) {
	f.Validators = append(f.Validators, validator)
}

// AddError is a shortcut for append.
func (f *Form) AddError(err string) {
	f.Errors = append(f.Errors, err)
}

// Validate initializes all form fields. Calling this method calls
// Field.Validate on all form fields. After calling this method Errors
// will be populated and IsValid() will return the correct value.
func (f *Form) Validate(getFormValue GetFormValue) bool {
	for _, field := range f.Fields {
		field.Validate(getFormValue(field.GetName()))
	}
	for _, validator := range f.Validators {
		validator(f)
	}
	return f.IsValid()
}

// IsValid can be called after calling Validate.
func (f *Form) IsValid() bool {
	for _, field := range f.Fields {
		if !field.IsValid() {
			return false
		}
	}
	return f.Errors == nil || len(f.Errors) == 0
}