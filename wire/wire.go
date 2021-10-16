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

type App struct {
	handler.HelloHandler
}

var Providers = wire.NewSet(
	handler.Providers,
	service.Providers,
	api.Providers,
)

func InitializeApp() App {
	wire.Build(
		wire.Struct(new(App), "*"),
		Providers,

		// Bindings
		wire.Bind(new(api2.BinanceApiInterface), new(api.BinanceApi)),
	)
	return App{}
}
