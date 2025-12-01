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

	// Validate flag combinations before execution
	it.validateFlagCombinations(arguments)

	// Check if this is a parallel command (either state command with --all or any command with --parallel=N)
	if it.isParallelCommand(arguments) {
		// For parallel commands, skip additional before steps as they don't make sense
		// when running across multiple directories
		err := it.parallelState.Execute(targetPath, arguments, dependencies)
		if err != nil {
			logger.Fatalf("Parallel command failed: %s", err)
		}
		return
	}

	// Normal execution path for non-parallel commands
	it.additionalBefore.Execute(targetPath, arguments)

	// Check if auto-answer flag is present and filter it out
	useInteractive := it.hasAutoAnswerFlag(arguments)
	autoAnswerValue := it.getAutoAnswerValue(arguments)
	filteredArguments := it.removeAutoAnswerFlag(arguments)
	
	// Remove --no-parallel-bypass flag before passing to terragrunt
	filteredArguments = RemoveNoParallelBypassFlag(filteredArguments)

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

// validateFlagCombinations validates that flag combinations are correct.
// Errors and exits if invalid combinations are detected.
func (it *RunFromRootCommand) validateFlagCombinations(arguments []string) {
	hasParallelFlag := HasParallelFlag(arguments)
	hasNoParallelBypass := HasNoParallelBypassFlag(arguments)
	hasAutoAnswerFlag := it.hasAutoAnswerFlag(arguments)
	hasAllFlag := HasAllFlag(arguments)
	isStateCommand := IsStateManipulationCommand(arguments)

	// If --parallel is used without --no-parallel-bypass, terra handles parallel execution
	// In this case, flags intended for terragrunt should not be used
	if hasParallelFlag && !hasNoParallelBypass {
		// Error if --auto-answer is used (intended for terragrunt, not terra's parallel execution)
		if hasAutoAnswerFlag {
			logger.Fatalf("Error: --auto-answer flag is intended for terragrunt and should only be used with --no-parallel-bypass. " +
				"When using --parallel without --no-parallel-bypass, terra handles parallel execution and --auto-answer is not applicable.")
		}

		// Error if --all is used with --parallel (redundant, since --parallel already implies --all behavior)
		if hasAllFlag {
			logger.Fatalf("Error: --all flag is not needed when using --parallel flag. " +
				"The --parallel flag already executes across all modules. Remove --all or use --no-parallel-bypass to forward --parallel to terragrunt.")
		}
	}

	// If --no-parallel-bypass is used with state commands, error out (terragrunt doesn't handle state commands)
	if hasNoParallelBypass && isStateCommand {
		logger.Fatalf("Error: --no-parallel-bypass cannot be used with state commands. " +
			"Terragrunt doesn't support state operations, so state commands must be handled by terra. " +
			"Remove --no-parallel-bypass to let terra handle the parallel execution.")
	}

	// If --no-parallel-bypass is used, --all is required (for non-state commands)
	if hasNoParallelBypass && hasParallelFlag && !isStateCommand {
		if !hasAllFlag {
			logger.Fatalf("Error: --all flag is required when using --no-parallel-bypass with --parallel. " +
				"Terragrunt needs --all to understand that it should apply to all modules.")
		}
	}
}

// isParallelCommand checks if the command should be executed in parallel.
// Returns true if:
// 1. It's a state command with --all flag (backward compatibility)
// 2. It has --parallel=N flag (new functionality for any command), UNLESS --no-parallel-bypass is present
// If --no-parallel-bypass is present, --parallel flag will be forwarded to terragrunt instead
func (it *RunFromRootCommand) isParallelCommand(arguments []string) bool {
	// Check if --no-parallel-bypass is present
	hasNoParallelBypass := HasNoParallelBypassFlag(arguments)
	
	// New: support parallel=N for any command, but only if --no-parallel-bypass is NOT present
	if HasParallelFlag(arguments) && !hasNoParallelBypass {
		return true
	}
	// Backward compatibility: state commands with --all flag (always handled by terra, regardless of --no-parallel-bypass)
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
