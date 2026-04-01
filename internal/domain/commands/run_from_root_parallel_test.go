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

}
