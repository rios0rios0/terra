//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Run("should execute all sub-commands and terragrunt command", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{Name: "Go", CLI: "go", VersionURL: "https://golang.org/dl/", RegexVersion: `go(\d+\.\d+\.\d+)`, BinaryURL: "https://golang.org/dl/go%s.linux-amd64.tar.gz"},
		}
		installCommand := NewInstallDependenciesCommandStub().WithSuccess()
		formatCommand := NewFormatFilesCommandStub().WithSuccess()
		additionalBefore := NewRunAdditionalBeforeCommandStub().WithSuccess()
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBefore, repository)

		// when
		command.Execute("target/path", []string{"apply"}, dependencies)

		// then
		assert.True(t, installCommand.Executed, "install command should be executed")
		assert.True(t, formatCommand.Executed, "format command should be executed")
		assert.True(t, additionalBefore.Executed, "additional before command should be executed")
		assert.True(t, repository.CommandExecuted("terragrunt", []string{"apply"}), "terragrunt command should be executed")
	})

	t.Run("should handle errors gracefully", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{Name: "Go", CLI: "go", VersionURL: "https://golang.org/dl/", RegexVersion: `go(\d+\.\d+\.\d+)`, BinaryURL: "https://golang.org/dl/go%s.linux-amd64.tar.gz"},
		}
		installCommand := NewInstallDependenciesCommandStub().WithError()
		formatCommand := NewFormatFilesCommandStub().WithError()
		additionalBefore := NewRunAdditionalBeforeCommandStub().WithError()
		repository := repositories.NewShellRepositoryStub().WithError()
		command := commands.NewRunFromRootCommand(installCommand, formatCommand, additionalBefore, repository)

		// when
		command.Execute("target/path", []string{"apply"}, dependencies)

		// then
		assert.False(t, installCommand.Executed, "install command should not be executed")
		assert.False(t, formatCommand.Executed, "format command should not be executed")
		assert.False(t, additionalBefore.Executed, "additional before command should not be executed")
		assert.False(t, repository.CommandExecuted("terragrunt", []string{"apply"}), "terragrunt command should not be executed")
	})
}

// Mock command stubs for testing
type InstallDependenciesCommandStub struct {
	Executed bool
}

func NewInstallDependenciesCommandStub() *InstallDependenciesCommandStub {
	return &InstallDependenciesCommandStub{}
}

func (cmd *InstallDependenciesCommandStub) WithSuccess() *InstallDependenciesCommandStub {
	cmd.Executed = true
	return cmd
}

func (cmd *InstallDependenciesCommandStub) WithError() *InstallDependenciesCommandStub {
	cmd.Executed = false
	return cmd
}

func (cmd *InstallDependenciesCommandStub) Execute(dependencies []entities.Dependency) {
	cmd.Executed = true
}

type FormatFilesCommandStub struct {
	Executed bool
}

func NewFormatFilesCommandStub() *FormatFilesCommandStub {
	return &FormatFilesCommandStub{}
}

func (cmd *FormatFilesCommandStub) WithSuccess() *FormatFilesCommandStub {
	cmd.Executed = true
	return cmd
}

func (cmd *FormatFilesCommandStub) WithError() *FormatFilesCommandStub {
	cmd.Executed = false
	return cmd
}

func (cmd *FormatFilesCommandStub) Execute(dependencies []entities.Dependency) {
	cmd.Executed = true
}

type RunAdditionalBeforeCommandStub struct {
	Executed bool
}

func NewRunAdditionalBeforeCommandStub() *RunAdditionalBeforeCommandStub {
	return &RunAdditionalBeforeCommandStub{}
}

func (cmd *RunAdditionalBeforeCommandStub) WithSuccess() *RunAdditionalBeforeCommandStub {
	cmd.Executed = true
	return cmd
}

func (cmd *RunAdditionalBeforeCommandStub) WithError() *RunAdditionalBeforeCommandStub {
	cmd.Executed = false
	return cmd
}

func (cmd *RunAdditionalBeforeCommandStub) Execute(targetPath string, arguments []string) {
	cmd.Executed = true
}
