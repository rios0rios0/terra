//go:build unit

package commands_test

import (
	"errors"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/domain/entitydoubles"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunAdditionalBeforeCommand_Execute_AccountChangeError(t *testing.T) {
	t.Run("should fatalf when account change command fails", func(t *testing.T) {
		// GIVEN: A command with CLI that can change account but the command fails
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		cli := &entitydoubles.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role"},
		}
		repository := &repositorydoubles.StubShellRepositoryForAdditional{
			ExecuteErrors: []error{errors.New("account change failed")},
		}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should log a fatal error about the account change failure
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "Error changing account")
	})
}

func TestRunAdditionalBeforeCommand_Execute_WorkspaceChangeError(t *testing.T) {
	t.Run("should fatalf when workspace change command fails", func(t *testing.T) {
		// GIVEN: A command with workspace configured but the workspace command fails
		hook, cleanup := setupFatalInterceptor()
		defer cleanup()

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("production").
			BuildSettings()
		// Use a repository that succeeds for init but fails for workspace change
		// The call order is: init (succeeds or not called), workspace (fails)
		repository := &repositorydoubles.StubShellRepositoryForAdditional{
			// First call: init (success), Second call: workspace (fail)
			ExecuteErrors: []error{nil, errors.New("workspace change failed")},
		}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		t.Setenv("TG_DOWNLOAD_DIR", "")
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should log a fatal error about the workspace change failure
		require.NotEmpty(t, hook.Entries)
		lastEntry := hook.LastEntry()
		assert.Equal(t, logger.FatalLevel, lastEntry.Level)
		assert.Contains(t, lastEntry.Message, "Error changing workspace")
	})
}

func TestRunAdditionalBeforeCommand_Execute_StateCommand(t *testing.T) {
	t.Run("should not init environment when command is state manipulation", func(t *testing.T) {
		// GIVEN: A command with a state manipulation argument
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		t.Setenv("TG_DOWNLOAD_DIR", "")
		arguments := []string{"import", "null_resource.test", "test-id"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init for state manipulation commands
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init for state manipulation commands")
			}
		}
	})
}
