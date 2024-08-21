package commands

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewDeleteCacheCommand,
	wire.Bind(new(interfaces.DeleteCache), new(*DeleteCacheCommand)),
	NewFormatFilesCommand,
	wire.Bind(new(interfaces.FormatFiles), new(*FormatFilesCommand)),
	NewInstallDependenciesCommand,
	wire.Bind(new(interfaces.InstallDependencies), new(*InstallDependenciesCommand)),
	NewRunAdditionalBeforeCommand,
	wire.Bind(new(interfaces.RunAdditionalBefore), new(*RunAdditionalBeforeCommand)),
	NewRunFromRootCommand,
	wire.Bind(new(interfaces.RunFromRoot), new(*RunFromRootCommand)),
)
