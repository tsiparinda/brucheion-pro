package admin

import (
	"brucheion/auth"
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/services"
)

type DatabaseHandler struct {
	models.Repository
	handling.URLGenerator
}

func (handler DatabaseHandler) GetData() actionresults.ActionResult {
	return actionresults.NewTemplateAction("admin_database.html", struct {
		InitUrl, SeedUrl string
	}{
		InitUrl: auth.MustGenerateUrl(handler.URLGenerator,
			DatabaseHandler.PostDatabaseInit),
		SeedUrl: auth.MustGenerateUrl(handler.URLGenerator,
			DatabaseHandler.PostDatabaseSeed),
	})
}

func (handler DatabaseHandler) PostDatabaseInit() actionresults.ActionResult {
	//handler.Repository.Init()
	services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	return actionresults.NewRedirectAction(auth.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Database"))
}

func (handler DatabaseHandler) PostDatabaseSeed() actionresults.ActionResult {
	// handler.Repository.Seed()
	return actionresults.NewRedirectAction(auth.MustGenerateUrl(handler.URLGenerator,
		AdminHandler.GetSection, "Database"))
}
