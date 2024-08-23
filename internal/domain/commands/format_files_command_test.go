//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/test/domain/entities/builders"
	"github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFilesCommand_Execute(t *testing.T) {
	t.Run("should format files for each dependency", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewOSRepositoryStub().WithSuccess()
		command := commands.NewFormatFilesCommand(repository)

		listeners := interfaces.FormatFilesListeners{
			OnSuccess: func() {
				// then
				assert.True(t, true, "the success listener should be called")
			},
		}

		// when
		command.Execute(dependencies, listeners)
	})

	t.Run("should throw an error when some unexpected condition happens", func(t *testing.T) {
		// given
		dependencies := builders.NewDependencyBuilder().BuildMany()
		repository := doubles.NewOSRepositoryStub().WithError()
		command := commands.NewFormatFilesCommand(repository)

		listeners := interfaces.FormatFilesListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called")
				assert.ErrorContains(t, err, "failed to format")
			},
		}

		// when
		command.Execute(dependencies, listeners)
	})
}
