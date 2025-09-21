package repositories_doubles

import "github.com/rios0rios0/terra/test"

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
		return test.NewStubError("stub execution error")
	}
	return nil
}