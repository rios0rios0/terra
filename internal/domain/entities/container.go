package entities

import (
	"github.com/google/wire"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewSettings,
	NewCLI,
)
