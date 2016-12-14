package auth

import (
	"github.com/boreq/blogs/forms"
)

func makeUsernameAndPasswordFields() (*forms.TextField, *forms.TextField) {
	usernameField := &forms.TextField{
		Name:     "username",
		Label:    "Username:",
		HelpText: "Your username is used to log in and identifies you on the website.",
		Validators: []forms.Validator{
			forms.MaxLength(50),
			forms.MinLength(3),
			forms.Regexp("^[A-Za-z0-9]+$"),
		},
	}
	usernameField.SetAttribute("class", "form-control")

	passwordField := forms.ToPasswordField(&forms.TextField{
		Name:     "password",
		Label:    "Password:",
		HelpText: "Choose your password, the passwords are hashed using bcrypt.",
		Validators: []forms.Validator{
			forms.MinLength(1),
		},
	})
	passwordField.SetAttribute("class", "form-control")

	return usernameField, passwordField
}

func makeRegisterForm() (forms.Form, forms.Field, forms.Field) {
	usernameField, passwordField := makeUsernameAndPasswordFields()

	passwordConfirmField := forms.ToPasswordField(&forms.TextField{
		Name:     "password_confirm",
		Label:    "Confirm password:",
		HelpText: "Confirm your password.",
	})
	passwordConfirmField.SetAttribute("class", "form-control")

	form := forms.Form{}
	form.AddField(usernameField)
	form.AddField(passwordField)
	form.AddField(passwordConfirmField)
	form.AddValidator(forms.FieldsEqual(passwordField, passwordConfirmField))
	return form, usernameField, passwordField
}

func makeLoginForm() (forms.Form, forms.Field, forms.Field) {
	usernameField, passwordField := makeUsernameAndPasswordFields()
	usernameField.HelpText = ""
	passwordField.HelpText = ""
	form := forms.Form{}
	form.AddField(usernameField)
	form.AddField(passwordField)
	return form, usernameField, passwordField
}
