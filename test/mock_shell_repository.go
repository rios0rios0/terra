package test

// MockShellRepository for testing shell-related commands
type MockShellRepository struct {
	ExecuteCallCount  int
	LastCommand       string
	LastArguments     []string
	LastDirectory     string
	ShouldReturnError bool
}

func (m *MockShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if m.ShouldReturnError {
		return &MockError{message: "mock execution error"}
	}
	return nil
}