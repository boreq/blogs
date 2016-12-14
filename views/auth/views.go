package auth

import (
	"github.com/boreq/blogs/auth"
	"github.com/boreq/blogs/database"
	bhttp "github.com/boreq/blogs/http"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/views/errors"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	form, usernameField, passwordField := makeRegisterForm()
	if r.Method == "POST" && form.Validate(r.FormValue) {
		err := auth.CreateUser(usernameField.GetValue(), passwordField.GetValue())
		if err == nil {
			auth.LoginUser(usernameField.GetValue(), passwordField.GetValue(), w)
			bhttp.Redirect(w, r, "/")
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
			bhttp.Redirect(w, r, "/")
			return
		}
	}

	var data = templates.GetDefaultData(r)
	data["form"] = form
	if err := templates.RenderTemplateSafe(w, "auth/login.tmpl", data); err != nil {
		errors.InternalServerErrorWithStack(w, r, err)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	auth.LogoutUser(w, r)
	bhttp.Redirect(w, r, "/")
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
