package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewFormatFilesCommand_ShouldCreateInstance_WhenRepositoryProvided(t *testing.T) {
	// GIVEN: A mock shell repository
	mockRepo := &MockShellRepository{}

	// WHEN: Creating a new format files command
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// THEN: Should create a valid command instance
	require.NotNil(t, cmd)
}

func TestFormatFilesCommand_ShouldExecuteFormatCommands_WhenDependenciesProvided(t *testing.T) {
	// GIVEN: A mock repository and dependencies with formatting commands
	mockRepo := &MockShellRepository{}
	terraformDep := entities.Dependency{
		Name:              "Terraform",
		CLI:               "terraform",
		FormattingCommand: []string{"fmt", "-recursive"},
	}
	terragruntDep := entities.Dependency{
		Name:              "Terragrunt",
		CLI:               "terragrunt",
		FormattingCommand: []string{"hcl", "format", "**/*.hcl"},
	}
	dependencies := []entities.Dependency{terraformDep, terragruntDep}
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// WHEN: Executing the format command
	cmd.Execute(dependencies)

	// THEN: Should execute command for each dependency
	assert.Equal(t, len(dependencies), mockRepo.ExecuteCallCount)
	assert.Equal(t, terragruntDep.CLI, mockRepo.LastCommand)
	assert.Equal(t, terragruntDep.FormattingCommand, mockRepo.LastArguments)
	assert.Equal(t, ".", mockRepo.LastDirectory)
}

func TestFormatFilesCommand_ShouldContinueExecution_WhenRepositoryReturnsError(t *testing.T) {
	// GIVEN: A mock repository that returns errors and a single dependency
	mockRepo := &MockShellRepository{ShouldReturnError: true}
	dependencies := []entities.Dependency{
		{
			Name:              "Terraform",
			CLI:               "terraform",
			FormattingCommand: []string{"fmt", "-recursive"},
		},
	}
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// WHEN: Executing the format command
	cmd.Execute(dependencies)

	// THEN: Should execute command despite the error (command handles errors gracefully)
	assert.Equal(t, 1, mockRepo.ExecuteCallCount)
}

func TestFormatFilesCommand_ShouldNotExecute_WhenNoDependenciesProvided(t *testing.T) {
	// GIVEN: A mock repository and empty dependencies list
	mockRepo := &MockShellRepository{}
	dependencies := []entities.Dependency{}
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// WHEN: Executing the format command
	cmd.Execute(dependencies)

	// THEN: Should not execute any commands
	assert.Equal(t, 0, mockRepo.ExecuteCallCount)
}

func TestFormatFilesCommand_ShouldExecuteWithEmptyArguments_WhenDependencyHasNoFormattingCommand(t *testing.T) {
	// GIVEN: A mock repository and dependency with empty formatting command
	mockRepo := &MockShellRepository{}
	dependencies := []entities.Dependency{
		{
			Name:              "SomeTool",
			CLI:               "sometool",
			FormattingCommand: []string{}, // Empty formatting command
		},
	}
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// WHEN: Executing the format command
	cmd.Execute(dependencies)

	// THEN: Should execute command with empty arguments
	assert.Equal(t, 1, mockRepo.ExecuteCallCount)
	assert.Equal(t, "sometool", mockRepo.LastCommand)
	assert.Empty(t, mockRepo.LastArguments)
}

func TestFormatFilesCommand_ShouldExecuteAllDependencies_WhenMultipleDependenciesProvided(t *testing.T) {
	// GIVEN: A recording mock repository and multiple dependencies
	mockRepo := &MockShellRepositoryWithRecording{}
	terraformDep := entities.Dependency{
		Name:              "Terraform",
		CLI:               "terraform",
		FormattingCommand: []string{"fmt", "-recursive"},
	}
	terragruntDep := entities.Dependency{
		Name:              "Terragrunt",
		CLI:               "terragrunt",
		FormattingCommand: []string{"hcl", "format", "**/*.hcl"},
	}
	customDep := entities.Dependency{
		Name:              "CustomTool",
		CLI:               "customtool",
		FormattingCommand: []string{"format", "--all"},
	}
	dependencies := []entities.Dependency{terraformDep, terragruntDep, customDep}
	cmd := commands.NewFormatFilesCommand(mockRepo)

	// WHEN: Executing the format command
	cmd.Execute(dependencies)

	// THEN: Should execute all dependencies in order
	require.Equal(t, len(dependencies), len(mockRepo.CallRecords))
	
	// Verify first call (Terraform)
	firstRecord := mockRepo.CallRecords[0]
	assert.Equal(t, terraformDep.CLI, firstRecord.Command)
	assert.Equal(t, terraformDep.FormattingCommand, firstRecord.Arguments)
	assert.Equal(t, ".", firstRecord.Directory)
	
	// Verify second call (Terragrunt)
	secondRecord := mockRepo.CallRecords[1]
	assert.Equal(t, terragruntDep.CLI, secondRecord.Command)
	assert.Equal(t, terragruntDep.FormattingCommand, secondRecord.Arguments)
	assert.Equal(t, ".", secondRecord.Directory)
	
	// Verify third call (CustomTool)
	thirdRecord := mockRepo.CallRecords[2]
	assert.Equal(t, customDep.CLI, thirdRecord.Command)
	assert.Equal(t, customDep.FormattingCommand, thirdRecord.Arguments)
	assert.Equal(t, ".", thirdRecord.Directory)
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
