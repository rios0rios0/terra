package commands

import "github.com/google/wire"

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewDeleteCacheCommand,
	wire.Bind(new(DeleteCache), new(*DeleteCacheCommand)),
	NewFormatFilesCommand,
	wire.Bind(new(FormatFiles), new(*FormatFilesCommand)),
	NewInstallDependenciesCommand,
	wire.Bind(new(InstallDependencies), new(*InstallDependenciesCommand)),
	NewRunAdditionalBeforeCommand,
	wire.Bind(new(RunAdditionalBefore), new(*RunAdditionalBeforeCommand)),
	NewRunFromRootCommand,
	wire.Bind(new(RunFromRoot), new(*RunFromRootCommand)),
)
