package tools

import (
	"brucheion/admin/auth"
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/services"
)

type PassageOverviewHandler struct {
	models.Repository
	handling.URLGenerator
}

func (handler PassageOverviewHandler) GetData() actionresults.ActionResult {
	return actionresults.NewTemplateAction("tools_passageoverview.html", struct {
		InitUrl, SeedUrl string
	}{
		InitUrl: auth.MustGenerateUrl(handler.URLGenerator,
			PassageOverviewHandler.PostDatabaseInit),
		SeedUrl: auth.MustGenerateUrl(handler.URLGenerator,
			PassageOverviewHandler.PostDatabaseSeed),
	})
}

func (handler PassageOverviewHandler) PostDatabaseInit() actionresults.ActionResult {
	//handler.Repository.Init()
	services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	return actionresults.NewRedirectAction(auth.MustGenerateUrl(handler.URLGenerator,
		ToolsHandler.GetSection, "PassageOverview"))
}

func (handler PassageOverviewHandler) PostDatabaseSeed() actionresults.ActionResult {
	// handler.Repository.Seed()
	return actionresults.NewRedirectAction(auth.MustGenerateUrl(handler.URLGenerator,
		ToolsHandler.GetSection, "PassageOverview"))
}
