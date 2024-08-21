//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	testcommands "github.com/rios0rios0/terra/test/domain/commands/doubles"
	testbuilders "github.com/rios0rios0/terra/test/domain/entities/builders"
	testrepositories "github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Run("should execute all sub-commands when requested successfully", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		installCommand := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
		formatCommand := testcommands.NewFormatFilesCommandStub().WithSuccess()
		additionalBeforeCommand := testcommands.NewRunAdditionalBeforeCommandStub().WithSuccess()
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBeforeCommand, repository)

		listeners := interfaces.RunFromRootListeners{
			OnSuccess: func() {
				// then
				assert.True(t, true, "the success listener should be called")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, dependencies, listeners)
	})

	t.Run("should return error when installing dependencies fails", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		installCommand := testcommands.NewInstallDependenciesCommandStub().WithError()
		formatCommand := testcommands.NewFormatFilesCommandStub().WithSuccess()
		additionalBeforeCommand := testcommands.NewRunAdditionalBeforeCommandStub().WithSuccess()
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBeforeCommand, repository)

		listeners := interfaces.RunFromRootListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, dependencies, listeners)
	})

	t.Run("should return error when formating files fails", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		installCommand := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
		formatCommand := testcommands.NewFormatFilesCommandStub().WithError()
		additionalBeforeCommand := testcommands.NewRunAdditionalBeforeCommandStub().WithSuccess()
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBeforeCommand, repository)

		listeners := interfaces.RunFromRootListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, dependencies, listeners)
	})

	t.Run("should return error when running additional commands fails", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		installCommand := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
		formatCommand := testcommands.NewFormatFilesCommandStub().WithSuccess()
		additionalBeforeCommand := testcommands.NewRunAdditionalBeforeCommandStub().WithError()
		repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBeforeCommand, repository)

		listeners := interfaces.RunFromRootListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, dependencies, listeners)
	})

	t.Run("should return error when the main command execution fails", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		installCommand := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
		formatCommand := testcommands.NewFormatFilesCommandStub().WithSuccess()
		additionalBeforeCommand := testcommands.NewRunAdditionalBeforeCommandStub().WithSuccess()
		repository := testrepositories.NewShellRepositoryStub().WithError()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBeforeCommand, repository)

		listeners := interfaces.RunFromRootListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called")
			},
		}

		// when
		command.Execute("target/path", []string{"apply"}, dependencies, listeners)
	})
}
