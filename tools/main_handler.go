package tools

import (
	//"github.com/gorilla/sessions"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/sessions"
)

var sectionNames = []string{"PassageOverview", "IngestCEX"}

type ToolsHandler struct {
	handling.URLGenerator
	sessions.Session
	identity.User
}

type ToolsTemplateContext struct {
	UserName       string
	Sections       []string
	ActiveSection  string
	SectionUrlFunc func(string) string
}

const USER_SESSION_KEY string = "USER" // for test - see below

func (handler ToolsHandler) GetSection(section string) actionresults.ActionResult {
	return actionresults.NewTemplateAction("tools.html", ToolsTemplateContext{
		UserName:      handler.User.GetDisplayName(), //( handler.Session.GetValueDefault(USER_SESSION_KEY, 0).(int)),
		Sections:      sectionNames,
		ActiveSection: section,
		SectionUrlFunc: func(sec string) string {
			sectionUrl, _ := handler.GenerateUrl(ToolsHandler.GetSection, sec)
			// fmt.Println(sectionUrl)
			// services.Call(func(logger logging.Logger) { logger.Debugf("section url:", sectionUrl) })
			return sectionUrl
		},
	})
}
func (handler ToolsHandler) GetTools() actionresults.ActionResult {
	return actionresults.NewTemplateAction("tools.html", ToolsTemplateContext{
		Sections:      sectionNames,
		ActiveSection: "",
		SectionUrlFunc: func(sec string) string {
			sectionUrl, _ := handler.GenerateUrl(ToolsHandler.GetTools)
			return sectionUrl
		},
	})
}
