//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunFromRootCommand_AutoAnswerEnhanced(t *testing.T) {
	t.Parallel()

	t.Run("should handle --auto-answer=y flag", func(t *testing.T) {
		// GIVEN: A command with --auto-answer=y flag
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
		arguments := []string{"apply", "--auto-answer=y", "--auto-approve"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository and configure it with "y"
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "y", interactiveRepository.GetAutoAnswerValue(), "Should configure auto-answer value to 'y'")
		assert.Equal(t, []string{"apply", "--auto-approve"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
	})

	t.Run("should handle --auto-answer=n flag", func(t *testing.T) {
		// GIVEN: A command with --auto-answer=n flag
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
		arguments := []string{"plan", "--auto-answer=n", "--detailed-exitcode"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository and configure it with "n"
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "n", interactiveRepository.GetAutoAnswerValue(), "Should configure auto-answer value to 'n'")
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
	})

	t.Run("should handle --auto-answer=yes flag", func(t *testing.T) {
		// GIVEN: A command with --auto-answer=yes flag
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
		arguments := []string{"apply", "--auto-answer=yes"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository and configure it with "y"
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "y", interactiveRepository.GetAutoAnswerValue(), "Should configure auto-answer value to 'y' for 'yes'")
	})

	t.Run("should handle --auto-answer=no flag", func(t *testing.T) {
		// GIVEN: A command with --auto-answer=no flag
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
		arguments := []string{"plan", "--auto-answer=no"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository and configure it with "n"
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "n", interactiveRepository.GetAutoAnswerValue(), "Should configure auto-answer value to 'n' for 'no'")
	})

	t.Run("should default to 'n' for invalid auto-answer values", func(t *testing.T) {
		// GIVEN: A command with invalid auto-answer value
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
		arguments := []string{"plan", "--auto-answer=invalid"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository and default to "n"
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "n", interactiveRepository.GetAutoAnswerValue(), "Should default to 'n' for invalid values")
	})

	t.Run("should maintain backward compatibility with --auto-answer flag", func(t *testing.T) {
		// GIVEN: A command with legacy --auto-answer flag (no value)
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

		// THEN: Should use interactive repository and default to "n" for backward compatibility
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "n", interactiveRepository.GetAutoAnswerValue(), "Should default to 'n' for backward compatibility")
	})

	t.Run("should maintain backward compatibility with -a flag", func(t *testing.T) {
		// GIVEN: A command with legacy -a flag
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

		// THEN: Should use interactive repository and default to "n" for backward compatibility
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
		assert.Equal(t, "n", interactiveRepository.GetAutoAnswerValue(), "Should default to 'n' for backward compatibility")
	})
}

func TestRunFromRootCommand_GetAutoAnswerValue(t *testing.T) {
	t.Parallel()

	// Create a command instance for testing the method indirectly
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

	testCases := []struct {
		name        string
		arguments   []string
		expected    string
		shouldUseInteractive bool
	}{
		{
			name:        "should return 'y' for --auto-answer=y",
			arguments:   []string{"plan", "--auto-answer=y", "--detailed-exitcode"},
			expected:    "y",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'n' for --auto-answer=n",
			arguments:   []string{"apply", "--auto-answer=n", "--auto-approve"},
			expected:    "n",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'y' for --auto-answer=yes",
			arguments:   []string{"plan", "--auto-answer=yes"},
			expected:    "y",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'n' for --auto-answer=no",
			arguments:   []string{"apply", "--auto-answer=no"},
			expected:    "n",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'n' for --auto-answer flag without value",
			arguments:   []string{"plan", "--auto-answer", "--detailed-exitcode"},
			expected:    "n",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'n' for -a flag",
			arguments:   []string{"apply", "-a", "--auto-approve"},
			expected:    "n",
			shouldUseInteractive: true,
		},
		{
			name:        "should return 'n' for invalid value",
			arguments:   []string{"plan", "--auto-answer=invalid"},
			expected:    "n",
			shouldUseInteractive: true,
		},
		{
			name:        "should return empty string when no auto-answer flag",
			arguments:   []string{"plan", "--detailed-exitcode"},
			expected:    "",
			shouldUseInteractive: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset state
			repository.ExecuteCallCount = 0
			interactiveRepository.ExecuteCallCount = 0

			// WHEN: Executing the command to test getAutoAnswerValue indirectly
			targetPath := "/test/path"
			dependencies := []entities.Dependency{}
			cmd.Execute(targetPath, tc.arguments, dependencies)

			// THEN: Verify correct repository was used and value was set
			if tc.shouldUseInteractive {
				assert.Equal(t, 0, repository.ExecuteCallCount, "Should use interactive repository")
				assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
				assert.Equal(t, tc.expected, interactiveRepository.GetAutoAnswerValue(), "Should set correct auto-answer value")
			} else {
				assert.Equal(t, 1, repository.ExecuteCallCount, "Should use normal repository")
				assert.Equal(t, 0, interactiveRepository.ExecuteCallCount, "Should not use interactive repository")
			}
		})
	}
}

func TestRunFromRootCommand_RemoveAutoAnswerFlagEnhanced(t *testing.T) {
	t.Parallel()

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

	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "should remove --auto-answer=y flag",
			input:    []string{"plan", "--auto-answer=y", "--detailed-exitcode"},
			expected: []string{"plan", "--detailed-exitcode"},
		},
		{
			name:     "should remove --auto-answer=n flag",
			input:    []string{"apply", "--auto-answer=n", "--auto-approve"},
			expected: []string{"apply", "--auto-approve"},
		},
		{
			name:     "should remove --auto-answer=yes flag",
			input:    []string{"plan", "--auto-answer=yes"},
			expected: []string{"plan"},
		},
		{
			name:     "should remove --auto-answer=no flag",
			input:    []string{"apply", "--auto-answer=no"},
			expected: []string{"apply"},
		},
		{
			name:     "should remove --auto-answer=invalid flag",
			input:    []string{"plan", "--auto-answer=invalid", "--out=plan.out"},
			expected: []string{"plan", "--out=plan.out"},
		},
		{
			name:     "should remove multiple auto-answer flags",
			input:    []string{"plan", "--auto-answer=y", "-a", "--auto-answer=n", "--out=plan.out"},
			expected: []string{"plan", "--out=plan.out"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset state
			interactiveRepository.ExecuteCallCount = 0
			interactiveRepository.LastArguments = nil

			// WHEN: Executing the command to test filtering
			targetPath := "/test/path"
			dependencies := []entities.Dependency{}
			cmd.Execute(targetPath, tc.input, dependencies)

			// THEN: Should filter arguments correctly
			require.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
			assert.Equal(t, tc.expected, interactiveRepository.LastArguments, "Should filter out auto-answer flags")
		})
	}
}