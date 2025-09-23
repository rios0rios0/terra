//go:build integration || unit || test

package repositorydoubles

import "github.com/rios0rios0/terra/internal/domain/repositories"

// ParallelStateCallRecord represents a single command execution call for parallel state testing
type ParallelStateCallRecord struct {
	Command   string
	Arguments []string
	Directory string
}

// StubShellRepositoryForParallelState is a test double for shell repository focused on parallel state testing
type StubShellRepositoryForParallelState struct {
	ExecuteCallCount int
	CallHistory      []ParallelStateCallRecord
	ShouldFail       bool
	FailureMessage   string
}

// Verify it implements the interface
var _ repositories.ShellRepository = (*StubShellRepositoryForParallelState)(nil)

func (stub *StubShellRepositoryForParallelState) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	stub.ExecuteCallCount++
	stub.CallHistory = append(stub.CallHistory, ParallelStateCallRecord{
		Command:   command,
		Arguments: make([]string, len(arguments)),
		Directory: directory,
	})
	
	// Copy arguments to avoid modification issues
	copy(stub.CallHistory[len(stub.CallHistory)-1].Arguments, arguments)
	
	if stub.ShouldFail {
		return &stubParallelStateError{message: stub.FailureMessage}
	}
	
	return nil
}

// stubParallelStateError represents a simple error for testing
type stubParallelStateError struct {
	message string
}

func (e *stubParallelStateError) Error() string {
	if e.message == "" {
		return "simulated command failure"
	}
	return e.message
}