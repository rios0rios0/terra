package commands

import (
	"strings"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

const (
	// AutoAnswerFlag represents the --auto-answer flag.
	AutoAnswerFlag = "--auto-answer"
	// AutoAnswerShortFlag represents the -a flag.
	AutoAnswerShortFlag = "-a"
)

type RunFromRootCommand struct {
	installCommand        InstallDependencies
	formatCommand         FormatFiles
	additionalBefore      RunAdditionalBefore
	parallelState         ParallelState
	repository            repositories.ShellRepository
	interactiveRepository repositories.InteractiveShellRepository
}

func NewRunFromRootCommand(
	installCommand InstallDependencies,
	formatCommand FormatFiles,
	additionalBefore RunAdditionalBefore,
	parallelState ParallelState,
	repository repositories.ShellRepository,
	interactiveRepository repositories.InteractiveShellRepository,
) *RunFromRootCommand {
	return &RunFromRootCommand{
		installCommand:        installCommand,
		formatCommand:         formatCommand,
		additionalBefore:      additionalBefore,
		parallelState:         parallelState,
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

	// Check if this is a parallel state manipulation command
	if it.isParallelStateCommand(arguments) {
		// For parallel state commands, skip additional before steps as they don't make sense
		// when running across multiple directories
		err := it.parallelState.Execute(targetPath, arguments, dependencies)
		if err != nil {
			logger.Fatalf("Parallel state command failed: %s", err)
		}
		return
	}

	// Normal execution path for non-parallel commands
	it.additionalBefore.Execute(targetPath, arguments)

	// Check if auto-answer flag is present and filter it out
	useInteractive := it.hasAutoAnswerFlag(arguments)
	autoAnswerValue := it.getAutoAnswerValue(arguments)
	filteredArguments := it.removeAutoAnswerFlag(arguments)

	var err error
	if useInteractive {
		logger.Infof("Using interactive mode with auto-answering (%s)", autoAnswerValue)
		err = it.interactiveRepository.ExecuteCommandWithAnswer(
			"terragrunt", filteredArguments, targetPath, autoAnswerValue)
	} else {
		err = it.repository.ExecuteCommand("terragrunt", filteredArguments, targetPath)
	}

	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

func (it *RunFromRootCommand) hasAutoAnswerFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == AutoAnswerFlag || arg == AutoAnswerShortFlag ||
			strings.HasPrefix(arg, AutoAnswerFlag+"=") ||
			strings.HasPrefix(arg, AutoAnswerShortFlag+"=") {
			return true
		}
	}
	return false
}

func (it *RunFromRootCommand) getAutoAnswerValue(arguments []string) string {
	for _, arg := range arguments {
		if arg == AutoAnswerFlag || arg == AutoAnswerShortFlag {
			return "n" // Default backward compatibility behavior
		}
		if strings.HasPrefix(arg, AutoAnswerFlag+"=") {
			return arg[len(AutoAnswerFlag+"="):]
		}
		if strings.HasPrefix(arg, AutoAnswerShortFlag+"=") {
			return arg[len(AutoAnswerShortFlag+"="):]
		}
	}
	return ""
}

func (it *RunFromRootCommand) removeAutoAnswerFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != AutoAnswerFlag && arg != AutoAnswerShortFlag &&
			!strings.HasPrefix(arg, AutoAnswerFlag+"=") &&
			!strings.HasPrefix(arg, AutoAnswerShortFlag+"=") {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// isParallelStateCommand checks if the command should be executed in parallel.
func (it *RunFromRootCommand) isParallelStateCommand(arguments []string) bool {
	return HasAllFlag(arguments) && IsStateManipulationCommand(arguments)
}

// HasAutoAnswerFlagPublic is a public wrapper for testing the private hasAutoAnswerFlag method.
func (it *RunFromRootCommand) HasAutoAnswerFlagPublic(arguments []string) bool {
	return it.hasAutoAnswerFlag(arguments)
}

// GetAutoAnswerValuePublic is a public wrapper for testing the private getAutoAnswerValue method.
func (it *RunFromRootCommand) GetAutoAnswerValuePublic(arguments []string) string {
	return it.getAutoAnswerValue(arguments)
}

// RemoveAutoAnswerFlagPublic is a public wrapper for testing the private removeAutoAnswerFlag method.
func (it *RunFromRootCommand) RemoveAutoAnswerFlagPublic(arguments []string) []string {
	return it.removeAutoAnswerFlag(arguments)
}
