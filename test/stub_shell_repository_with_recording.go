package test

// CallRecord represents a single repository call
type CallRecord struct {
	Command   string
	Arguments []string
	Directory string
}

// StubShellRepositoryWithRecording for testing with call recording
type StubShellRepositoryWithRecording struct {
	CallRecords       []CallRecord
	ShouldReturnError bool
}

func (m *StubShellRepositoryWithRecording) ExecuteCommand(
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
		return &StubError{message: "stub execution error"}
	}
	return nil
}