package commands

import "github.com/google/wire"

var Container = wire.NewSet(
	NewDeleteCacheCommand,
	wire.Bind(new(DeleteCache), new(*DeleteCacheCommand)),
	NewFormatFilesCommand,
	wire.Bind(new(FormatFiles), new(*FormatFilesCommand)),
	NewInstallDependenciesCommand,
	wire.Bind(new(InstallDependencies), new(*InstallDependenciesCommand)),
)
