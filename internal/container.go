package internal

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	commands.Container,
	controllers.Container,
	repositories.Container,
	NewAppCLI,
)
