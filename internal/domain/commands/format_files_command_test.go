//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFilesCommand_Execute(t *testing.T) {
	t.Run("should format files for each dependency", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{CLI: "go", FormattingCommand: "fmt"},
			{CLI: "python", FormattingCommand: "black"},
		}
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewFormatFilesCommand(repository)

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.True(t, repository.CommandExecuted(dependency.CLI, dependency.FormattingCommand), "command should be executed for: "+dependency.CLI)
		}
	})

	t.Run("should handle errors gracefully", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{CLI: "go", FormattingCommand: "fmt"},
		}
		repository := repositories.NewShellRepositoryStub().WithError()
		command := commands.NewFormatFilesCommand(repository)

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.False(t, repository.CommandExecuted(dependency.CLI, dependency.FormattingCommand), "command should not be executed for: "+dependency.CLI)
		}
	})
}
