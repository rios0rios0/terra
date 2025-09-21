package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockInstallDependenciesCommand is a mock implementation of the InstallDependencies interface
type MockInstallDependenciesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *MockInstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}