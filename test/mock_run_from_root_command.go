package test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockRunFromRootCommand is a mock implementation of the RunFromRoot interface
type MockRunFromRootCommand struct {
	ExecuteCallCount int
	LastTargetPath   string
	LastArguments    []string
	LastDependencies []entities.Dependency
}

func (m *MockRunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	m.ExecuteCallCount++
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
	m.LastDependencies = dependencies
}