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
			&entities.Settings{},
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
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
			&entities.Settings{},
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
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
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			&entities.Settings{},
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
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
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute normal terragrunt command")

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
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			&entities.Settings{},
			installCommand,
			formatCommand,
			additionalBefore,
			parallelState,
			repository,
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
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute normal terragrunt command")

		// Should NOT execute parallel state command
		assert.False(t, parallelState.ExecuteCalled, "Should not execute parallel state command for plan")
	})
}
