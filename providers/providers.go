package providers

import (
	"github.com/google/wire"
	"github.com/t-kuni/go-cli-app-skeleton/domain/service"
	"github.com/t-kuni/go-cli-app-skeleton/infrastructure/api"
	"github.com/t-kuni/go-cli-app-skeleton/interface/handler"
)

var Providers = wire.NewSet(
	handler.Providers,
	service.Providers,
	api.Providers,
)
