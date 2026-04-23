package commands

import (
	"fmt"
	"slices"
	"strings"
)

// sensitiveFlagPrefixes lists terraform/terragrunt flags whose values commonly
// carry secrets (variable values, backend credentials). Both single-dash and
// double-dash variants are handled. When these flags appear in argv, the error
// echo replaces their values with "<redacted>" so credentials cannot leak into
// error logs or CI output. Note: "-var-file" is NOT included because its value
// is a filename, not a secret.
func sensitiveFlagPrefixes() []string {
	return []string{
		"-var",
		"--var",
		"-backend-config",
		"--backend-config",
	}
}

// extractSubcommand returns the terragrunt command prefix from the argument list.
// For most commands this is the first non-flag argument (apply, plan, destroy,
// import, etc.). For multi-word state commands it preserves the second token as
// well (for example: "state rm", "state mv") so generated suggestions remain
// copy-pasteable. Returns an empty string when no non-flag argument is present.
func extractSubcommand(arguments []string) string {
	for index, arg := range arguments {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if arg != "state" {
			return arg
		}

		subcommandParts := []string{arg}
		for _, nextArg := range arguments[index+1:] {
			if strings.HasPrefix(nextArg, "-") {
				continue
			}

			subcommandParts = append(subcommandParts, nextArg)
			break
		}

		return strings.Join(subcommandParts, " ")
	}
	return ""
}

// sanitizeArgvForEcho returns a copy of arguments in which the values of known
// sensitive flags (see sensitiveFlagPrefixes) are replaced with "<redacted>".
// Both "-var=key=secret" and "-var key=secret" forms are handled. Filenames
// (e.g. "-var-file=prod.tfvars") are left untouched.
func sanitizeArgvForEcho(arguments []string) []string {
	sanitized := make([]string, 0, len(arguments))
	redactNext := false
	for _, arg := range arguments {
		if redactNext {
			sanitized = append(sanitized, "<redacted>")
			redactNext = false
			continue
		}
		if prefix, matched := matchSensitiveFlagWithValue(arg); matched {
			sanitized = append(sanitized, prefix+"=<redacted>")
			continue
		}
		if isSensitiveFlagWithoutValue(arg) {
			sanitized = append(sanitized, arg)
			redactNext = true
			continue
		}
		sanitized = append(sanitized, arg)
	}
	return sanitized
}

// matchSensitiveFlagWithValue returns the flag prefix (e.g. "-var") and true when
// the argument matches the "<flag>=<value>" form of a sensitive flag.
func matchSensitiveFlagWithValue(arg string) (string, bool) {
	for _, prefix := range sensitiveFlagPrefixes() {
		if strings.HasPrefix(arg, prefix+"=") {
			return prefix, true
		}
	}
	return "", false
}

// isSensitiveFlagWithoutValue returns true when the argument is exactly a sensitive
// flag name (e.g. "-var"), indicating the next argument carries its value.
func isSensitiveFlagWithoutValue(arg string) bool {
	return slices.Contains(sensitiveFlagPrefixes(), arg)
}

// buildEchoedCommand reconstructs the command the user typed so the error message
// can quote it back. Values of known sensitive flags (-var, -backend-config) are
// redacted so credentials cannot leak into error output. Format:
// "terra <sanitized arguments joined> <targetPath>".
func buildEchoedCommand(arguments []string, targetPath string) string {
	parts := []string{"terra"}
	parts = append(parts, sanitizeArgvForEcho(arguments)...)
	if targetPath != "" {
		parts = append(parts, targetPath)
	}
	return strings.Join(parts, " ")
}

// buildParallelSuggestion builds the terra-managed parallel form of the user's intent,
// suitable as a copy-pasteable example in a validation error. --yes is appended only
// when the command is interactive (apply/destroy), because terra rejects those without
// a confirmation flag.
func buildParallelSuggestion(
	subcommand string,
	onlyValues, skipValues []string,
	needsConfirmation bool,
	targetPath string,
) string {
	parts := []string{"terra"}
	if subcommand != "" {
		parts = append(parts, subcommand)
	}
	parts = append(parts, "--parallel=5")
	if len(onlyValues) > 0 {
		parts = append(parts, "--only="+strings.Join(onlyValues, ","))
	}
	if len(skipValues) > 0 {
		parts = append(parts, "--skip="+strings.Join(skipValues, ","))
	}
	if needsConfirmation {
		parts = append(parts, "--yes")
	}
	if targetPath != "" {
		parts = append(parts, targetPath)
	}
	return strings.Join(parts, " ")
}

