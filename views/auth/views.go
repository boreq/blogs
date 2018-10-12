package auth

import (
	"github.com/boreq/blogs/auth"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/forms"
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/http/context"
	//"github.com/boreq/blogs/http/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func New(prefix string) *Auth {
	rv := &Auth{
		Prefix: prefix,
	}
	return rv
}

type Auth struct {
	Prefix string
}

func (a *Auth) Register(router *httprouter.Router) {
	router.POST(a.Prefix+"/register", api.Wrap(a.register))
	router.POST(a.Prefix+"/login", api.Wrap(a.login))
	router.POST(a.Prefix+"/logout", api.Wrap(a.logout))
	router.GET(a.Prefix+"/check-login", api.Wrap(a.checkLogin))
}

type formError struct {
	Errors      []string            `json:"errors,omitempty"`
	FieldErrors map[string][]string `json:"fieldErrors,omitempty"`
}

func toFormError(form forms.Form) formError {
	var formError formError
	formError.FieldErrors = make(map[string][]string)
	formError.Errors = form.Errors
	for _, field := range form.Fields {
		if len(field.Errors()) > 0 {
			formError.FieldErrors[field.GetName()] = field.Errors()
		}
	}
	return formError
}

type userWithSessionKey struct {
	User       *database.User `json:"user"`
	SessionKey string         `json:"token"`
}

func (a *Auth) register(r *http.Request, _ httprouter.Params) (api.Response, api.Error) {
	getFormValue, err := forms.GetJsonFormValue(r)
	if err != nil {
		return nil, api.BadRequestError
	}
	form, usernameField, passwordField := makeRegisterForm()
	if form.Validate(getFormValue) {
		err := auth.CreateUser(usernameField.GetValue(), passwordField.GetValue())
		if err == nil {
			user, sessionKey, err := auth.LoginUser(usernameField.GetValue(), passwordField.GetValue())
			response := userWithSessionKey{
				User:       user,
				SessionKey: sessionKey,
			}
			if err == nil {
				return api.NewResponseOk(response), nil
			}
			return nil, api.NewError(http.StatusInternalServerError, "Registration successful but login failed.")
		} else {
			if err == auth.UsernameTakenError {
				usernameField.AddError("Username is already taken")
			} else {
				return nil, api.InternalServerError
			}
		}
	}
	errors := toFormError(form)
	return api.NewResponse(http.StatusBadRequest, errors), nil
}

func (a *Auth) login(r *http.Request, _ httprouter.Params) (api.Response, api.Error) {
	getFormValue, err := forms.GetJsonFormValue(r)
	if err != nil {
		return nil, api.BadRequestError
	}
	form, usernameField, passwordField := makeLoginForm()
	if form.Validate(getFormValue) {
		user, sessionKey, err := auth.LoginUser(usernameField.GetValue(), passwordField.GetValue())
		if err != nil {
			if err == auth.InvalidUsernameOrPasswordError {
				form.AddError("Invalid username or password")
			} else {
				// TODO log
				return nil, api.InternalServerError
			}
		} else {
			response := userWithSessionKey{
				User:       user,
				SessionKey: sessionKey,
			}
			return api.NewResponseOk(response), nil
		}
	}
	errors := toFormError(form)
	return api.NewResponse(http.StatusBadRequest, errors), nil
}

func (a *Auth) checkLogin(r *http.Request, _ httprouter.Params) (api.Response, api.Error) {
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user := ctx.User.GetUser()
		return api.NewResponseOk(user), nil
	}
	return nil, api.UnauthorizedError
}

func (a *Auth) logout(r *http.Request, _ httprouter.Params) (api.Response, api.Error) {
	err := auth.LogoutUser(r)
	if err != nil {
		// TODO log
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}

func settingsSessions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//	userId := context.Get(r).User.GetUser().ID
	//
	////	currentSessionKey := ""
	//	if sessionCookie, err := r.Cookie(auth.SessionKeyCookieName); err == nil {
	//		currentSessionKey = sessionCookie.Value
	//	}
	//
	//	var sessions []database.UserSession
	//	if err := database.DB.Select(&sessions,
	//		`SELECT *
	//		FROM user_session
	//		WHERE user_id=$1
	//		ORDER BY last_seen DESC`, userId); err != nil {
	//		//		verrors.InternalServerErrorWithStack(w, r, err)
	//		return
	//	}

	//var data = templates.GetDefaultData(r)
	//data["sessions"] = sessions
	//data["currentSessionKey"] = currentSessionKey
	//if err := templates.RenderTemplateSafe(w, "auth/settings_sessions.tmpl", data); err != nil {
	//	verrors.InternalServerErrorWithStack(w, r, err)
	//	return
	//}
}
