package entities

import (
	"github.com/google/wire"
)

//nolint:gochecknoglobals // Wire dependency injection container
var Container = wire.NewSet(
	NewSettings,
	NewCLI,
)
