package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockFormatFilesCommand is a mock implementation of the FormatFiles interface
type MockFormatFilesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *MockFormatFilesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}