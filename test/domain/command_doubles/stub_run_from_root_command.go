package command_doubles

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// StubRunFromRootCommand is a stub implementation of the RunFromRoot interface.
type StubRunFromRootCommand struct {
	ExecuteCallCount int
	LastTargetPath   string
	LastArguments    []string
	LastDependencies []entities.Dependency
}

func (m *StubRunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	m.ExecuteCallCount++
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
	m.LastDependencies = dependencies
}
