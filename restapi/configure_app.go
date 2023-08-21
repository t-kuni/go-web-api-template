// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/di"
	useCaseTodos "github.com/t-kuni/go-web-api-template/domain/usecases/todos"
	"github.com/t-kuni/go-web-api-template/logger"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/t-kuni/go-web-api-template/restapi/operations"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
)

//go:generate swagger generate server --target ../../go-web-api-template --name App --spec ../swagger.yml --model-package restapi/models --principal interface{}

func configureFlags(api *operations.AppAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AppAPI) http.Handler {
	app := di.NewApp()

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

	if api.TodosAddOneHandler == nil {
		api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.AddOne has not yet been implemented")
		})
	}
	if api.TodosDestroyOneHandler == nil {
		api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.DestroyOne has not yet been implemented")
		})
	}
	api.TodosFindTodosHandler = todos.FindTodosHandlerFunc(do.MustInvoke[*useCaseTodos.Find](app).Main)
	if api.TodosUpdateOneHandler == nil {
		api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams) middleware.Responder {
			return middleware.NotImplemented("operation todos.UpdateOne has not yet been implemented")
		})
	}

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
	godotenv.Load("..")

	if err := logger.SetupLogger(); err != nil {
		log.Fatalf("Logger initialization failed: %+v", err)
		os.Exit(1)
	}
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
