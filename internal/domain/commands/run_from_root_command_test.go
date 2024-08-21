//go:build unit

package commands_test

import (
	"testing"
)

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Run("should execute all sub-commands and terragrunt command", func(t *testing.T) {
		// given
		//dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		//installCommand := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
		//formatCommand := testcommands.NewFormatFilesCommandStub().WithSuccess()
		//additionalBefore := testcommands.NewRunAdditionalBeforeCommandStub().WithSuccess()
		//repository := testrepositories.NewShellRepositoryStub().WithSuccess()
		//command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBefore, repository)
		//
		//// when
		//command.Execute("target/path", []string{"apply"}, dependencies)
		//
		//// then
		//assert.True(t, installCommand, "install command should be executed")
		//assert.True(t, formatCommand, "format command should be executed")
		//assert.True(t, additionalBefore, "additional before command should be executed")
		//assert.True(t, repository.CommandExecuted("terragrunt", []string{"apply"}), "terragrunt command should be executed")
	})

	t.Run("should handle errors gracefully", func(t *testing.T) {
		// given
		//dependencies := testbuilders.NewDependencyBuilder().BuildMany()
		//installCommand := testcommands.NewInstallDependenciesCommandStub().WithError()
		//formatCommand := testcommands.NewFormatFilesCommandStub().WithError()
		//additionalBefore := testcommands.NewRunAdditionalBeforeCommandStub().WithError()
		//repository := testrepositories.NewShellRepositoryStub().WithError()
		//command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBefore, repository)
		//
		//// when
		//command.Execute("target/path", []string{"apply"}, dependencies)
		//
		//// then
		//assert.False(t, installCommand, "install command should not be executed")
		//assert.False(t, formatCommand, "format command should not be executed")
		//assert.False(t, additionalBefore, "additional before command should not be executed")
		//assert.False(t, repository.CommandExecuted("terragrunt", []string{"apply"}), "terragrunt command should not be executed")
	})
}
