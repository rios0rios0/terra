//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/test/domain/entities/builders"
	"github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Run("should install dependencies if not available", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewShellRepositoryStub().WithSuccess()
		command := commands.NewInstallDependenciesCommand()

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.True(t, repository.ExecuteCommand(dependency.CLI, "install"), "command should be executed for: "+dependency.CLI)
		}
	})

	t.Run("should handle already installed dependencies gracefully", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewShellRepositoryStub().WithError()
		command := commands.NewInstallDependenciesCommand()

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.False(t, repository.ExecuteCommand(dependency.CLI, "install"), "command should not be executed for: "+dependency.CLI)
		}
	})
}
