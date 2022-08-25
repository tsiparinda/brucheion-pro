package admin

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/services"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
)

type DatabaseHandler struct {
	models.Repository
	handling.URLGenerator
}

func (handler DatabaseHandler) GetData() actionresults.ActionResult {
	return actionresults.NewTemplateAction("admin_database.html", struct {
		InitUrl, SeedUrl string
	}{
		InitUrl: mustGenerateUrl(handler.URLGenerator,
			DatabaseHandler.PostDatabaseInit),
		SeedUrl: mustGenerateUrl(handler.URLGenerator,
			DatabaseHandler.PostDatabaseSeed),
	})
}

func (handler DatabaseHandler) PostDatabaseInit() actionresults.ActionResult {
	//handler.Repository.Init()
	services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	return actionresults.NewRedirectAction(mustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Database"))
}

func (handler DatabaseHandler) PostDatabaseSeed() actionresults.ActionResult {
	// handler.Repository.Seed()
	return actionresults.NewRedirectAction(mustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Database"))
}
