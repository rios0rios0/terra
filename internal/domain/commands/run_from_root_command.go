package commands

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

type RunFromRootCommand struct {
	installCommand   InstallDependencies
	formatCommand    FormatFiles
	additionalBefore RunAdditionalBefore
	repository       repositories.ShellRepository
}

func NewRunFromRootCommand(
	installCommand InstallDependencies,
	formatCommand FormatFiles,
	additionalBefore RunAdditionalBefore,
	repository repositories.ShellRepository,
) *RunFromRootCommand {
	return &RunFromRootCommand{
		installCommand:   installCommand,
		formatCommand:    formatCommand,
		additionalBefore: additionalBefore,
		repository:       repository,
	}
}

func (it *RunFromRootCommand) Execute(targetPath string, arguments []string, dependencies []entities.Dependency) {
	// ensure that all dependencies are installed
	it.installCommand.Execute(dependencies)
	it.formatCommand.Execute(dependencies)
	it.additionalBefore.Execute(targetPath, arguments)

	err := it.repository.ExecuteCommand("terragrunt", arguments, targetPath)
	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}
