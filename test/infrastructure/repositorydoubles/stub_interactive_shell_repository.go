//go:build integration || unit || test

package repositorydoubles //nolint:staticcheck // Test package naming follows established project structure

// StubInteractiveShellRepository is a stub implementation of InteractiveShellRepository.
type StubInteractiveShellRepository struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ExecuteErrors    []error
	CallHistory      []struct {
		Command   string
		Arguments []string
		Directory string
	}
	autoAnswerValue string
}

func (m *StubInteractiveShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallHistory = append(m.CallHistory, struct {
		Command   string
		Arguments []string
		Directory string
	}{
		Command:   command,
		Arguments: arguments,
		Directory: directory,
	})

	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if len(m.ExecuteErrors) > 0 {
		err := m.ExecuteErrors[0]
		m.ExecuteErrors = m.ExecuteErrors[1:]
		return err
	}
	return nil
}

// SetAutoAnswerValue sets the auto-answer preference (for testing compatibility)
func (m *StubInteractiveShellRepository) SetAutoAnswerValue(value string) {
	m.autoAnswerValue = value
}

// GetAutoAnswerValue returns the configured auto-answer value (for testing)
func (m *StubInteractiveShellRepository) GetAutoAnswerValue() string {
	return m.autoAnswerValue
}