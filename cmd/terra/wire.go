package main

import (
	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"go.uber.org/dig"
)

func injectAppContext() entities.AppContext {
	container := dig.New()
	
	// Register all providers
	if err := internal.RegisterProviders(container); err != nil {
		panic(err)
	}
	
	// Invoke to get AppInternal and convert to AppContext
	var appInternal *internal.AppInternal
	if err := container.Invoke(func(ai *internal.AppInternal) {
		appInternal = ai
	}); err != nil {
		panic(err)
	}
	
	return newAppContext(appInternal)
}

func newAppContext(appInternal *internal.AppInternal) entities.AppContext {
	return appInternal
}

func injectRootController() entities.Controller {
	container := dig.New()
	
	// Register all providers
	if err := internal.RegisterProviders(container); err != nil {
		panic(err)
	}
	
	// Invoke to get RunFromRootController and convert to Controller
	var rootController *controllers.RunFromRootController
	if err := container.Invoke(func(rc *controllers.RunFromRootController) {
		rootController = rc
	}); err != nil {
		panic(err)
	}
	
	return newRootController(rootController)
}

func newRootController(rootController *controllers.RunFromRootController) entities.Controller {
	return rootController
}
