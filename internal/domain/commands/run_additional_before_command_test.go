//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunAdditionalBeforeCommand_Execute(t *testing.T) {
	t.Run("should change account if CLI can change account", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		cli := entities.NewCLIStub().WithCanChangeAccount(true)
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.CommandExecuted(cli.GetName(), cli.GetCommandChangeAccount()), "command should be executed to change account")
	})

	t.Run("should init environment if necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "default"}
		cli := entities.NewCLIStub().WithCanChangeAccount(false)
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.CommandExecuted("terragrunt", []string{"init"}), "command should be executed to init environment")
	})

	t.Run("should change workspace if necessary", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraTerraformWorkspace: "new-workspace"}
		cli := entities.NewCLIStub().WithCanChangeAccount(false)
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// when
		command.Execute("target/path", []string{"apply"})

		// then
		assert.True(t, repository.CommandExecuted("terragrunt", []string{"workspace", "select", "-or-create", "new-workspace"}), "command should be executed to change workspace")
	})
}

// Mock functions for testing
func shouldInitEnvironment(arguments []string) bool {
	// Mock implementation for testing
	return true
}
