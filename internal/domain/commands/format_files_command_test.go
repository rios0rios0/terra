package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

// MockShellRepository for testing FormatFilesCommand
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

// MockError implements the error interface
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestNewFormatFilesCommand(t *testing.T) {
	mockRepo := &MockShellRepository{}
	cmd := NewFormatFilesCommand(mockRepo)

	if cmd == nil {
		t.Fatal("NewFormatFilesCommand returned nil")
	}

	if cmd.repository != mockRepo {
		t.Error("Repository was not set correctly")
	}
}

func TestFormatFilesCommand_Execute(t *testing.T) {
	mockRepo := &MockShellRepository{}
	dependencies := []entities.Dependency{
		{
			Name:              "Terraform",
			CLI:               "terraform",
			FormattingCommand: []string{"fmt", "-recursive"},
		},
		{
			Name:              "Terragrunt",
			CLI:               "terragrunt",
			FormattingCommand: []string{"hcl", "format", "**/*.hcl"},
		},
	}

	cmd := NewFormatFilesCommand(mockRepo)
	cmd.Execute(dependencies)

	// Verify that ExecuteCommand was called for each dependency
	expectedCallCount := len(dependencies)
	if mockRepo.ExecuteCallCount != expectedCallCount {
		t.Errorf(
			"Expected ExecuteCommand to be called %d times, got %d",
			expectedCallCount,
			mockRepo.ExecuteCallCount,
		)
	}

	// Since we can't easily verify all calls in this simple mock,
	// we can at least verify the last call was correct
	lastDep := dependencies[len(dependencies)-1]
	if mockRepo.LastCommand != lastDep.CLI {
		t.Errorf("Expected last command to be %s, got %s", lastDep.CLI, mockRepo.LastCommand)
	}

	if len(mockRepo.LastArguments) != len(lastDep.FormattingCommand) {
		t.Errorf(
			"Expected %d arguments, got %d",
			len(lastDep.FormattingCommand),
			len(mockRepo.LastArguments),
		)
	}

	for i, expectedArg := range lastDep.FormattingCommand {
		if i < len(mockRepo.LastArguments) && mockRepo.LastArguments[i] != expectedArg {
			t.Errorf(
				"Expected argument %s at index %d, got %s",
				expectedArg,
				i,
				mockRepo.LastArguments[i],
			)
		}
	}

	if mockRepo.LastDirectory != "." {
		t.Errorf("Expected directory to be '.', got %s", mockRepo.LastDirectory)
	}
}

func TestFormatFilesCommand_ExecuteWithError(t *testing.T) {
	mockRepo := &MockShellRepository{ShouldReturnError: true}
	dependencies := []entities.Dependency{
		{
			Name:              "Terraform",
			CLI:               "terraform",
			FormattingCommand: []string{"fmt", "-recursive"},
		},
	}

	cmd := NewFormatFilesCommand(mockRepo)

	// This should not panic even if the repository returns an error
	// The command should log the error but continue execution
	cmd.Execute(dependencies)

	// Verify that ExecuteCommand was called
	if mockRepo.ExecuteCallCount != 1 {
		t.Errorf("Expected ExecuteCommand to be called once, got %d", mockRepo.ExecuteCallCount)
	}
}

func TestFormatFilesCommand_ExecuteWithEmptyDependencies(t *testing.T) {
	mockRepo := &MockShellRepository{}
	dependencies := []entities.Dependency{}

	cmd := NewFormatFilesCommand(mockRepo)
	cmd.Execute(dependencies)

	// Verify that ExecuteCommand was not called
	if mockRepo.ExecuteCallCount != 0 {
		t.Errorf(
			"Expected ExecuteCommand to not be called, got %d calls",
			mockRepo.ExecuteCallCount,
		)
	}
}

func TestFormatFilesCommand_ExecuteWithDependencyWithoutFormattingCommand(t *testing.T) {
	mockRepo := &MockShellRepository{}
	dependencies := []entities.Dependency{
		{
			Name:              "SomeTool",
			CLI:               "sometool",
			FormattingCommand: []string{}, // Empty formatting command
		},
	}

	cmd := NewFormatFilesCommand(mockRepo)
	cmd.Execute(dependencies)

	// Verify that ExecuteCommand was called even with empty arguments
	if mockRepo.ExecuteCallCount != 1 {
		t.Errorf("Expected ExecuteCommand to be called once, got %d", mockRepo.ExecuteCallCount)
	}

	if mockRepo.LastCommand != "sometool" {
		t.Errorf("Expected command to be 'sometool', got %s", mockRepo.LastCommand)
	}

	if len(mockRepo.LastArguments) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(mockRepo.LastArguments))
	}
}

func TestFormatFilesCommand_ExecuteWithMultipleDependencies(t *testing.T) {
	// Create a new mock that records all calls
	mockRepo := &MockShellRepositoryWithRecording{}

	dependencies := []entities.Dependency{
		{
			Name:              "Terraform",
			CLI:               "terraform",
			FormattingCommand: []string{"fmt", "-recursive"},
		},
		{
			Name:              "Terragrunt",
			CLI:               "terragrunt",
			FormattingCommand: []string{"hcl", "format", "**/*.hcl"},
		},
		{
			Name:              "CustomTool",
			CLI:               "customtool",
			FormattingCommand: []string{"format", "--all"},
		},
	}

	cmd := NewFormatFilesCommand(mockRepo)
	cmd.Execute(dependencies)

	// Verify correct number of calls
	if len(mockRepo.CallRecords) != len(dependencies) {
		t.Errorf("Expected %d calls, got %d", len(dependencies), len(mockRepo.CallRecords))
	}

	// Verify each call
	for i, dep := range dependencies {
		if i >= len(mockRepo.CallRecords) {
			continue
		}

		record := mockRepo.CallRecords[i]
		if record.Command != dep.CLI {
			t.Errorf("Call %d: expected command %s, got %s", i, dep.CLI, record.Command)
		}

		if len(record.Arguments) != len(dep.FormattingCommand) {
			t.Errorf(
				"Call %d: expected %d arguments, got %d",
				i,
				len(dep.FormattingCommand),
				len(record.Arguments),
			)
		}

		for j, expectedArg := range dep.FormattingCommand {
			if j < len(record.Arguments) && record.Arguments[j] != expectedArg {
				t.Errorf(
					"Call %d, arg %d: expected %s, got %s",
					i,
					j,
					expectedArg,
					record.Arguments[j],
				)
			}
		}

		if record.Directory != "." {
			t.Errorf("Call %d: expected directory '.', got %s", i, record.Directory)
		}
	}
}

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
