//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	testrepositories "github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAdditionalBeforeCommand_Execute(t *testing.T) {
	t.Run("should perform all actions before executing the next command successfully", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		repository := testrepositories.NewOSRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{OnSuccess: func() {
			// then
			assert.True(t, true, "the success listener should be called")
		}}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})

	t.Run("should return an error when changing account fails", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		repository := testrepositories.NewOSRepositoryStub().WithError()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called when changing account fails")
				assert.ErrorContains(t, err, "error changing account")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})

	t.Run("should return an error when initializing the environment fails", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		repository := testrepositories.NewOSRepositoryStub().WithError()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called when initializing the environment fails")
				assert.ErrorContains(t, err, "error initializing the environment")
			},
		}

		// when
		command.Execute("target/path", []string{"apply", "init"}, listeners)
	})

	t.Run("should return an error when changing workspace fails", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "new-workspace"}
		repository := testrepositories.NewOSRepositoryStub().WithError()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called when changing workspace fails")
				assert.ErrorContains(t, err, "error changing workspace")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})
}
