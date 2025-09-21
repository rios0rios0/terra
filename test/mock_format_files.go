package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockFormatFiles is a mock implementation for FormatFiles interface.
type MockFormatFiles struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *MockFormatFiles) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}