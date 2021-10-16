package handler

import "github.com/google/wire"

var Providers = wire.NewSet(
	ProvideHello,
)
