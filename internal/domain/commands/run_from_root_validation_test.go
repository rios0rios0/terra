//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"
	logger "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupFatalInterceptor configures logrus to capture Fatal-level log entries instead
// of calling os.Exit(1). It returns the test hook containing captured entries and a
// cleanup function that must be deferred.
func setupFatalInterceptor() (*test.Hook, func()) {
	hook := test.NewLocal(logger.StandardLogger())
	originalExitFunc := logger.StandardLogger().ExitFunc
	logger.StandardLogger().ExitFunc = func(int) {}
	cleanup := func() {
		logger.StandardLogger().ExitFunc = originalExitFunc
		hook.Reset()
	}
	return hook, cleanup
}

// newRunFromRootForValidation creates a RunFromRootCommand with stub dependencies
// suitable for testing validation paths.
func newRunFromRootForValidation() *commands.RunFromRootCommand {
	return commands.NewRunFromRootCommand(
		entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir("/tmp/terra-test-modules").
			WithTerraProviderCacheDir("/tmp/terra-test-providers").
			BuildSettings(),
		&commanddoubles.StubInstallDependencies{},
		&commanddoubles.StubFormatFiles{},
		&commanddoubles.StubRunAdditionalBefore{},
		&commanddoubles.StubParallelState{},
		&repositorydoubles.StubShellRepositoryForRoot{},
		&repositorydoubles.StubUpgradeShellRepository{},
		&repositorydoubles.StubInteractiveShellRepository{},
	)
}

func TestRunFromRootCommand_validateDeprecatedFlags(t *testing.T) {
	t.Run("should fatalf when -a short flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated -a short flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "-a=y"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the removed -a flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "-a short flag has been removed")
	})

	t.Run("should fatalf when -a boolean flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated -a boolean flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "-a"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the removed -a flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "-a short flag has been removed")
	})

	t.Run("should fatalf when --auto-answer flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated --auto-answer flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--auto-answer=y"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the renamed --auto-answer flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--auto-answer has been renamed to --reply")
	})

	t.Run("should fatalf when --auto-answer boolean flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated --auto-answer boolean flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--auto-answer"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the renamed --auto-answer flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--auto-answer has been renamed to --reply")
	})

	t.Run("should fatalf when --all is used with state commands", func(t *testing.T) {
		// GIVEN: Arguments containing --all with a state manipulation command
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"import", "--all", "null_resource.test", "test-id"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about --all with state commands
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--all cannot be used with state commands")
	})

	t.Run("should fatalf when --no-parallel-bypass flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the removed --no-parallel-bypass flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--no-parallel-bypass"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the removed --no-parallel-bypass flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--no-parallel-bypass has been removed")
	})

	t.Run("should fatalf when --include flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated --include flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--include=mod1,mod2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the renamed --include flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--include has been renamed to --only")
	})

	t.Run("should fatalf when --exclude flag is used", func(t *testing.T) {
		// GIVEN: Arguments containing the deprecated --exclude flag
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--exclude=mod1,mod2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the renamed --exclude flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--exclude has been renamed to --skip")
	})
}

func TestRunFromRootCommand_validateFlagCombinations(t *testing.T) {
	t.Run("should fatalf when --parallel and --all are used together", func(t *testing.T) {
		// GIVEN: Arguments containing both --parallel and --all
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--parallel=5", "--all"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about conflicting flags
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--parallel and --all cannot be used together")
	})

	t.Run("should fatalf when --parallel is used with apply without --reply", func(t *testing.T) {
		// GIVEN: Arguments containing --parallel with apply but no --reply
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"apply", "--parallel=2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about missing --reply flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--reply is required when using --parallel with apply or destroy")
	})

	t.Run("should fatalf when --parallel is used with destroy without --reply", func(t *testing.T) {
		// GIVEN: Arguments containing --parallel with destroy but no --reply
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"destroy", "--parallel=3"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about missing --reply flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--reply is required when using --parallel with apply or destroy")
	})
}

func TestRunFromRootCommand_validateSelectionFlags(t *testing.T) {
	t.Run("should fatalf when --only is used without --parallel", func(t *testing.T) {
		// GIVEN: Arguments containing --only without --parallel
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--only=mod1,mod2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about missing --parallel flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--only/--skip flags require --parallel=N")
	})

	t.Run("should fatalf when --skip is used without --parallel", func(t *testing.T) {
		// GIVEN: Arguments containing --skip without --parallel
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--skip=mod1"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about missing --parallel flag
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--only/--skip flags require --parallel=N")
	})
}

