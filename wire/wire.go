//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/t-kuni/go-cli-app-skeleton/domain/service"
	"github.com/t-kuni/go-cli-app-skeleton/infrastructure/api"
)

type App struct {
	ExampleService service.ExampleService
}

func InitializeApp() App {
	wire.Build(
		wire.Struct(new(App), "*"),
		api.Providers,
		service.Providers,
	)
	return App{}
}
