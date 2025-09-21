package repositories_doubles

import "github.com/rios0rios0/terra/test/domain/entities_doubles"

// StubShellRepository for testing shell-related commands
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
		return entities_doubles.NewStubError("stub execution error")
	}
	return nil
}