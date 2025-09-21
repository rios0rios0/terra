package repositories

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/repositories"
)

//nolint:gochecknoglobals // Wire dependency injection container
var Container = wire.NewSet(
	NewStdShellRepository,
	NewInteractiveShellRepository,
	wire.Bind(new(repositories.ShellRepository), new(*StdShellRepository)),
)
