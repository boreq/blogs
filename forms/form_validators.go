package forms

// FieldsEqual attaches an error to the second field if the values of the
// fields are not equal.
func FieldsEqual(a Field, b Field) FormValidator {
	return func(form *Form) {
		if a.GetValue() != b.GetValue() {
			err := "The value of this field is different."
			b.AddError(err)
		}
	}
}
