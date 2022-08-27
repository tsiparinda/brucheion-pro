package auth

import (
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/sessions"
)

type AuthenticationHandler struct {
	identity.User
	identity.SignInManager
	identity.UserStore
	sessions.Session
	handling.URLGenerator
}

const SIGNIN_MSG_KEY string = "signin_message"

func (handler AuthenticationHandler) GetSignIn() actionresults.ActionResult {
	message := handler.Session.GetValueDefault(SIGNIN_MSG_KEY, "").(string)
	return actionresults.NewTemplateAction("signin.html", message)
}

type Credentials struct {
	Username string
	Password string
}

func (handler AuthenticationHandler) PostSignIn(creds Credentials) actionresults.ActionResult {
	if creds.Password == "mysecret" {
		user, ok := handler.UserStore.GetUserByName(creds.Username)
		if ok {
			if user.InRole("ToolsUsers") {
				handler.Session.SetValue(SIGNIN_MSG_KEY, "")
				handler.SignInManager.SignIn(user)
				return actionresults.NewRedirectAction("/tools/section/")
			}
			if user.InRole("Administrators") {
				handler.Session.SetValue(SIGNIN_MSG_KEY, "")
				handler.SignInManager.SignIn(user)
				return actionresults.NewRedirectAction("/admin/section/")
			}

		}
	}
	handler.Session.SetValue(SIGNIN_MSG_KEY, "Access Denied")
	return actionresults.NewRedirectAction(MustGenerateUrl(handler.URLGenerator,
		AuthenticationHandler.GetSignIn))
}

func (handler AuthenticationHandler) PostSignOut(creds Credentials) actionresults.ActionResult {
	handler.SignInManager.SignOut(handler.User)
	return actionresults.NewRedirectAction("/")
}

func MustGenerateUrl(gen handling.URLGenerator, target interface{},
	data ...interface{}) string {
	url, err := gen.GenerateUrl(target, data...)
	if err != nil {
		panic(err)
	}
	return url
}
