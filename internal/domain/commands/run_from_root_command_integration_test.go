//go:build integration

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

func TestRunFromRootCommand_AutoAnswer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("should use different repositories based on auto-answer flag", func(t *testing.T) {
		// GIVEN: A command with both regular and interactive repositories
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
		dependencies := []entities.Dependency{}

		// Test case 1: Without auto-answer flag
		argumentsWithoutFlag := []string{"plan", "--detailed-exitcode"}
		cmd.Execute(targetPath, argumentsWithoutFlag, dependencies)
		
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should use regular repository without auto-answer flag")
		assert.Equal(t, 0, interactiveRepository.ExecuteCallCount, "Should not use interactive repository without auto-answer flag")
		
		// Reset counters
		repository.ExecuteCallCount = 0
		interactiveRepository.ExecuteCallCount = 0

		// Test case 2: With --auto-answer flag
		argumentsWithFlag := []string{"plan", "--auto-answer", "--detailed-exitcode"}
		cmd.Execute(targetPath, argumentsWithFlag, dependencies)
		
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use regular repository with auto-answer flag")
		assert.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository with auto-answer flag")
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
	})

	t.Run("should properly filter auto-answer flags from arguments", func(t *testing.T) {
		// GIVEN: A command with interactive repository stub to capture arguments
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
				name:     "should filter --auto-answer flag",
				input:    []string{"plan", "--auto-answer", "--detailed-exitcode"},
				expected: []string{"plan", "--detailed-exitcode"},
			},
			{
				name:     "should filter -a flag",
				input:    []string{"apply", "-a", "--auto-approve"},
				expected: []string{"apply", "--auto-approve"},
			},
			{
				name:     "should filter both flags",
				input:    []string{"plan", "--auto-answer", "-a", "--out=plan.out"},
				expected: []string{"plan", "--out=plan.out"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Reset counters
				interactiveRepository.ExecuteCallCount = 0
				interactiveRepository.LastArguments = nil

				// WHEN: Executing with auto-answer flag
				cmd.Execute("/test/path", tc.input, []entities.Dependency{})

				// THEN: Should use interactive repository with filtered arguments
				require.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
				assert.Equal(t, tc.expected, interactiveRepository.LastArguments, "Should filter out auto-answer flags")
			})
		}
	})
}

func TestInteractiveShellRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("should create interactive repository successfully", func(t *testing.T) {
		// GIVEN: No preconditions
		
		// WHEN: Creating an interactive shell repository
		repo := infrastructure_repositories.NewInteractiveShellRepository()
		
		// THEN: Should return a valid instance
		require.NotNil(t, repo, "Should create repository instance")
	})

	t.Run("should handle simple echo command without hanging", func(t *testing.T) {
		// GIVEN: An interactive shell repository and a simple command
		repo := infrastructure_repositories.NewInteractiveShellRepository()
		
		// WHEN: Executing a simple echo command
		err := repo.ExecuteCommand("echo", []string{"test"}, ".")
		
		// THEN: Should execute without error
		assert.NoError(t, err, "Should execute simple command without error")
	})

	// Note: Testing actual pattern matching would require complex setup with mock processes
	// that send specific prompts. For now, we verify the repository can be created and
	// execute simple commands without hanging or crashing.
}

func TestRunFromRootCommand_AutoAnswerEnhanced_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("should configure interactive repository with different auto-answer values", func(t *testing.T) {
		// GIVEN: Commands with different auto-answer values
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
		dependencies := []entities.Dependency{}

		testCases := []struct {
			name      string
			arguments []string
			expected  string
		}{
			{
				name:      "should configure 'y' for --auto-answer=y",
				arguments: []string{"apply", "--auto-answer=y"},
				expected:  "y",
			},
			{
				name:      "should configure 'n' for --auto-answer=n",
				arguments: []string{"plan", "--auto-answer=n"},
				expected:  "n",
			},
			{
				name:      "should configure 'y' for --auto-answer=yes",
				arguments: []string{"apply", "--auto-answer=yes"},
				expected:  "y",
			},
			{
				name:      "should configure 'n' for --auto-answer=no",
				arguments: []string{"plan", "--auto-answer=no"},
				expected:  "n",
			},
			{
				name:      "should configure 'n' for legacy --auto-answer",
				arguments: []string{"plan", "--auto-answer"},
				expected:  "n",
			},
			{
				name:      "should configure 'n' for legacy -a",
				arguments: []string{"apply", "-a"},
				expected:  "n",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Reset state
				interactiveRepository.ExecuteCallCount = 0

				// WHEN: Executing with the specific auto-answer value
				cmd.Execute(targetPath, tc.arguments, dependencies)

				// THEN: Should configure the repository correctly
				require.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
				assert.Equal(t, tc.expected, interactiveRepository.GetAutoAnswerValue(), "Should configure correct auto-answer value")
			})
		}
	})

	t.Run("should properly filter enhanced auto-answer flags", func(t *testing.T) {
		// GIVEN: A command with enhanced auto-answer flags
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
			autoAnswer string
		}{
			{
				name:     "should filter --auto-answer=y and preserve other flags",
				input:    []string{"apply", "--auto-answer=y", "--auto-approve", "--force"},
				expected: []string{"apply", "--auto-approve", "--force"},
				autoAnswer: "y",
			},
			{
				name:     "should filter --auto-answer=n and preserve other flags",
				input:    []string{"plan", "--auto-answer=n", "--detailed-exitcode", "--out=plan.out"},
				expected: []string{"plan", "--detailed-exitcode", "--out=plan.out"},
				autoAnswer: "n",
			},
			{
				name:     "should filter multiple auto-answer flags",
				input:    []string{"plan", "--auto-answer=y", "-a", "--auto-answer=n", "--detailed-exitcode"},
				expected: []string{"plan", "--detailed-exitcode"},
				autoAnswer: "y", // First match wins
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Reset state
				interactiveRepository.ExecuteCallCount = 0
				interactiveRepository.LastArguments = nil

				// WHEN: Executing with enhanced auto-answer flags
				cmd.Execute("/test/path", tc.input, []entities.Dependency{})

				// THEN: Should filter arguments and configure auto-answer value
				require.Equal(t, 1, interactiveRepository.ExecuteCallCount, "Should use interactive repository")
				assert.Equal(t, tc.expected, interactiveRepository.LastArguments, "Should filter out auto-answer flags")
				assert.Equal(t, tc.autoAnswer, interactiveRepository.GetAutoAnswerValue(), "Should configure correct auto-answer value")
			})
		}
	})
}