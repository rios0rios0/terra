package commands

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
)

type RunFromRootCommand struct {
	installCommand   interfaces.InstallDependencies
	formatCommand    interfaces.FormatFiles
	additionalBefore interfaces.RunAdditionalBefore
	repository       repositories.OSRepository
}

func NewRunFromRootCommand(
	installCommand interfaces.InstallDependencies,
	formatCommand interfaces.FormatFiles,
	additionalBefore interfaces.RunAdditionalBefore,
	repository repositories.OSRepository,
) *RunFromRootCommand {
	return &RunFromRootCommand{
		installCommand:   installCommand,
		formatCommand:    formatCommand,
		additionalBefore: additionalBefore,
		repository:       repository,
	}
}

func (it *RunFromRootCommand) Execute(targetPath string, arguments []string, dependencies []entities.Dependency, listeners interfaces.RunFromRootListeners) {
	listeners3 := interfaces.RunAdditionalBeforeListeners{
		OnSuccess: func() {
			err := it.repository.ExecuteCommand("terragrunt", arguments, targetPath)
			if err != nil {
				listeners.OnError(fmt.Errorf("error happened while executing Terragrunt: %w", err))
				return
			}
			listeners.OnSuccess()
		},
		OnError: listeners.OnError,
	}

	listeners2 := interfaces.FormatFilesListeners{
		OnSuccess: func() {
			it.additionalBefore.Execute(targetPath, arguments, listeners3)
		},
		OnError: listeners.OnError,
	}

	// ensure that all dependencies are installed
	listeners1 := interfaces.InstallDependenciesListeners{
		OnSuccess: func() {
			it.formatCommand.Execute(dependencies, listeners2)
		},
		OnError: listeners.OnError,
	}
	it.installCommand.Execute(dependencies, listeners1)
}
