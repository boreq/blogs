package auth

import (
	"github.com/boreq/blogs/auth"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/forms"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views/errors"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
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

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	form, usernameField, passwordField := makeRegisterForm()
	if r.Method == "POST" && form.Validate(r.FormValue) {
		err := auth.CreateUser(usernameField.GetValue(), passwordField.GetValue())
		if err == nil {
			auth.LoginUser(usernameField.GetValue(), passwordField.GetValue(), w)
			http.Redirect(w, r, "/", 302)
			return
		} else {
			if err == auth.UsernameTakenError {
				usernameField.AddError("Username is already taken")
			} else {
				errors.InternalServerErrorWithStack(w, r, err)
				return
			}
		}
	}

	// Render
	var data = templates.GetDefaultData(r)
	data["form"] = form
	if err := templates.RenderTemplateSafe(w, "auth/register.tmpl", data); err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	form, usernameField, passwordField := makeLoginForm()
	if r.Method == "POST" && form.Validate(r.FormValue) {
		err := auth.LoginUser(usernameField.GetValue(), passwordField.GetValue(), w)
		if err != nil {
			if err == auth.InvalidUsernameOrPasswordError {
				form.AddError("Invalid username or password")
			} else {
				errors.InternalServerErrorWithStack(w, r, err)
				return
			}
		} else {
			http.Redirect(w, r, "/", 302)
			return
		}
	}

	// Render
	var data = templates.GetDefaultData(r)
	data["form"] = form
	if err := templates.RenderTemplateSafe(w, "auth/login.tmpl", data); err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	auth.LogoutUser(w, r)
	http.Redirect(w, r, "/", 302)
}

func settingsSessions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userId := context.Get(r).User.GetUser().ID

	currentSessionKey := ""
	if sessionCookie, err := r.Cookie(auth.SessionKeyCookieName); err == nil {
		currentSessionKey = sessionCookie.Value
	}

	var sessions []database.UserSession
	if err := database.DB.Select(&sessions,
		`SELECT *
		FROM user_session
		WHERE user_id=$1
		ORDER BY last_seen DESC`, userId); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["sessions"] = sessions
	data["currentSessionKey"] = currentSessionKey
	if err := templates.RenderTemplateSafe(w, "auth/settings_sessions.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}
