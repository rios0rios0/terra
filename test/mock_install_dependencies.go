package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockInstallDependencies is a mock implementation for InstallDependencies interface.
type MockInstallDependencies struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *MockInstallDependencies) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}