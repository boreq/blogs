package auth

import (
	"github.com/boreq/blogs/forms"
	"github.com/boreq/blogs/templates"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func makeRegisterForm() forms.Form {
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
	return form
}

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	form := makeRegisterForm()
	if r.Method == "POST" {
		if form.Validate(r.FormValue) {
		}
	}

	// Render
	var data = make(map[string]interface{})
	data["form"] = form
	return templates.RenderTemplate(w, "auth/register.tmpl", data)
}
