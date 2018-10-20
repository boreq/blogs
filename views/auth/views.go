package auth

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/forms"
	"github.com/boreq/blogs/http/api"
	"github.com/boreq/blogs/logging"
	"github.com/boreq/blogs/service/auth"
	"github.com/boreq/blogs/service/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var log = logging.New("views/auth")

func New(prefix string, authService *auth.AuthService, contextService *context.ContextService) *Auth {
	rv := &Auth{
		Prefix:         prefix,
		authService:    authService,
		contextService: contextService,
	}
	return rv
}

type Auth struct {
	Prefix         string
	authService    *auth.AuthService
	contextService *context.ContextService
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
		err := a.authService.CreateUser(usernameField.GetValue(), passwordField.GetValue())
		if err == nil {
			user, sessionKey, err := a.authService.LoginUser(usernameField.GetValue(), passwordField.GetValue())
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
				log.Error("register error", "err", err)
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
		user, sessionKey, err := a.authService.LoginUser(usernameField.GetValue(), passwordField.GetValue())
		if err != nil {
			if err == auth.InvalidUsernameOrPasswordError {
				form.AddError("Invalid username or password")
			} else {
				log.Error("login error", "err", err)
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
	ctx := a.contextService.Get(r)
	if ctx.User.IsAuthenticated() {
		user := ctx.User.GetUser()
		return api.NewResponseOk(user), nil
	}
	return nil, api.UnauthorizedError
}

func (a *Auth) logout(r *http.Request, _ httprouter.Params) (api.Response, api.Error) {
	err := a.authService.LogoutUser(r)
	if err != nil {
		log.Error("logout error", "err", err)
		return nil, api.InternalServerError
	}
	return api.NewResponseOk(nil), nil
}
