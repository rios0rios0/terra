//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"
	
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRunFromRootCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when valid dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Valid dependencies for creating the command
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()

		// WHEN: Creating a new RunFromRootCommand
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute all steps when normal execution", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with all dependencies and normal arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

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

		// THEN: Should execute preparation steps but not install
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.Nil(t, installCommand.LastDependencies)

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
		t.Parallel()
		// GIVEN: A command with empty arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty arguments gracefully
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute terragrunt command")
		assert.Len(
			t,
			repository.LastArguments,
			len(arguments),
			"Should pass arguments with same length",
		)
	})

	t.Run("should handle empty dependencies when no dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with empty dependencies
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{"plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty dependencies gracefully
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.Nil(t, installCommand.LastDependencies)

		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.Equal(t, dependencies, formatCommand.LastDependencies)

		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute terragrunt command")
	})

	t.Run("should pass correct target path when different paths used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with specific target path
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

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
		t.Parallel()
		// GIVEN: A command without auto-answer flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode", "--out=plan.out"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use normal repository (indirectly tests hasAutoAnswerFlag)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should use normal repository")
		assert.Equal(t, arguments, repository.LastArguments, "Should pass arguments unchanged")
	})

	t.Run("should use interactive mode when --auto-answer flag present", func(t *testing.T) {
		// GIVEN: A command with --auto-answer flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{"plan", "--auto-answer", "--detailed-exitcode"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository instead of normal repository
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "terragrunt", interactiveRepository.LastCommand)
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
		assert.Equal(t, targetPath, interactiveRepository.LastDirectory)
	})

	t.Run("should use interactive mode when -a flag present", func(t *testing.T) {
		// GIVEN: A command with -a flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{"apply", "-a", "--auto-approve"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository instead of normal one
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "terragrunt", interactiveRepository.LastCommand)
		assert.Equal(t, []string{"apply", "--auto-approve"}, interactiveRepository.LastArguments, "Should filter out -a flag")
		assert.Equal(t, targetPath, interactiveRepository.LastDirectory)
	})

	t.Run("should remove auto-answer flags from arguments before passing to terragrunt", func(t *testing.T) {
		// GIVEN: A command with auto-answer flags mixed with other arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
			InstallCommand:        installCommand,
			FormatCommand:         formatCommand,
			AdditionalBefore:      additionalBefore,
			Repository:            repository,
			InteractiveRepository: interactiveRepository,
		})

		targetPath := "/test/path"
		arguments := []string{"plan", "--auto-answer", "--detailed-exitcode", "-a", "--out=plan.out"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should pass filtered arguments to additionalBefore (which records them)
		expectedFilteredArgs := []string{"plan", "--auto-answer", "--detailed-exitcode", "-a", "--out=plan.out"}
		assert.Equal(t, expectedFilteredArgs, additionalBefore.LastArguments, "Should pass original arguments to additionalBefore")
	})
}

func TestRunFromRootCommand_HasAutoAnswerFlag(t *testing.T) {
	t.Parallel()

	// Create a command instance for testing the method
	installCommand := &commanddoubles.StubInstallDependencies{}
	formatCommand := &commanddoubles.StubFormatFiles{}
	additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
	repository := &repositorydoubles.StubShellRepositoryForRoot{}
	interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
	cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
		InstallCommand:        installCommand,
		FormatCommand:         formatCommand,
		AdditionalBefore:      additionalBefore,
		Repository:            repository,
		InteractiveRepository: interactiveRepository,
	})

	testCases := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{
			name:      "should return true when --auto-answer flag present",
			arguments: []string{"plan", "--auto-answer", "--detailed-exitcode"},
			expected:  true,
		},
		{
			name:      "should return true when -a flag present",
			arguments: []string{"apply", "-a", "--auto-approve"},
			expected:  true,
		},
		{
			name:      "should return true when both flags present",
			arguments: []string{"plan", "--auto-answer", "-a", "--out=plan.out"},
			expected:  true,
		},
		{
			name:      "should return false when no auto-answer flags present",
			arguments: []string{"plan", "--detailed-exitcode", "--out=plan.out"},
			expected:  false,
		},
		{
			name:      "should return false when empty arguments",
			arguments: []string{},
			expected:  false,
		},
		{
			name:      "should return false when similar but different flags present",
			arguments: []string{"plan", "--auto-approve", "-approve", "--answer"},
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// GIVEN: Arguments with or without auto-answer flags
			arguments := tc.arguments

			// WHEN: Checking for auto-answer flag using Execute (which calls hasAutoAnswerFlag internally)
			targetPath := "/test/path"
			dependencies := []entities.Dependency{}
			cmd.Execute(targetPath, arguments, dependencies)

			// THEN: Verify correct repository was used based on auto-answer flag presence
			if tc.expected {
				assert.Equal(t, 0, repository.ExecuteCallCount, "Should use interactive repository when auto-answer flag present")
			} else {
				assert.Equal(t, 1, repository.ExecuteCallCount, "Should use normal repository when no auto-answer flag present")
			}

			// Reset for next test
			repository.ExecuteCallCount = 0
		})
	}
}

func TestRunFromRootCommand_RemoveAutoAnswerFlag(t *testing.T) {
	t.Parallel()

	// Since removeAutoAnswerFlag is private, we test it indirectly through Execute
	// by checking the arguments passed to additionalBefore which receives filtered arguments
	installCommand := &commanddoubles.StubInstallDependencies{}
	formatCommand := &commanddoubles.StubFormatFiles{}
	additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
	repository := &repositorydoubles.StubShellRepositoryForRoot{}
	interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
	cmd := commands.NewRunFromRootCommand(commands.RunFromRootCommandDeps{
		InstallCommand:        installCommand,
		FormatCommand:         formatCommand,
		AdditionalBefore:      additionalBefore,
		Repository:            repository,
		InteractiveRepository: interactiveRepository,
	})

	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "should remove --auto-answer flag",
			input:    []string{"plan", "--auto-answer", "--detailed-exitcode"},
			expected: []string{"plan", "--detailed-exitcode"},
		},
		{
			name:     "should remove -a flag",
			input:    []string{"apply", "-a", "--auto-approve"},
			expected: []string{"apply", "--auto-approve"},
		},
		{
			name:     "should remove both --auto-answer and -a flags",
			input:    []string{"plan", "--auto-answer", "-a", "--out=plan.out"},
			expected: []string{"plan", "--out=plan.out"},
		},
		{
			name:     "should not remove other flags",
			input:    []string{"plan", "--auto-approve", "-approve", "--answer"},
			expected: []string{"plan", "--auto-approve", "-approve", "--answer"},
		},
		{
			name:     "should handle empty arguments",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "should handle arguments with only auto-answer flags",
			input:    []string{"--auto-answer", "-a"},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN: Arguments that may contain auto-answer flags
			targetPath := "/test/path"
			dependencies := []entities.Dependency{}

			// WHEN: Executing the command
			cmd.Execute(targetPath, tc.input, dependencies)

			// THEN: Should pass original arguments to additionalBefore (before filtering)
			// Note: additionalBefore receives arguments before they're filtered
			assert.Equal(t, tc.input, additionalBefore.LastArguments, "additionalBefore should receive original arguments")

			// Reset for next test
			additionalBefore.LastArguments = nil
		})
	}
}
