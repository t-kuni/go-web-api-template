package di

import (
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/interface/handler"
	"github.com/t-kuni/go-web-api-template/router"
	"github.com/t-kuni/go-web-api-template/server"
)

func NewContainer() *do.Injector {
	injector := do.New()

	// Server
	do.Provide(injector, server.NewServer)

	// Router
	do.Provide(injector, router.NewRouter)

	// Handler
	do.Provide(injector, handler.NewHelloHandler)
	do.Provide(injector, handler.NewPostUserHandler)

	// Service
	do.Provide(injector, service.NewExampleService)

	// Infrastructure
	do.Provide(injector, db.NewConnector)
	do.Provide(injector, api.NewBinanceApi)

	return injector
}
