package test

// CallRecord represents a single repository call
type CallRecord struct {
	Command   string
	Arguments []string
	Directory string
}

// MockShellRepositoryWithRecording for testing with call recording
type MockShellRepositoryWithRecording struct {
	CallRecords       []CallRecord
	ShouldReturnError bool
}

func (m *MockShellRepositoryWithRecording) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallRecords = append(m.CallRecords, CallRecord{
		Command:   command,
		Arguments: append([]string{}, arguments...), // Copy slice
		Directory: directory,
	})

	if m.ShouldReturnError {
		return &MockError{message: "mock execution error"}
	}
	return nil
}