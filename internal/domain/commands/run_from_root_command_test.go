package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockInstallDependencies is a mock implementation for InstallDependencies interface.
type MockInstallDependencies struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *MockInstallDependencies) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}

// MockFormatFiles is a mock implementation for FormatFiles interface.
type MockFormatFiles struct {
	ExecuteCalled    bool
	LastDependencies []entities.Dependency
}

func (m *MockFormatFiles) Execute(dependencies []entities.Dependency) {
	m.ExecuteCalled = true
	m.LastDependencies = dependencies
}

// MockRunAdditionalBefore is a mock implementation for RunAdditionalBefore interface.
type MockRunAdditionalBefore struct {
	ExecuteCalled  bool
	LastTargetPath string
	LastArguments  []string
}

func (m *MockRunAdditionalBefore) Execute(targetPath string, arguments []string) {
	m.ExecuteCalled = true
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
}

// MockShellRepositoryForRoot is a mock implementation of repositories.ShellRepository
type MockShellRepositoryForRoot struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ExecuteErrors    []error
	CallHistory      []struct {
		Command   string
		Arguments []string
		Directory string
	}
}

func (m *MockShellRepositoryForRoot) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallHistory = append(m.CallHistory, struct {
		Command   string
		Arguments []string
		Directory string
	}{
		Command:   command,
		Arguments: arguments,
		Directory: directory,
	})

	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if len(m.ExecuteErrors) > 0 {
		err := m.ExecuteErrors[0]
		m.ExecuteErrors = m.ExecuteErrors[1:]
		return err
	}
	return nil
}

// MockRootError implements the error interface.
type MockRootError struct {
	message string
}

func (e *MockRootError) Error() string {
	return e.message
}

func TestNewRunFromRootCommand(t *testing.T) {
	t.Parallel()
	
	t.Run("should create instance when valid dependencies provided", func(t *testing.T) {
		// GIVEN: Valid dependencies for creating the command
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()

		// WHEN: Creating a new RunFromRootCommand
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Parallel()
	
	t.Run("should execute all steps when normal execution", func(t *testing.T) {
		// GIVEN: A command with all dependencies and normal arguments
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode"}
		dependencies := []entities.Dependency{
			{
				Name: "terraform",
				CLI:  "terraform",
			},
			{
				Name: "terragrunt",
				CLI:  "terragrunt",
			},
		}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute all preparation steps
		assert.True(t, installCommand.ExecuteCalled, "Should execute install command")
		assert.Equal(t, dependencies, installCommand.LastDependencies)

		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.Equal(t, dependencies, formatCommand.LastDependencies)

		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, targetPath, additionalBefore.LastTargetPath)
		assert.Equal(t, arguments, additionalBefore.LastArguments)

		// Should execute terragrunt with normal repository (not interactive)
		assert.Equal(t, 1, repository.ExecuteCallCount)
		assert.Equal(t, "terragrunt", repository.LastCommand)
		assert.Equal(t, arguments, repository.LastArguments)
		assert.Equal(t, targetPath, repository.LastDirectory)
	})
	
	t.Run("should handle empty arguments when no arguments provided", func(t *testing.T) {
		// GIVEN: A command with empty arguments
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty arguments gracefully
		assert.True(t, installCommand.ExecuteCalled, "Should execute install command")
		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute terragrunt command")
		assert.Len(t, repository.LastArguments, len(arguments), "Should pass arguments with same length")
	})
	
	t.Run("should handle empty dependencies when no dependencies provided", func(t *testing.T) {
		// GIVEN: A command with empty dependencies
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty dependencies gracefully
		assert.True(t, installCommand.ExecuteCalled, "Should execute install command")
		assert.Equal(t, dependencies, installCommand.LastDependencies)
		
		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.Equal(t, dependencies, formatCommand.LastDependencies)

		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute terragrunt command")
	})
	
	t.Run("should pass correct target path when different paths used", func(t *testing.T) {
		// GIVEN: A command with specific target path
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		targetPath := "/custom/terraform/modules/vpc"
		arguments := []string{"validate"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should pass correct target path to all components
		assert.Equal(t, targetPath, additionalBefore.LastTargetPath)
		assert.Equal(t, targetPath, repository.LastDirectory)
	})
	
	t.Run("should not use interactive mode when no auto answer flag", func(t *testing.T) {
		// GIVEN: A command without auto-answer flag in arguments
		installCommand := &MockInstallDependencies{}
		formatCommand := &MockFormatFiles{}
		additionalBefore := &MockRunAdditionalBefore{}
		repository := &MockShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			installCommand,
			formatCommand,
			additionalBefore,
			repository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode", "--out=plan.out"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use normal repository (indirectly tests hasAutoAnswerFlag)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should use normal repository")
		assert.Equal(t, arguments, repository.LastArguments, "Should pass arguments unchanged")
	})
}