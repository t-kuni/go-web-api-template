// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	useCaseCompanies "github.com/t-kuni/go-web-api-template/application/handler/companies"
	useCaseTodos "github.com/t-kuni/go-web-api-template/application/handler/todos"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/logger"
	middleware2 "github.com/t-kuni/go-web-api-template/middleware"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
)

//go:generate swagger generate server --target ../ --name App --spec ../swagger.yml --model-package restapi/models --principal interface{}

var app *do.Injector

func configureFlags(api *operations.AppAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AppAPI) http.Handler {
	godotenv.Load()

	if err := logger.SetupLogger(); err != nil {
		log.Fatalf("Logger initialization failed: %+v", err)
		os.Exit(1)
	}

	app = di.NewApp()

	// configure the api here
	api.ServeError = errors.ServeError

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

	api.TodosGetTodosHandler = todos.GetTodosHandlerFunc(do.MustInvoke[*useCaseTodos.ListTodos](app).Main)
	api.CompaniesGetCompaniesHandler = companies.GetCompaniesHandlerFunc(do.MustInvoke[*useCaseCompanies.GetCompanies](app).Main)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		app.Shutdown()
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
	recoverHandler := do.MustInvoke[*middleware2.Recover](app).Recover
	accessLog := do.MustInvoke[*middleware2.AccessLog](app).AccessLog
	return recoverHandler(accessLog(handler))
}
