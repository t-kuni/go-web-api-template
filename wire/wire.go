//go:build wireinject

//go:generate wire

package wire

import (
	"github.com/google/wire"
	api2 "github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	db2 "github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"github.com/t-kuni/go-web-api-skeleton/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/infrastructure/db"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
)

var Providers = wire.NewSet(
	handler.Providers,
	service.Providers,
	api.Providers,
	db.Providers,
)

type App struct {
	DBConnector db2.ConnectorInterface

	HelloHandler handler.HelloHandler
}

func InitializeApp() (App, func(), error) {
	wire.Build(
		wire.Struct(new(App), "*"),
		Providers,

		//
		// Bindings
		//

		// Service
		wire.Bind(new(service.ExampleServiceInterface), new(service.ExampleService)),

		// api
		wire.Bind(new(api2.BinanceApiInterface), new(api.BinanceApi)),

		// db
		wire.Bind(new(db2.ConnectorInterface), new(*db.Connector)),
	)
	return App{}, nil, nil
}
