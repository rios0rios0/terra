//go:build integration || unit || test

package repository_doubles //nolint:staticcheck // Test package naming follows established project structure

import "github.com/rios0rios0/terra/test/domain/entity_doubles"

// CallRecord represents a single repository call.
type CallRecord struct {
	Command   string
	Arguments []string
	Directory string
}

// StubShellRepositoryWithRecording for testing with call recording.
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
		Arguments: append([]string{}, arguments...), // Copy slice.
		Directory: directory,
	})

	if m.ShouldReturnError {
		return entity_doubles.NewStubError("stub execution error")
	}
	return nil
}
