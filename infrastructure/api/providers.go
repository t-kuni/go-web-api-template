package api

import "github.com/google/wire"

var Providers = wire.NewSet(
	ProvideBinanceApi,
)
