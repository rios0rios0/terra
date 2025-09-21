package commands_doubles

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// StubFormatFiles is a stub implementation for FormatFiles interface.
type StubFormatFiles struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *StubFormatFiles) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}