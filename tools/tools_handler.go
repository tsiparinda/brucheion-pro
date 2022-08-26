package tools

import (
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
)

var sectionNames = []string{"PassageOverview"}

type ToolsHandler struct {
	handling.URLGenerator
}

type ToolsTemplateContext struct {
	Sections       []string
	ActiveSection  string
	SectionUrlFunc func(string) string
}

func (handler ToolsHandler) GetSection(section string) actionresults.ActionResult {
	return actionresults.NewTemplateAction("tools.html", ToolsTemplateContext{
		Sections:      sectionNames,
		ActiveSection: section,
		SectionUrlFunc: func(sec string) string {
			sectionUrl, _ := handler.GenerateUrl(ToolsHandler.GetSection, sec)
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
