//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"
	"github.com/stretchr/testify/assert"
)

func TestRunFromRootCommand_ExecuteParallelState(t *testing.T) {
	t.Run("should execute parallel state command when import with --all", func(t *testing.T) {
		// GIVEN: A command with parallel state import arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			&repositorydoubles.StubUpgradeShellRepository{},
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"import", "--all", "null_resource.test", "test-id"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should NOT execute format command (state commands skip formatting)
		assert.False(t, formatCommand.ExecuteCalled, "Should not execute format command for state commands")

		// Should NOT execute additional before (skipped for parallel state)
		assert.False(t, additionalBefore.ExecuteCalled, "Should not execute additional before for parallel state")

		// Should execute parallel state command
		assert.True(t, parallelState.ExecuteCalled, "Should execute parallel state command")
		assert.Equal(t, targetPath, parallelState.LastTargetPath)
		assert.Equal(t, arguments, parallelState.LastArguments)

		// Should NOT execute normal terragrunt command
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not execute normal terragrunt command for parallel state")
	})

	t.Run("should execute parallel state command when state rm with --all", func(t *testing.T) {
		// GIVEN: A command with parallel state rm arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			&repositorydoubles.StubUpgradeShellRepository{},
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"state", "rm", "--all", "null_resource.test"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute parallel state command
		assert.True(t, parallelState.ExecuteCalled, "Should execute parallel state command")
		assert.Equal(t, targetPath, parallelState.LastTargetPath)
		assert.Equal(t, arguments, parallelState.LastArguments)

		// Should NOT execute normal flow
		assert.False(t, additionalBefore.ExecuteCalled, "Should not execute additional before for parallel state")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not execute normal terragrunt command for parallel state")
	})

	t.Run("should execute normal flow when import without --all", func(t *testing.T) {
		// GIVEN: A command with normal import arguments (no --all)
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"import", "null_resource.test", "test-id"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should NOT execute format command (state commands skip formatting)
		assert.False(t, formatCommand.ExecuteCalled, "Should not execute format command for state commands")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should execute normal terragrunt command")

		// Should NOT execute parallel state command
		assert.False(t, parallelState.ExecuteCalled, "Should not execute parallel state command for normal import")
	})

	t.Run("should execute normal flow when plan with --all", func(t *testing.T) {
		// GIVEN: A command with plan arguments and --all (not a state manipulation command)
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--all"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute normal flow
		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should execute normal terragrunt command")

		// Should NOT execute parallel state command
		assert.False(t, parallelState.ExecuteCalled, "Should not execute parallel state command for plan")
	})

	t.Run("should execute parallel command when --parallel=N used with non-state command", func(t *testing.T) {
		// GIVEN: A command with --parallel=2 flag on a plan command (not a state command)
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute parallel state command via isParallelCommand
		assert.True(t, parallelState.ExecuteCalled, "Should execute parallel command for --parallel=N")
		assert.Equal(t, targetPath, parallelState.LastTargetPath)
		assert.Equal(t, arguments, parallelState.LastArguments)

		// Should NOT execute normal flow
		assert.False(t, additionalBefore.ExecuteCalled, "Should not execute additional before for parallel command")
		assert.Equal(t, 0, upgradeRepository.ExecuteCallCount, "Should not execute normal terragrunt command for parallel")
	})

	t.Run("should execute normal flow when --parallel=N used with --no-parallel-bypass", func(t *testing.T) {
		// GIVEN: A command with both --parallel=2 and --no-parallel-bypass and --all flags
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--parallel=2", "--no-parallel-bypass", "--all"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should bypass terra's parallel execution and forward to terragrunt
		assert.False(t, parallelState.ExecuteCalled, "Should not execute parallel command when --no-parallel-bypass is present")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before in normal flow")
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should execute normal terragrunt command")
	})

	t.Run("should execute parallel state command when state mv with --all", func(t *testing.T) {
		// GIVEN: A command with state mv arguments and --all flag
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		parallelState := &commanddoubles.StubParallelState{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"state", "mv", "--all", "old_resource", "new_resource"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute parallel state command (backward compat: state + --all)
		assert.True(t, parallelState.ExecuteCalled, "Should execute parallel state command")
		assert.Equal(t, targetPath, parallelState.LastTargetPath)
		assert.Equal(t, arguments, parallelState.LastArguments)

		// Should NOT execute normal flow
		assert.False(t, additionalBefore.ExecuteCalled, "Should not execute additional before for parallel state")
		assert.Equal(t, 0, upgradeRepository.ExecuteCallCount, "Should not execute normal terragrunt command")
	})
}
