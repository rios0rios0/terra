package commands

import (
	"fmt"
	"strings"
)

// extractSubcommand returns the first non-flag argument from the argument list,
// which is the terragrunt subcommand (apply, plan, destroy, import, etc.). Returns
// an empty string when no non-flag argument is present.
func extractSubcommand(arguments []string) string {
	for _, arg := range arguments {
		if !strings.HasPrefix(arg, "-") {
			return arg
		}
	}
	return ""
}

// buildEchoedCommand reconstructs the command the user typed so the error message
// can quote it back verbatim. Format: "terra <arguments joined> <targetPath>".
func buildEchoedCommand(arguments []string, targetPath string) string {
	parts := []string{"terra"}
	parts = append(parts, arguments...)
	if targetPath != "" {
		parts = append(parts, targetPath)
	}
	return strings.Join(parts, " ")
}

// buildParallelSuggestion builds the terra-managed parallel form of the user's intent,
// suitable as a copy-pasteable example in a validation error. --reply is appended only
// when the command is interactive (apply/destroy), because terra rejects those without it.
func buildParallelSuggestion(
	subcommand string,
	onlyValues, skipValues []string,
	needsReply bool,
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
	if needsReply {
		parts = append(parts, "--reply")
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
	needsReply := IsInteractiveCommand(arguments)

	echoed := buildEchoedCommand(arguments, targetPath)
	parallelSuggestion := buildParallelSuggestion(
		subcommand, onlyValues, skipValues, needsReply, targetPath,
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
	needsReply := IsInteractiveCommand(arguments)

	echoed := buildEchoedCommand(arguments, targetPath)
	parallelSuggestion := buildParallelSuggestion(
		subcommand, onlyValues, skipValues, needsReply, targetPath,
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
