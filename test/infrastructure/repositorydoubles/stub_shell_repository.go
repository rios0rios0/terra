//go:build integration || unit || test

package repositorydoubles //nolint:staticcheck // Test package naming follows established project structure

import "github.com/rios0rios0/terra/test/domain/entitydoubles"

// StubShellRepository for testing shell-related commands.
type StubShellRepository struct {
	ExecuteCallCount  int
	LastCommand       string
	LastArguments     []string
	LastDirectory     string
	ShouldReturnError bool
}

func (m *StubShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if m.ShouldReturnError {
		return entitydoubles.NewStubError("stub execution error")
	}
	return nil
}
