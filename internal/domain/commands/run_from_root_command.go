package commands

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"
	logger "github.com/sirupsen/logrus"
)

type RunFromRootCommand struct {
	installCommand        InstallDependencies
	formatCommand         FormatFiles
	additionalBefore      RunAdditionalBefore
	repository            repositories.ShellRepository
	interactiveRepository *infrastructure_repositories.InteractiveShellRepository
}

func NewRunFromRootCommand(
	installCommand InstallDependencies,
	formatCommand FormatFiles,
	additionalBefore RunAdditionalBefore,
	repository repositories.ShellRepository,
	interactiveRepository *infrastructure_repositories.InteractiveShellRepository,
) *RunFromRootCommand {
	return &RunFromRootCommand{
		installCommand:        installCommand,
		formatCommand:         formatCommand,
		additionalBefore:      additionalBefore,
		repository:            repository,
		interactiveRepository: interactiveRepository,
	}
}

func (it *RunFromRootCommand) Execute(targetPath string, arguments []string, dependencies []entities.Dependency) {
	// ensure that all dependencies are installed
	it.installCommand.Execute(dependencies)
	it.formatCommand.Execute(dependencies)
	it.additionalBefore.Execute(targetPath, arguments)

	// Check if auto-answer flag is present and filter it out
	useInteractive := it.hasAutoAnswerFlag(arguments)
	filteredArguments := it.removeAutoAnswerFlag(arguments)

	var err error
	if useInteractive {
		logger.Info("Using interactive mode with auto-answering")
		err = it.interactiveRepository.ExecuteCommand("terragrunt", filteredArguments, targetPath)
	} else {
		err = it.repository.ExecuteCommand("terragrunt", filteredArguments, targetPath)
	}

	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

func (it *RunFromRootCommand) hasAutoAnswerFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == "--auto-answer" || arg == "-a" {
			return true
		}
	}
	return false
}

func (it *RunFromRootCommand) removeAutoAnswerFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != "--auto-answer" && arg != "-a" {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
