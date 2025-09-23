//go:build integration || unit || test

package commanddoubles

import "github.com/rios0rios0/terra/internal/domain/entities"

// StubParallelState is a test double for parallel state commands
type StubParallelState struct {
	ExecuteCalled       bool
	LastTargetPath      string
	LastArguments       []string
	LastDependencies    []entities.Dependency
	ShouldReturnError   bool
	ErrorMessage        string
}

func (stub *StubParallelState) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) error {
	stub.ExecuteCalled = true
	stub.LastTargetPath = targetPath
	stub.LastArguments = make([]string, len(arguments))
	copy(stub.LastArguments, arguments)
	stub.LastDependencies = make([]entities.Dependency, len(dependencies))
	copy(stub.LastDependencies, dependencies)
	
	if stub.ShouldReturnError {
		return &stubError{message: stub.ErrorMessage}
	}
	
	return nil
}

// stubError represents a simple error for testing
type stubError struct {
	message string
}

func (e *stubError) Error() string {
	if e.message == "" {
		return "simulated parallel state error"
	}
	return e.message
}