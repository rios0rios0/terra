//go:build !test && wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
)

func injectAppContext() entities.AppContext {
	wire.Build(
		internal.Container,
		newAppContext,
	)
	return nil
}

func newAppContext(appInternal *internal.AppInternal) entities.AppContext {
	return appInternal
}

func injectRootController() entities.Controller {
	// this way you avoid repeating the root controller with the other controllers
	wire.Build(
		internal.Container,
		controllers.NewRunFromRootController,
		newRootController,
	)
	return nil
}

func newRootController(rootController *controllers.RunFromRootController) entities.Controller {
	return rootController
}
