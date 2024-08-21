//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/test/domain/entities/builders"
	"github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFilesCommand_Execute(t *testing.T) {
	t.Run("should format files for each dependency", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewShellRepositoryStub().WithSuccess()
		command := commands.NewFormatFilesCommand(repository)

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.Nil(t, repository.ExecuteCommand(dependency.CLI, dependency.FormattingCommand, "."), "command should be executed for: "+dependency.CLI)
		}
	})

	t.Run("should handle errors gracefully", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewShellRepositoryStub().WithError()
		command := commands.NewFormatFilesCommand(repository)

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.Error(t, repository.ExecuteCommand(dependency.CLI, dependency.FormattingCommand, "."), "command should not be executed for: "+dependency.CLI)
		}
	})
}
