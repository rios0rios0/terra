//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	testentities "github.com/rios0rios0/terra/test/domain/entities/doubles"
	testrepositories "github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAdditionalBeforeCommand_Execute(t *testing.T) {
	t.Run("should change account if CLI can change account", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		cli := testentities.NewCLIStub().WithCanChangeAccount(true)
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.ExecuteCommand(cli.GetName(), cli.GetCommandChangeAccount()), "command should be executed to change account")
	})

	t.Run("should init environment if necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		cli := testentities.NewCLIStub().WithCanChangeAccount(false)
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.ExecuteCommand("terragrunt", []string{"init"}), "command should be executed to init environment")
	})

	t.Run("should change workspace if necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "new-workspace"}
		cli := testentities.NewCLIStub().WithCanChangeAccount(false)
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.ExecuteCommand("terragrunt", []string{"workspace", "select", "-or-create", "new-workspace"}), "command should be executed to change workspace")
	})
}
