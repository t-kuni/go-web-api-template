package providers

import (
	"github.com/google/wire"
	"github.com/t-kuni/go-cli-app-skeleton/domain/service"
	"github.com/t-kuni/go-cli-app-skeleton/infrastructure/api"
)

var Providers = wire.NewSet(
	api.Providers,
	service.Providers,
)
