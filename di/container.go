package di

import (
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/application/handler/companies"
	"github.com/t-kuni/go-web-api-template/application/handler/todos"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/infrastructure/system"
	"github.com/t-kuni/go-web-api-template/middleware"
	"github.com/t-kuni/go-web-api-template/validator"
	"go.uber.org/fx"
)

func NewApp(opts ...fx.Option) *fx.App {
	mergedOpts := []fx.Option{
		//fx.WithLogger(func(log *logger.Logger) fxevent.Logger {
		//	return log
		//}),
		fx.Provide(

			// Validator
			validator.NewCustomValidator,

			// Middleware
			middleware.NewRecover,
			middleware.NewAccessLog,

			// Handler
			handler.NewHelloHandler,

			// Service
			service.NewExampleService,

			// Infrastructure
			db.NewConnector,
			api.NewBinanceApi,
			system.NewTimer,

			// UseCase
			todos.NewListTodos,
			companies.NewGetCompanies,
		),
	}
	mergedOpts = append(mergedOpts, opts...)

	return fx.New(mergedOpts...)
}
