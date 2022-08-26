package auth

import (
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
)

type SignOutHandler struct {
	identity.User
	handling.URLGenerator
}

func (handler SignOutHandler) GetUserWidget() actionresults.ActionResult {
	return actionresults.NewTemplateAction("user_widget.html", struct {
		identity.User
		SignoutUrl string
	}{
		handler.User,
		MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.PostSignOut),
	})
}
