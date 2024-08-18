package controllers

import "github.com/google/wire"

var Container = wire.NewSet(
	NewDeleteCacheController,
	NewFormatFilesController,
	NewInstallDependenciesController,
)
