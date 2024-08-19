package repositories

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/repositories"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewShellRepository,
	wire.Bind(new(repositories.ShellRepository), new(*ShellRepository)),
)