func TestRunFromRootCommand_validateSelectionFlagValues(t *testing.T) {
	t.Run("should fatalf when --only flag has empty value", func(t *testing.T) {
		// GIVEN: Arguments containing --only= with empty value and --parallel
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--parallel=2", "--only="}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about empty --only values
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--only flag is present but has no values")
	})

	t.Run("should fatalf when --skip flag has empty value", func(t *testing.T) {
		// GIVEN: Arguments containing --skip= with empty value and --parallel
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--parallel=2", "--skip="}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about empty --skip values
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "--skip flag is present but has no values")
	})
}

func TestRunFromRootCommand_validateSelectionFlagConflicts(t *testing.T) {
	t.Run("should fatalf when same module appears in both --only and --skip", func(t *testing.T) {
		// GIVEN: Arguments containing the same module in both --only and --skip
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()
		cmd := newRunFromRootForValidation()
		arguments := []string{"plan", "--parallel=2", "--only=mod1,mod2", "--skip=mod2,mod3"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about conflicting module
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "appears in both --only and --skip")
		assert.Contains(t, lastEntry.Message, "mod2")
	})

	t.Run("should not fatalf when --only and --skip have no overlapping modules", func(t *testing.T) {
		// GIVEN: Arguments with non-overlapping --only and --skip values
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		parallelState := &commanddoubles.StubParallelState{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			parallelState,
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)
		arguments := []string{"plan", "--parallel=2", "--only=mod1,mod2", "--skip=mod3"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should not log any fatal error (validation passes, parallel command executes)
		for _, entry := range hook.Entries {
			assert.NotEqual(t, logger.FatalLevel, entry.Level,
				"Should not produce a fatal log entry when --only and --skip do not overlap")
		}
		assert.True(t, parallelState.ExecuteCalled, "Should proceed to parallel execution")
	})

	t.Run("should not fatalf when only --only is used without --skip", func(t *testing.T) {
		// GIVEN: Arguments with only --only flag (no --skip)
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		parallelState := &commanddoubles.StubParallelState{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			parallelState,
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)
		arguments := []string{"plan", "--parallel=2", "--only=mod1"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should not log any fatal error
		for _, entry := range hook.Entries {
			assert.NotEqual(t, logger.FatalLevel, entry.Level,
				"Should not produce a fatal log entry for --only without --skip")
		}
		assert.True(t, parallelState.ExecuteCalled, "Should proceed to parallel execution")
	})
}

func TestRunFromRootCommand_validateDeprecatedFlags_stateUtils(t *testing.T) {
	t.Run("should not fatalf when --all is used with non-state command", func(t *testing.T) {
		// GIVEN: Arguments containing --all with a non-state command (plan)
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			upgradeRepository,
			&repositorydoubles.StubInteractiveShellRepository{},
		)
		arguments := []string{"plan", "--all"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should not log any fatal error (--all with non-state commands is valid)
		for _, entry := range hook.Entries {
			assert.NotEqual(t, logger.FatalLevel, entry.Level,
				"Should not produce a fatal log entry for --all with non-state command")
		}
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should proceed to normal execution")
	})
}

func TestRunFromRootCommand_Execute_parallelStateFails(t *testing.T) {
	t.Run("should fatalf when parallel state command returns error", func(t *testing.T) {
		// GIVEN: A parallel state stub that returns an error
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		parallelState := &commanddoubles.StubParallelState{
			ShouldReturnError: true,
			ErrorMessage:      "simulated parallel failure",
		}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			parallelState,
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the parallel command failure
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "Parallel command failed")
	})
}

func TestRunFromRootCommand_Execute_terragruntFails(t *testing.T) {
	t.Run("should fatalf when upgrade-aware repository returns error", func(t *testing.T) {
		// GIVEN: An upgrade repository that returns an error
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{
			ErrorToReturn: assert.AnError,
		}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			upgradeRepository,
			&repositorydoubles.StubInteractiveShellRepository{},
		)
		arguments := []string{"plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the terragrunt command failure
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "Terragrunt command failed")
	})

	t.Run("should fatalf when interactive repository returns error", func(t *testing.T) {
		// GIVEN: An interactive repository that returns an error
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{
			ShouldReturnError: true,
			ExecuteError:      assert.AnError,
		}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().
				WithTerraModuleCacheDir("/tmp/terra-test-modules").
				WithTerraProviderCacheDir("/tmp/terra-test-providers").
				BuildSettings(),
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			interactiveRepository,
		)
		arguments := []string{"--reply=y", "plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute("/test/path", arguments, dependencies)

		// THEN: Should log a fatal error about the terragrunt command failure
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "Terragrunt command failed")
	})
}
