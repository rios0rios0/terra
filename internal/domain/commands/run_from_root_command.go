package commands

import (
	"os"
	"strings"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

const (
	// ReplyFlag represents the --reply flag for auto-answering terragrunt prompts.
	ReplyFlag = "--reply"
	// ReplyShortFlag represents the -r short flag for --reply.
	ReplyShortFlag = "-r"
	// DeprecatedAutoAnswerFlag is the renamed --auto-answer flag (now --reply).
	DeprecatedAutoAnswerFlag = "--auto-answer"
	// DeprecatedAutoAnswerShortFlag is the removed -a short flag (collides with terragrunt's -a for --all).
	DeprecatedAutoAnswerShortFlag = "-a"
)

type RunFromRootCommand struct {
	settings              *entities.Settings
	installCommand        InstallDependencies
	formatCommand         FormatFiles
	additionalBefore      RunAdditionalBefore
	parallelState         ParallelState
	repository            repositories.ShellRepository
	upgradeRepository     repositories.UpgradeShellRepository
	interactiveRepository repositories.InteractiveShellRepository
}

func NewRunFromRootCommand(
	settings *entities.Settings,
	installCommand InstallDependencies,
	formatCommand FormatFiles,
	additionalBefore RunAdditionalBefore,
	parallelState ParallelState,
	repository repositories.ShellRepository,
	upgradeRepository repositories.UpgradeShellRepository,
	interactiveRepository repositories.InteractiveShellRepository,
) *RunFromRootCommand {
	return &RunFromRootCommand{
		settings:              settings,
		installCommand:        installCommand,
		formatCommand:         formatCommand,
		additionalBefore:      additionalBefore,
		parallelState:         parallelState,
		repository:            repository,
		upgradeRepository:     upgradeRepository,
		interactiveRepository: interactiveRepository,
	}
}

func (it *RunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	// Configure centralized cache directories before any Terragrunt invocation
	it.configureCacheEnvironment()

	// Skip formatting for state commands: state operations (mv, rm, etc.) don't modify
	// source code, so formatting is unnecessary. Skipping it also avoids file contention
	// when multiple terra processes run concurrently from the same repository.
	if !IsStateManipulationCommand(arguments) {
		it.formatCommand.Execute(dependencies)
	}

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

	// Check if reply flag is present and filter it out
	useInteractive := it.hasReplyFlag(arguments)
	replyValue := it.getReplyValue(arguments)
	filteredArguments := it.removeReplyFlag(arguments)

	var err error
	if useInteractive {
		logger.Infof("Using interactive mode with auto-replying (%s)", replyValue)
		err = it.interactiveRepository.ExecuteCommandWithAnswer(
			"terragrunt", filteredArguments, targetPath, replyValue)
	} else {
		// Use upgrade-aware repository: automatically detects when init --upgrade
		// is needed, runs it, and retries the original command.
		err = it.upgradeRepository.ExecuteCommandWithUpgrade(
			"terragrunt", filteredArguments, targetPath)
	}

	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

func (it *RunFromRootCommand) hasReplyFlag(arguments []string) bool {
	return HasReplyFlag(arguments)
}

func (it *RunFromRootCommand) getReplyValue(arguments []string) string {
	return GetReplyValue(arguments)
}

func (it *RunFromRootCommand) removeReplyFlag(arguments []string) []string {
	return RemoveReplyFlag(arguments)
}

// validateFlagCombinations validates that flag combinations are correct.
// Errors and exits if invalid combinations are detected.
func (it *RunFromRootCommand) validateFlagCombinations(arguments []string) {
	it.validateDeprecatedFlags(arguments)

	hasParallelFlag := HasParallelFlag(arguments)
	hasAllFlag := HasAllFlag(arguments)

	// --parallel and --all cannot be used together (competing execution strategies)
	if hasParallelFlag && hasAllFlag {
		logger.Fatalf(
			"Error: --parallel and --all cannot be used together. " +
				"Use --parallel=N for terra's parallel execution, or --all for terragrunt's run-all.",
		)
	}

	// --reply is required when --parallel is used with interactive commands (apply, destroy)
	// because parallel workers cannot handle stdin prompts
	if hasParallelFlag && IsInteractiveCommand(arguments) && !HasReplyFlag(arguments) {
		logger.Fatalf(
			"Error: --reply is required when using --parallel with apply or destroy. " +
				"Parallel workers cannot handle interactive prompts. Example: --reply=y",
		)
	}

	it.validateSelectionFlags(arguments, hasParallelFlag)
}

// validateDeprecatedFlags detects removed/renamed flags and exits with migration guidance.
func (it *RunFromRootCommand) validateDeprecatedFlags(arguments []string) {
	// Detect -a short flag (removed: collides with terragrunt's -a for --all)
	for _, arg := range arguments {
		if arg == DeprecatedAutoAnswerShortFlag || strings.HasPrefix(arg, DeprecatedAutoAnswerShortFlag+"=") {
			logger.Fatalf(
				"Error: the -a short flag has been removed (conflicts with terragrunt's -a for --all). " +
					"Use --reply or -r instead. Example: --reply=y",
			)
		}
	}

	// Detect --auto-answer (renamed to --reply)
	for _, arg := range arguments {
		if arg == DeprecatedAutoAnswerFlag || strings.HasPrefix(arg, DeprecatedAutoAnswerFlag+"=") {
			logger.Fatalf(
				"Error: --auto-answer has been renamed to --reply. " +
					"Use --reply or -r instead. Example: --reply=y",
			)
		}
	}

	// Detect --all with state commands (no longer intercepted by terra)
	if HasAllFlag(arguments) && IsStateManipulationCommand(arguments) {
		logger.Fatalf(
			"Error: --all cannot be used with state commands (terragrunt does not support this). " +
				"Use --parallel=5 instead. Example: terra import --parallel=5 <address> <id> <directory>",
		)
	}

	// Detect --no-parallel-bypass (removed entirely)
	if HasDeprecatedNoParallelBypassFlag(arguments) {
		logger.Fatalf(
			"Error: --no-parallel-bypass has been removed. " +
				"Use terragrunt's --parallelism=N directly for terragrunt-managed parallelism.",
		)
	}

	// Detect --include= (renamed to --only=)
	if HasDeprecatedIncludeFlag(arguments) {
		logger.Fatalf(
			"Error: --include has been renamed to --only. " +
				"Use --only=mod1,mod2 instead.",
		)
	}

	// Detect --exclude= (renamed to --skip=)
	if HasDeprecatedExcludeFlag(arguments) {
		logger.Fatalf(
			"Error: --exclude has been renamed to --skip. " +
				"Use --skip=mod1,mod2 instead.",
		)
	}
}

// validateSelectionFlags validates --only/--skip flag usage.
func (it *RunFromRootCommand) validateSelectionFlags(arguments []string, hasParallelFlag bool) {
	hasOnlyFlag := HasOnlyFlag(arguments)
	hasSkipFlag := HasSkipFlag(arguments)

	if !hasOnlyFlag && !hasSkipFlag {
		return
	}

	it.validateSelectionFlagValues(arguments, hasOnlyFlag, hasSkipFlag)

	// --only/--skip require --parallel=N
	if !hasParallelFlag {
		logger.Fatalf("Error: --only/--skip flags require --parallel=N.")
	}

	it.validateSelectionFlagConflicts(arguments, hasOnlyFlag, hasSkipFlag)
}

// validateSelectionFlagValues ensures present --only/--skip flags have non-empty values.
func (it *RunFromRootCommand) validateSelectionFlagValues(
	arguments []string,
	hasOnlyFlag, hasSkipFlag bool,
) {
	if hasOnlyFlag {
		if values, found := GetOnlyValues(arguments); !found || len(values) == 0 {
			logger.Fatalf("Error: --only flag is present but has no values. " +
				"Provide comma-separated module names, e.g. --only=mod1,mod2.")
		}
	}

	if hasSkipFlag {
		if values, found := GetSkipValues(arguments); !found || len(values) == 0 {
			logger.Fatalf("Error: --skip flag is present but has no values. " +
				"Provide comma-separated module names, e.g. --skip=mod1,mod2.")
		}
	}
}

// validateSelectionFlagConflicts detects modules appearing in both --only and --skip.
func (it *RunFromRootCommand) validateSelectionFlagConflicts(
	arguments []string,
	hasOnlyFlag, hasSkipFlag bool,
) {
	if !hasOnlyFlag || !hasSkipFlag {
		return
	}

	onlyValues, _ := GetOnlyValues(arguments)
	skipValues, _ := GetSkipValues(arguments)

	for _, only := range onlyValues {
		for _, skip := range skipValues {
			if only == skip {
				logger.Fatalf(
					"Error: module %q appears in both --only and --skip. Remove it from one flag.", only,
				)
			}
		}
	}
}

// isParallelCommand checks if the command should be executed in parallel by terra.
// Returns true if --parallel=N flag is present.
func (it *RunFromRootCommand) isParallelCommand(arguments []string) bool {
	return HasParallelFlag(arguments)
}

// configureCacheEnvironment sets environment variables for centralized Terragrunt module
// and provider caching. This ensures all stacks share a single download directory and
// provider cache, avoiding redundant downloads. It enables the Terragrunt Provider Cache
// Server (TG_PROVIDER_CACHE) by default for concurrent-safe provider deduplication, and
// the CAS (Content Addressable Store) experiment for Git clone deduplication.
func (it *RunFromRootCommand) configureCacheEnvironment() {
	const dirPermissions = 0o755

	moduleDir, moduleDirErr := it.settings.GetModuleCacheDir()
	if moduleDirErr != nil {
		logger.Warnf("Could not determine module cache directory: %s", moduleDirErr)
	} else if mkdirErr := os.MkdirAll(moduleDir, dirPermissions); mkdirErr != nil { // nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		logger.Warnf("Could not create module cache directory %s: %s", moduleDir, mkdirErr)
	} else if setenvErr := os.Setenv("TG_DOWNLOAD_DIR", moduleDir); setenvErr != nil {
		logger.Warnf("Could not set TG_DOWNLOAD_DIR: %s", setenvErr)
	} else {
		logger.Debugf("Module cache directory set to %s", moduleDir)
	}

	providerDir, providerDirErr := it.settings.GetProviderCacheDir()
	if providerDirErr != nil {
		logger.Warnf("Could not determine provider cache directory: %s", providerDirErr)
	} else if mkdirErr := os.MkdirAll(providerDir, dirPermissions); mkdirErr != nil { // nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		logger.Warnf("Could not create provider cache directory %s: %s", providerDir, mkdirErr)
	} else if setenvErr := os.Setenv("TG_PROVIDER_CACHE_DIR", providerDir); setenvErr != nil {
		logger.Warnf("Could not set TG_PROVIDER_CACHE_DIR: %s", setenvErr)
	} else {
		logger.Debugf("Provider cache directory set to %s", providerDir)
	}

	// Explicitly unset TF_PLUGIN_CACHE_DIR to prevent conflicts from inherited environment.
	// TF_PLUGIN_CACHE_DIR causes "text file busy" errors during parallel execution because
	// Terraform directly creates symlinks without file locking.
	if err := os.Unsetenv("TF_PLUGIN_CACHE_DIR"); err != nil {
		logger.Warnf("Could not unset TF_PLUGIN_CACHE_DIR: %s", err)
	}

	// Enable Terragrunt Provider Cache Server by default.
	// The server starts a localhost proxy that downloads each provider once with file
	// locking and creates symlinks for subsequent modules. This is safe for concurrent
	// access from parallel goroutines, unlike TF_PLUGIN_CACHE_DIR which causes "text
	// file busy" errors during parallel execution.
	setOrUnsetEnv("TG_PROVIDER_CACHE", "1", it.settings.TerraNoProviderCache)

	// Enable Terragrunt CAS (Content Addressable Store) experiment by default.
	// CAS deduplicates Git clones via hard links, reducing disk usage and clone times.
	setOrUnsetEnv("TG_EXPERIMENT", "cas", it.settings.TerraNoCAS)

	// Enable Terragrunt Partial Parse Config Cache by default.
	// Caches parsed HCL configs across modules sharing the same root include,
	// speeding up config parsing in large codebases.
	setOrUnsetEnv(
		"TG_USE_PARTIAL_PARSE_CONFIG_CACHE",
		"true",
		it.settings.TerraNoPartialParseCache,
	)
}

// setOrUnsetEnv sets the environment variable to the given value when disabled is false,
// or unsets it when disabled is true to ensure the opt-out is deterministic.
func setOrUnsetEnv(key, value string, disabled bool) {
	if disabled {
		if err := os.Unsetenv(key); err != nil {
			logger.Warnf("Could not unset %s: %s", key, err)
		}
		return
	}

	if err := os.Setenv(key, value); err != nil {
		logger.Warnf("Could not set %s: %s", key, err)
	} else {
		logger.Debugf("%s set to %s", key, value)
	}
}

// ConfigureCacheEnvironmentPublic is a public wrapper for testing the private configureCacheEnvironment method.
func (it *RunFromRootCommand) ConfigureCacheEnvironmentPublic() {
	it.configureCacheEnvironment()
}

// HasReplyFlagPublic is a public wrapper for testing the private hasReplyFlag method.
func (it *RunFromRootCommand) HasReplyFlagPublic(arguments []string) bool {
	return it.hasReplyFlag(arguments)
}

// GetReplyValuePublic is a public wrapper for testing the private getReplyValue method.
func (it *RunFromRootCommand) GetReplyValuePublic(arguments []string) string {
	return it.getReplyValue(arguments)
}

// RemoveReplyFlagPublic is a public wrapper for testing the private removeReplyFlag method.
func (it *RunFromRootCommand) RemoveReplyFlagPublic(arguments []string) []string {
	return it.removeReplyFlag(arguments)
}
