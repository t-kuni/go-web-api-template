package providers

import (
	"github.com/google/wire"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"github.com/t-kuni/go-web-api-skeleton/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
)

var Providers = wire.NewSet(
	handler.Providers,
	service.Providers,
	api.Providers,
)
