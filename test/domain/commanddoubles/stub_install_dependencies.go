//go:build integration || unit || test

package commanddoubles //nolint:staticcheck // Test package naming follows established project structure

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
