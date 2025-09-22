package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/infrastructure/repository_doubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatFilesCommand(t *testing.T) {
	t.Parallel()
	
	t.Run("should create instance when repository provided", func(t *testing.T) {
		// GIVEN: A mock shell repository
		mockRepo := &repository_doubles.StubShellRepository{}

		// WHEN: Creating a new format files command
		cmd := commands.NewFormatFilesCommand(mockRepo)

		// THEN: Should create a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestFormatFilesCommand_Execute(t *testing.T) {
	t.Parallel()
	
	t.Run("should execute format commands when dependencies provided", func(t *testing.T) {
		// GIVEN: A mock repository and dependencies with formatting commands
		mockRepo := &repository_doubles.StubShellRepository{}
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
	})
	
	t.Run("should continue execution when repository returns error", func(t *testing.T) {
		// GIVEN: A mock repository that returns errors and a single dependency
		mockRepo := &repository_doubles.StubShellRepository{ShouldReturnError: true}
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
	})
	
	t.Run("should not execute when no dependencies provided", func(t *testing.T) {
		// GIVEN: A mock repository and empty dependencies list
		mockRepo := &repository_doubles.StubShellRepository{}
		dependencies := []entities.Dependency{}
		cmd := commands.NewFormatFilesCommand(mockRepo)

		// WHEN: Executing the format command
		cmd.Execute(dependencies)

		// THEN: Should not execute any commands
		assert.Equal(t, 0, mockRepo.ExecuteCallCount)
	})
	
	t.Run("should execute with empty arguments when dependency has no formatting command", func(t *testing.T) {
		// GIVEN: A mock repository and dependency with empty formatting command
		mockRepo := &repository_doubles.StubShellRepository{}
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
	})
	
	t.Run("should execute all dependencies when multiple dependencies provided", func(t *testing.T) {
		// GIVEN: A recording mock repository and multiple dependencies
		mockRepo := &repository_doubles.StubShellRepositoryWithRecording{}
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
	})
}
