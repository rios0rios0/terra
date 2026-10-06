package repositorydoubles

import (
	"sync"

	"github.com/rios0rios0/terra/internal/domain/repositories"
)

// ParallelStateCallRecord represents a single command execution call for parallel state testing.
type ParallelStateCallRecord struct {
	Command   string
	Arguments []string
	Directory string
	Prefix    string
}

// StubShellRepositoryForParallelState is a test double for shell repository focused on parallel state testing.
// The command under test calls it from several workers at once, so it records each call under a lock; a test
// reads the recorded calls once the command has returned, by which time every worker has finished.
type StubShellRepositoryForParallelState struct {
	ExecuteCallCount int
	CallHistory      []ParallelStateCallRecord
	ShouldFail       bool
	FailureMessage   string

	mu sync.Mutex
}

// Verify it implements the interface.
var _ repositories.ParallelShellRepository = (*StubShellRepositoryForParallelState)(nil)

func (stub *StubShellRepositoryForParallelState) ExecuteCommandWithPrefix(
	command string,
	arguments []string,
	directory string,
	prefix string,
) error {
	stub.mu.Lock()
	defer stub.mu.Unlock()

	stub.ExecuteCallCount++
	stub.CallHistory = append(stub.CallHistory, ParallelStateCallRecord{
		Command:   command,
		Arguments: make([]string, len(arguments)),
		Directory: directory,
		Prefix:    prefix,
	})

	// Copy arguments to avoid modification issues
	copy(stub.CallHistory[len(stub.CallHistory)-1].Arguments, arguments)

	if stub.ShouldFail {
		return &stubParallelStateError{message: stub.FailureMessage}
	}

	return nil
}

// stubParallelStateError represents a simple error for testing.
type stubParallelStateError struct {
	message string
}

func (e *stubParallelStateError) Error() string {
	if e.message == "" {
		return "simulated command failure"
	}
	return e.message
}
