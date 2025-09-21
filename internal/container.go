package internal

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

//nolint:gochecknoglobals // Wire dependency injection container
var Container = wire.NewSet(
	repositories.Container,
	entities.Container,
	commands.Container,
	controllers.Container,
	NewAppInternal,
)
