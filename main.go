package main

import (
	"brucheion/admin"
	"brucheion/models"
	"brucheion/models/repo"
	"net/http"
	"sync"

	//"github.com/vedicsociety/platform/http"
	//"github.com/vedicsociety/platform/http/handling"
	//"github.com/vedicsociety/platform/http/handling"

	"brucheion/admin/auth"
	"brucheion/store"

	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/authorization"
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
			authorization.NewRoleCondition("Administrator"),
			admin.AdminHandler{},
			admin.ProductsHandler{},
			admin.CategoriesHandler{},
			admin.OrdersHandler{},
			admin.DatabaseHandler{},
			admin.SignOutHandler{},
		).AddFallback("/admin/section/", "^/admin[/]?$"),

		handling.NewRouter(
			handling.HandlerEntry{"", store.ProductHandler{}},
			handling.HandlerEntry{"", store.CategoryHandler{}},
			//handling.HandlerEntry{"", store.CartHandler{}},
			//handling.HandlerEntry{"", store.OrderHandler{}},
			// handling.HandlerEntry{ "admin", admin.AdminHandler{}},
			// handling.HandlerEntry{ "admin", admin.ProductsHandler{}},
			// handling.HandlerEntry{ "admin", admin.CategoriesHandler{}},
			// handling.HandlerEntry{ "admin", admin.OrdersHandler{}},
			// handling.HandlerEntry{ "admin", admin.DatabaseHandler{}},
			handling.HandlerEntry{"", admin.AuthenticationHandler{}},
			//handling.HandlerEntry{"api", store.RestHandler{}},
			// ).AddMethodAlias("/", store.ProductHandler.GetProducts, 0, 1).
			//     AddMethodAlias("/products[/]?[A-z0-9]*?",
			//         store.ProductHandler.GetProducts, 0, 1),    )
		).AddMethodAlias("/", admin.AdminHandler.GetAdmin))
}

func main() {

	registerServices()
	services.Call(func(repo models.Repository) { repo.LoadMigrations() })
	//repo.LoadMigrations( config.Configuration,  logging.Logger)

	results, err := services.Call(http.Serve, createPipeline())
	if err == nil {
		(results[0].(*sync.WaitGroup)).Wait()
	} else {
		panic(err)
	}

}
