package commands

import (
	"slices"
	"strconv"
	"strings"
)

// StateCommandConstants defines constants for state manipulation commands and flags.
const (
	// AllFlag represents the --all flag used with state commands.
	AllFlag = "--all"
	// ParallelFlagPrefix represents the prefix for the parallel flag.
	ParallelFlagPrefix = "--parallel="
	// NoParallelBypassFlag represents the --no-parallel-bypass flag.
	NoParallelBypassFlag = "--no-parallel-bypass"
	// IncludeFlagPrefix represents the prefix for the include flag.
	IncludeFlagPrefix = "--include="
	// ExcludeFlagPrefix represents the prefix for the exclude flag.
	ExcludeFlagPrefix = "--exclude="
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

// HasNoParallelBypassFlag checks if the --no-parallel-bypass flag is present in arguments.
func HasNoParallelBypassFlag(arguments []string) bool {
	return slices.Contains(arguments, NoParallelBypassFlag)
}

// RemoveNoParallelBypassFlag removes --no-parallel-bypass flag from arguments.
func RemoveNoParallelBypassFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != NoParallelBypassFlag {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// FilterValues represents separated inclusions and exclusions from the include/exclude flags.
type FilterValues struct {
	Inclusions []string
	Exclusions []string
}

// getCommaSeparatedFlagValues extracts comma-separated values from a flag with the given prefix.
func getCommaSeparatedFlagValues(arguments []string, prefix string) ([]string, bool) {
	for _, arg := range arguments {
		if after, ok := strings.CutPrefix(arg, prefix); ok {
			if after == "" {
				return nil, false
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

			return nil, false
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

// HasIncludeFlag checks if the --include= flag is present in arguments.
func HasIncludeFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, IncludeFlagPrefix)
}

// GetIncludeValues extracts comma-separated values from --include=value1,value2 flag.
func GetIncludeValues(arguments []string) ([]string, bool) {
	return getCommaSeparatedFlagValues(arguments, IncludeFlagPrefix)
}

// RemoveIncludeFlag removes --include= flag from arguments.
func RemoveIncludeFlag(arguments []string) []string {
	return removeFlagWithPrefix(arguments, IncludeFlagPrefix)
}

// HasExcludeFlag checks if the --exclude= flag is present in arguments.
func HasExcludeFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, ExcludeFlagPrefix)
}

// GetExcludeValues extracts comma-separated values from --exclude=value1,value2 flag.
func GetExcludeValues(arguments []string) ([]string, bool) {
	return getCommaSeparatedFlagValues(arguments, ExcludeFlagPrefix)
}

// RemoveExcludeFlag removes --exclude= flag from arguments.
func RemoveExcludeFlag(arguments []string) []string {
	return removeFlagWithPrefix(arguments, ExcludeFlagPrefix)
}

// GetFilterValues extracts inclusions and exclusions from --include= and --exclude= flags.
func GetFilterValues(arguments []string) FilterValues {
	var result FilterValues

	if values, found := GetIncludeValues(arguments); found {
		result.Inclusions = values
	}

	if values, found := GetExcludeValues(arguments); found {
		result.Exclusions = values
	}

	return result
}

// RemoveFilterFlags removes both --include= and --exclude= flags from arguments.
func RemoveFilterFlags(arguments []string) []string {
	filtered := RemoveIncludeFlag(arguments)
	return RemoveExcludeFlag(filtered)
}
