//go:build wireinject

//go:generate wire

package wire

import (
	"github.com/google/wire"
	api2 "github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"github.com/t-kuni/go-web-api-skeleton/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
)

var Providers = wire.NewSet(
	handler.Providers,
	service.Providers,
	api.Providers,
)

type App struct {
	handler.HelloHandler
}

func InitializeApp() App {
	wire.Build(
		wire.Struct(new(App), "*"),
		Providers,

		// Binding
		wire.Bind(new(service.ExampleServiceInterface), new(service.ExampleService)),
		wire.Bind(new(api2.BinanceApiInterface), new(api.BinanceApi)),
	)
	return App{}
}
