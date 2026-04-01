package commands

import (
	"slices"
	"strconv"
	"strings"
)

// StateCommandConstants defines constants for state manipulation commands and flags.
const (
	// AllFlag represents the --all flag (forwarded to terragrunt for non-state commands).
	AllFlag = "--all"
	// ParallelFlagPrefix represents the prefix for the parallel flag.
	ParallelFlagPrefix = "--parallel="
	// OnlyFlagPrefix represents the prefix for the --only flag (select specific modules).
	OnlyFlagPrefix = "--only="
	// SkipFlagPrefix represents the prefix for the --skip flag (exclude specific modules).
	SkipFlagPrefix = "--skip="

	// DeprecatedNoParallelBypassFlag is the removed --no-parallel-bypass flag.
	DeprecatedNoParallelBypassFlag = "--no-parallel-bypass"
	// DeprecatedIncludeFlagPrefix is the renamed --include flag (now --only).
	DeprecatedIncludeFlagPrefix = "--include="
	// DeprecatedExcludeFlagPrefix is the renamed --exclude flag (now --skip).
	DeprecatedExcludeFlagPrefix = "--exclude="
)

// IsStateManipulationCommand checks if the given arguments represent a state manipulation command.
func IsStateManipulationCommand(arguments []string) bool {
	if len(arguments) == 0 {
		return false
	}

	// State manipulation commands
	stateCommands := []string{
		"import", "state",
	}

	firstArg := arguments[0]
	if slices.Contains(stateCommands, firstArg) {
		return true
	}

	// Check for state subcommands (e.g., "state rm", "state mv").
	if len(arguments) >= 2 && firstArg == "state" {
		stateSubcommands := []string{
			"rm", "mv", "pull", "push", "show",
		}
		secondArg := arguments[1]
		if slices.Contains(stateSubcommands, secondArg) {
			return true
		}
	}

	return false
}

// HasAllFlag checks if the --all flag is present in arguments.
func HasAllFlag(arguments []string) bool {
	return slices.Contains(arguments, AllFlag)
}

// HasParallelFlag checks if the --parallel=N flag is present in arguments.
func HasParallelFlag(arguments []string) bool {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, ParallelFlagPrefix) {
			return true
		}
	}
	return false
}

// GetParallelValue extracts the number N from --parallel=N flag.
// Returns the value and true if found, or 0 and false if not found or invalid.
func GetParallelValue(arguments []string) (int, bool) {
	for _, arg := range arguments {
		if after, ok := strings.CutPrefix(arg, ParallelFlagPrefix); ok {
			valueStr := after
			value, err := strconv.Atoi(valueStr)
			if err != nil || value <= 0 {
				return 0, false
			}
			return value, true
		}
	}
	return 0, false
}

// RemoveParallelFlag removes --parallel=N flag from arguments.
func RemoveParallelFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, ParallelFlagPrefix) {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// HasDeprecatedNoParallelBypassFlag checks if the removed --no-parallel-bypass flag is present.
func HasDeprecatedNoParallelBypassFlag(arguments []string) bool {
	return slices.Contains(arguments, DeprecatedNoParallelBypassFlag)
}

// HasDeprecatedIncludeFlag checks if the renamed --include= flag is present.
func HasDeprecatedIncludeFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, DeprecatedIncludeFlagPrefix)
}

// HasDeprecatedExcludeFlag checks if the renamed --exclude= flag is present.
func HasDeprecatedExcludeFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, DeprecatedExcludeFlagPrefix)
}

// SelectionValues represents separated --only and --skip values for module selection.
type SelectionValues struct {
	Only []string
	Skip []string
}

// getCommaSeparatedFlagValues extracts comma-separated values from a flag with the given prefix.
// If a matching flag has an empty value, it is skipped so that a later valid occurrence can be used.
func getCommaSeparatedFlagValues(arguments []string, prefix string) ([]string, bool) {
	for _, arg := range arguments {
		if after, ok := strings.CutPrefix(arg, prefix); ok {
			if after == "" {
				continue
			}

			parts := strings.Split(after, ",")

			var result []string
			for _, part := range parts {
				trimmed := strings.TrimSpace(part)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}

			if len(result) > 0 {
				return result, true
			}
		}
	}

	return nil, false
}

// hasFlagWithPrefix checks if any argument starts with the given prefix.
func hasFlagWithPrefix(arguments []string, prefix string) bool {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}

	return false
}

// removeFlagWithPrefix removes arguments that start with the given prefix.
func removeFlagWithPrefix(arguments []string, prefix string) []string {
	var filtered []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, prefix) {
			filtered = append(filtered, arg)
		}
	}

	return filtered
}

// HasOnlyFlag checks if the --only= flag is present in arguments.
func HasOnlyFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, OnlyFlagPrefix)
}

// GetOnlyValues extracts comma-separated values from --only=value1,value2 flag.
func GetOnlyValues(arguments []string) ([]string, bool) {
	return getCommaSeparatedFlagValues(arguments, OnlyFlagPrefix)
}

// RemoveOnlyFlag removes --only= flag from arguments.
func RemoveOnlyFlag(arguments []string) []string {
	return removeFlagWithPrefix(arguments, OnlyFlagPrefix)
}

// HasSkipFlag checks if the --skip= flag is present in arguments.
func HasSkipFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, SkipFlagPrefix)
}

// GetSkipValues extracts comma-separated values from --skip=value1,value2 flag.
func GetSkipValues(arguments []string) ([]string, bool) {
	return getCommaSeparatedFlagValues(arguments, SkipFlagPrefix)
}

// RemoveSkipFlag removes --skip= flag from arguments.
func RemoveSkipFlag(arguments []string) []string {
	return removeFlagWithPrefix(arguments, SkipFlagPrefix)
}

// GetSelectionValues extracts --only and --skip values for module selection.
func GetSelectionValues(arguments []string) SelectionValues {
	var result SelectionValues

	if values, found := GetOnlyValues(arguments); found {
		result.Only = values
	}

	if values, found := GetSkipValues(arguments); found {
		result.Skip = values
	}

	return result
}

// RemoveSelectionFlags removes both --only= and --skip= flags from arguments.
func RemoveSelectionFlags(arguments []string) []string {
	filtered := RemoveOnlyFlag(arguments)
	return RemoveSkipFlag(filtered)
}
