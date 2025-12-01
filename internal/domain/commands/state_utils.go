package commands

import (
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
	// FilterFlagPrefix represents the prefix for the filter flag.
	FilterFlagPrefix = "--filter="
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
	for _, cmd := range stateCommands {
		if firstArg == cmd {
			return true
		}
	}

	// Check for state subcommands (e.g., "state rm", "state mv").
	if len(arguments) >= 2 && firstArg == "state" {
		stateSubcommands := []string{
			"rm", "mv", "pull", "push", "show",
		}
		secondArg := arguments[1]
		for _, subcmd := range stateSubcommands {
			if secondArg == subcmd {
				return true
			}
		}
	}

	return false
}

// HasAllFlag checks if the --all flag is present in arguments.
func HasAllFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == AllFlag {
			return true
		}
	}
	return false
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
		if strings.HasPrefix(arg, ParallelFlagPrefix) {
			valueStr := strings.TrimPrefix(arg, ParallelFlagPrefix)
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
	for _, arg := range arguments {
		if arg == NoParallelBypassFlag {
			return true
		}
	}
	return false
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

// HasFilterFlag checks if the --filter= flag is present in arguments.
func HasFilterFlag(arguments []string) bool {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, FilterFlagPrefix) {
			return true
		}
	}
	return false
}

// FilterValues represents separated inclusions and exclusions from the filter flag.
type FilterValues struct {
	Inclusions []string // Values without ! prefix
	Exclusions []string // Values with ! prefix (without the !)
}

// GetFilterValue extracts the comma-separated list from --filter=value1,value2,!value3 flag.
// Returns the separated inclusions and exclusions, and true if found, or nil and false if not found.
func GetFilterValue(arguments []string) ([]string, bool) {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, FilterFlagPrefix) {
			valueStr := strings.TrimPrefix(arg, FilterFlagPrefix)
			if valueStr == "" {
				return nil, false
			}
			// Split by comma and trim whitespace
			filters := strings.Split(valueStr, ",")
			var result []string
			for _, filter := range filters {
				trimmed := strings.TrimSpace(filter)
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

// ParseFilterValues separates filter values into inclusions and exclusions.
func ParseFilterValues(filterValues []string) FilterValues {
	var result FilterValues
	for _, value := range filterValues {
		if strings.HasPrefix(value, "!") {
			// Exclusion: remove the ! prefix
			exclusion := strings.TrimPrefix(value, "!")
			exclusion = strings.TrimSpace(exclusion)
			if exclusion != "" {
				result.Exclusions = append(result.Exclusions, exclusion)
			}
		} else {
			// Inclusion
			inclusion := strings.TrimSpace(value)
			if inclusion != "" {
				result.Inclusions = append(result.Inclusions, inclusion)
			}
		}
	}
	return result
}

// RemoveFilterFlag removes --filter= flag from arguments.
func RemoveFilterFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, FilterFlagPrefix) {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