// buildAllWithFilterSuggestion builds the terragrunt-managed run-all form of the user's
// intent, using terragrunt's --filter as the module selector. --only values map to
// positive filters, --skip values map to negated filters prefixed with '!'.
func buildAllWithFilterSuggestion(
	subcommand string,
	onlyValues, skipValues []string,
	targetPath string,
) string {
	parts := []string{"terra"}
	if subcommand != "" {
		parts = append(parts, subcommand)
	}
	parts = append(parts, "--all")
	for _, value := range onlyValues {
		parts = append(parts, fmt.Sprintf("--filter='%s'", value))
	}
	for _, value := range skipValues {
		parts = append(parts, fmt.Sprintf("--filter='!%s'", value))
	}
	if targetPath != "" {
		parts = append(parts, targetPath)
	}
	return strings.Join(parts, " ")
}

// BuildSelectionFlagsError produces an educational error message when --only or --skip
// is used without --parallel=N. It echoes the user's command verbatim and shows both
// valid alternative forms (terra-managed --parallel and terragrunt-managed --all --filter)
// with the user's own module values substituted.
func BuildSelectionFlagsError(arguments []string, targetPath string) string {
	subcommand := extractSubcommand(arguments)
	onlyValues, _ := GetOnlyValues(arguments)
	skipValues, _ := GetSkipValues(arguments)
	needsConfirmation := IsInteractiveCommand(arguments)

	echoed := buildEchoedCommand(arguments, targetPath)
	parallelSuggestion := buildParallelSuggestion(
		subcommand, onlyValues, skipValues, needsConfirmation, targetPath,
	)
	allSuggestion := buildAllWithFilterSuggestion(subcommand, onlyValues, skipValues, targetPath)

	return fmt.Sprintf(
		"Error: --only/--skip are terra-managed flags and require --parallel=N.\n"+
			"You used: %s\n"+
			"\n"+
			"Terra has two separate parallel-execution strategies. "+
			"--only/--skip only work with the first:\n"+
			"\n"+
			"  1) Terra-managed worker pool (simple basename matching across the tree):\n"+
			"       %s\n"+
			"\n"+
			"  2) Terragrunt-managed run-all "+
			"(DAG-aware; supports path globs, graph and git-diff expressions):\n"+
			"       %s\n"+
			"\n"+
			"Terragrunt's --filter is strictly more expressive than terra's --skip/--only.\n"+
			"See docs/parallel-execution.md for the full comparison.",
		echoed, parallelSuggestion, allSuggestion,
	)
}

// BuildParallelAllConflictError produces an educational error message when --parallel
// and --all are combined. It echoes the user's command and shows both valid forms so
// the user can pick one without re-reading the documentation.
func BuildParallelAllConflictError(arguments []string, targetPath string) string {
	subcommand := extractSubcommand(arguments)
	onlyValues, _ := GetOnlyValues(arguments)
	skipValues, _ := GetSkipValues(arguments)
	needsConfirmation := IsInteractiveCommand(arguments)

	echoed := buildEchoedCommand(arguments, targetPath)
	parallelSuggestion := buildParallelSuggestion(
		subcommand, onlyValues, skipValues, needsConfirmation, targetPath,
	)
	allSuggestion := buildAllWithFilterSuggestion(subcommand, onlyValues, skipValues, targetPath)

	return fmt.Sprintf(
		"Error: --parallel and --all cannot be used together "+
			"(competing execution strategies).\n"+
			"You used: %s\n"+
			"\n"+
			"Pick one strategy:\n"+
			"\n"+
			"  1) Terra-managed worker pool (drop --all):\n"+
			"       %s\n"+
			"\n"+
			"  2) Terragrunt-managed run-all (drop --parallel=N):\n"+
			"       %s\n"+
			"\n"+
			"See docs/parallel-execution.md for when to use each.",
		echoed, parallelSuggestion, allSuggestion,
	)
}
