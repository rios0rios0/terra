package interfaces

import "github.com/rios0rios0/terra/internal/domain/entities"

type InstallDependencies interface {
	Execute(dependencies []entities.Dependency, listeners InstallDependenciesListeners)
}

type InstallDependenciesListeners struct {
	OnSuccess func()
	OnError   func(error)
}
