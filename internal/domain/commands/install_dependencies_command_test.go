//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/test/domain/entities/builders"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Run("should install dependencies if not available", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		command := commands.NewInstallDependenciesCommand()

		listeners := interfaces.InstallDependenciesListeners{OnSuccess: func() {
			// then
			assert.True(t, true, "the success listener should be called")
		}}

		// when
		command.Execute(dependencies, listeners)
	})

	t.Run("should throw an error when some unexpected condition happens", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		command := commands.NewInstallDependenciesCommand()

		listeners := interfaces.InstallDependenciesListeners{OnError: func(err error) {
			// then
			assert.Error(t, err, "the error listener should be called")
		}}

		// when
		command.Execute(dependencies, listeners)
	})
}
