package commands_doubles

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// StubInstallDependencies is a stub implementation for InstallDependencies interface.
type StubInstallDependencies struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *StubInstallDependencies) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}