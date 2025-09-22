package commands

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

// RunFromRootCommandDeps contains dependencies for RunFromRootCommand
type RunFromRootCommandDeps struct {
	InstallCommand        InstallDependencies
	FormatCommand         FormatFiles
	AdditionalBefore      RunAdditionalBefore
	Repository            repositories.ShellRepository
	InteractiveRepository repositories.ShellRepository
}

type RunFromRootCommand struct {
	installCommand        InstallDependencies
	formatCommand         FormatFiles
	additionalBefore      RunAdditionalBefore
	repository            repositories.ShellRepository
	interactiveRepository repositories.ShellRepository
}

func NewRunFromRootCommand(deps RunFromRootCommandDeps) *RunFromRootCommand {
	return &RunFromRootCommand{
		installCommand:        deps.InstallCommand,
		formatCommand:         deps.FormatCommand,
		additionalBefore:      deps.AdditionalBefore,
		repository:            deps.Repository,
		interactiveRepository: deps.InteractiveRepository,
	}
}

func (it *RunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
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
