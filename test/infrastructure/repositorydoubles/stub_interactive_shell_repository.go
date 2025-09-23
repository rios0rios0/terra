//go:build integration || unit || test

package repositorydoubles

import "github.com/rios0rios0/terra/internal/domain/repositories"

// StubInteractiveShellRepository is a test stub for InteractiveShellRepository.
type StubInteractiveShellRepository struct {
	ExecuteCallCount            int
	ExecuteWithAnswerCallCount  int
	LastCommand                 string
	LastArguments               []string
	LastDirectory               string
	LastAutoAnswer              string
	ShouldReturnError           bool
	ExecuteError                error
}

// Ensure StubInteractiveShellRepository implements the interface
var _ repositories.InteractiveShellRepository = (*StubInteractiveShellRepository)(nil)

func (s *StubInteractiveShellRepository) ExecuteCommand(command string, arguments []string, directory string) error {
	s.ExecuteCallCount++
	s.LastCommand = command
	s.LastArguments = arguments
	s.LastDirectory = directory
	
	if s.ShouldReturnError {
		return s.ExecuteError
	}
	return nil
}

func (s *StubInteractiveShellRepository) ExecuteCommandWithAnswer(command string, arguments []string, directory string, autoAnswer string) error {
	s.ExecuteWithAnswerCallCount++
	s.LastCommand = command
	s.LastArguments = arguments
	s.LastDirectory = directory
	s.LastAutoAnswer = autoAnswer
	
	if s.ShouldReturnError {
		return s.ExecuteError
	}
	return nil
}