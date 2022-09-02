package tools

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
)

type IngestCEXHandler struct {
	models.Repository
	handling.URLGenerator
}

func (handler IngestCEXHandler) GetData() actionresults.ActionResult {
	return actionresults.NewTemplateAction("tools_ingestcex.html", struct {
		Title string
	}{

		Title: "IngestCEX title",
	})
}
