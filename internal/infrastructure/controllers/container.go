package controllers

import "github.com/google/wire"

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewDeleteCacheController,
	NewFormatFilesController,
	NewInstallDependenciesController,
	// the root command is not included here as controller
)
