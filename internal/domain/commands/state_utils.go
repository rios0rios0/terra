package commands

// StateCommandConstants defines constants for state manipulation commands and flags.
const (
	// AllFlag represents the --all flag used with state commands.
	AllFlag = "--all"
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
