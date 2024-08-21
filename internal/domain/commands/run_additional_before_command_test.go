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
	t.Run("should change account when CLI can change the account", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{OnSuccess: func() {
			// then
			assert.True(t, true, "the success listener should be called")
		}}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})

	t.Run("should initialize the environment when it's necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{OnSuccess: func() {
			// then
			assert.True(t, true, "the success listener should be called")
		}}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})

	t.Run("should change the workspace when it's necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "new-workspace"}
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{OnSuccess: func() {
			// then
			assert.True(t, true, "the success listener should be called")
		}}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})

	t.Run("should throw an error when some unexpected condition happens", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "new-workspace"}
		repository := testrepositories.NewShellRepositoryStub().WithError()
		command := commands.NewRunAdditionalBeforeCommand(settings, repository)

		listeners := interfaces.RunAdditionalBeforeListeners{OnError: func(err error) {
			// then
			assert.Error(t, err, "the error listener should be called")
		}}

		// when
		command.Execute("target/path", []string{"apply"}, listeners)
	})
}
