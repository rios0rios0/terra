package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// StubFormatFilesCommand is a stub implementation of the FormatFiles interface
type StubFormatFilesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *StubFormatFilesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}