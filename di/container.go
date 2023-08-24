package di

import (
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/application/handler/companies"
	"github.com/t-kuni/go-web-api-template/application/handler/todos"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/middleware"
	"github.com/t-kuni/go-web-api-template/validator"
)

func NewApp() *do.Injector {
	injector := do.New()

	// Validator
	do.Provide(injector, validator.NewCustomValidator)

	// Middleware
	do.Provide(injector, middleware.NewRecover)
	do.Provide(injector, middleware.NewAccessLog)

	// Handler
	do.Provide(injector, handler.NewHelloHandler)
	do.Provide(injector, handler.NewPostUserHandler)

	// Service
	do.Provide(injector, service.NewExampleService)

	// Infrastructure
	do.Provide(injector, db.NewConnector)
	do.Provide(injector, api.NewBinanceApi)

	// UseCase
	do.Provide(injector, todos.NewFind)
	do.Provide(injector, companies.NewGetCompanies)

	return injector
}
