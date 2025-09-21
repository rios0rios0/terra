package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// StubInstallDependenciesCommand is a stub implementation of the InstallDependencies interface
type StubInstallDependenciesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *StubInstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}