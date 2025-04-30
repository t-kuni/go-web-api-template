// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"github.com/joho/godotenv"
	useCaseCompanies "github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	middleware2 "github.com/t-kuni/go-web-api-template/middleware"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
	"github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"go.uber.org/fx"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
)

//go:generate go tool swagger generate server --target ../ --name App --spec ../swagger.yml --model-package restapi/models --principal interface{}

type middleware func(http.Handler) http.Handler

var app *fx.App
var middlewares struct {
	recoverHandler middleware
	accessLog      middleware
}

func configureFlags(api *operations.AppAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AppAPI) http.Handler {
	godotenv.Load()

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	ctx := context.Background()
	app = di.NewApp(fx.Invoke(func(
		recoverHandler *middleware2.Recover,
		accessLog *middleware2.AccessLog,
		logger system.ILogger,
		customServeError func(http.ResponseWriter, *http.Request, error),

		listTodos *useCaseCompanies.ListTodos,
		getCompanies *useCaseCompanies.GetCompanies,
		getCompaniesUsers *useCaseCompanies.GetCompaniesUsers,
		getUsers *useCaseCompanies.GetUsers,
		postUser *useCaseCompanies.PostUser,
	) {
		api.ServeError = customServeError
		middlewares.recoverHandler = recoverHandler.Recover
		middlewares.accessLog = accessLog.AccessLog

		api.TodosGetTodosHandler = todos.GetTodosHandlerFunc(listTodos.Main)
		api.CompaniesGetCompaniesHandler = companies.GetCompaniesHandlerFunc(getCompanies.Main)
		api.CompaniesGetCompaniesUsersHandler = companies.GetCompaniesUsersHandlerFunc(getCompaniesUsers.Main)
		api.UserGetUsersHandler = user.GetUsersHandlerFunc(getUsers.Main)
		api.UserPostUsersHandler = user.PostUsersHandlerFunc(postUser.Main)
	}))
	err := app.Start(ctx)
	if err != nil {
		log.Fatalf("App initialization failed: %+v", err)
		os.Exit(1)
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		app.Stop(ctx)
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return middlewares.recoverHandler(middlewares.accessLog(handler))
}
