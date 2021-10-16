package db

import "github.com/google/wire"

var Providers = wire.NewSet(
	ProvideConnector,
)
