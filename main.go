package main

import (
	"brucheion/api"
	"brucheion/models/repo"
	"sync"

	"github.com/vedicsociety/platform/http"

	"brucheion/admin"
	"brucheion/auth"
	"brucheion/tools"

	//"brucheion/store"

	"github.com/vedicsociety/platform/authorization"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/pipeline"
	"github.com/vedicsociety/platform/pipeline/basic"
	"github.com/vedicsociety/platform/services"
	"github.com/vedicsociety/platform/sessions"
)

func registerServices() {
	services.RegisterDefaultServices()
	//repo.RegisterMemoryRepoService()
	repo.RegisterSqlRepositoryService()
	sessions.RegisterSessionService()
	//cart.RegisterCartService()
	authorization.RegisterDefaultSignInService()
	authorization.RegisterDefaultUserService()
	auth.RegisterUserStoreService()
}

func createPipeline() pipeline.RequestPipeline {
	return pipeline.CreatePipeline(
		&basic.ServicesComponent{},
		&basic.LoggingComponent{},
		&basic.ErrorComponent{},
		&basic.StaticFileComponent{},
		&sessions.SessionComponent{},

		authorization.NewAuthComponent(
			"admin",
			authorization.NewRoleCondition("Administrators"),
			admin.AdminHandler{},
			admin.DatabaseHandler{},
			auth.SignOutHandler{},
		).AddFallback("/admin/section/", "^/admin[/]?$"),

		authorization.NewAuthComponent(
			"tools",
			authorization.NewRoleCondition("ToolsUsers"),
			tools.ToolsHandler{},
			tools.PassageOverviewHandler{},
			tools.IngestCEXHandler{},
			auth.SignOutHandler{},
		).AddFallback("/tools/section/", "^/tools[/]?$"),

		authorization.NewAuthComponent(
			"api/v1",
			authorization.NewRoleCondition("ToolsUsers"),
			api.PassageHandler{},
			api.UserHandler{},
			api.CEXuploadHandler{},
		).AddFallback("/api/v1/", "^/api[/]?$"),

		handling.NewRouter(
			handling.HandlerEntry{"", auth.AuthenticationHandler{}},
			// handling.HandlerEntry{"api/v1", api.PassageHandler{}},
			// handling.HandlerEntry{"api/v1", api.UserHandler{}},
			// handling.HandlerEntry{"api/v1", api.CEXuploadHandler{}},
		).AddMethodAlias("/", tools.ToolsHandler.GetTools))
}

func main() {

	registerServices()
	//services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	//repo.LoadMigrations( config.Configuration,  logging.Logger)

	results, err := services.Call(http.Serve, createPipeline())
	if err == nil {
		(results[0].(*sync.WaitGroup)).Wait()
	} else {
		panic(err)
	}

}
