package commands

import (
	"strings"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"
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

func (it *RunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	it.formatCommand.Execute(dependencies)
	it.additionalBefore.Execute(targetPath, arguments)

	// Check if auto-answer flag is present and extract its value
	autoAnswerValue := it.getAutoAnswerValue(arguments)
	useInteractive := autoAnswerValue != ""
	filteredArguments := it.removeAutoAnswerFlag(arguments)

	var err error
	if useInteractive {
		logger.Infof("Using interactive mode with auto-answering (response: %s)", autoAnswerValue)
		
		// Configure the interactive repository with the auto-answer value
		it.interactiveRepository.SetAutoAnswerValue(autoAnswerValue)
		
		err = it.interactiveRepository.ExecuteCommand("terragrunt", filteredArguments, targetPath)
	} else {
		err = it.repository.ExecuteCommand("terragrunt", filteredArguments, targetPath)
	}

	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

// getAutoAnswerValue extracts the auto-answer value from arguments
// Returns "n" for --auto-answer or -a (default), "y" for --auto-answer=y, "n" for --auto-answer=n
// Returns empty string if no auto-answer flag is present
func (it *RunFromRootCommand) getAutoAnswerValue(arguments []string) string {
	for _, arg := range arguments {
		if arg == "--auto-answer" || arg == "-a" {
			return "n" // Default to "n" for backward compatibility
		}
		if strings.HasPrefix(arg, "--auto-answer=") {
			value := strings.TrimPrefix(arg, "--auto-answer=")
			if value == "y" || value == "yes" {
				return "y"
			}
			if value == "n" || value == "no" {
				return "n"
			}
			// Invalid value, default to "n"
			logger.Warnf("Invalid auto-answer value '%s', defaulting to 'n'", value)
			return "n"
		}
	}
	return ""
}

func (it *RunFromRootCommand) hasAutoAnswerFlag(arguments []string) bool {
	return it.getAutoAnswerValue(arguments) != ""
}

func (it *RunFromRootCommand) removeAutoAnswerFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != "--auto-answer" && arg != "-a" && !strings.HasPrefix(arg, "--auto-answer=") {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
