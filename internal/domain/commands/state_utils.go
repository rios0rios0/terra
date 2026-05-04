package commands

import (
	"slices"
	"strconv"
	"strings"
)

// StateCommandConstants defines constants for state manipulation commands and flags.
const (
	// TerraCommandName is the name of the terra binary, used when echoing commands
	// back to the user inside educational error messages.
	TerraCommandName = "terra"
	// StateCommand is the terragrunt/terraform top-level "state" subcommand
	// (e.g. "state rm", "state mv"). Both IsStateManipulationCommand and the
	// error builders inspect arguments for this verb.
	StateCommand = "state"

	// AllFlag represents the --all flag (forwarded to terragrunt for non-state commands).
	AllFlag = "--all"
	// ParallelFlagPrefix represents the prefix for the parallel flag.
	ParallelFlagPrefix = "--parallel="
	// OnlyFlagPrefix represents the prefix for the --only flag (select specific modules).
	OnlyFlagPrefix = "--only="
	// SkipFlagPrefix represents the prefix for the --skip flag (exclude specific modules).
	SkipFlagPrefix = "--skip="

	// YesFlag represents the --yes flag (auto-approve, non-interactive).
	YesFlag = "--yes"
	// YesShortFlag represents the -y short flag for --yes.
	YesShortFlag = "-y"
	// NoFlag represents the --no flag (non-interactive, abort on confirm prompts).
	NoFlag = "--no"
	// NoShortFlag represents the -n short flag for --no.
	NoShortFlag = "-n"

	// DeprecatedNoParallelBypassFlag is the removed --no-parallel-bypass flag.
	DeprecatedNoParallelBypassFlag = "--no-parallel-bypass"
	// DeprecatedIncludeFlagPrefix is the renamed --include flag (now --only).
	DeprecatedIncludeFlagPrefix = "--include="
	// DeprecatedExcludeFlagPrefix is the renamed --exclude flag (now --skip).
	DeprecatedExcludeFlagPrefix = "--exclude="

	// FilterFlagPrefix represents the prefix for terragrunt's --filter flag.
	// Forwarded to terragrunt as-is; only meaningful with --all.
	FilterFlagPrefix = "--filter="
	// QueueExcludeDirFlag represents terragrunt's --queue-exclude-dir flag.
	// Forwarded to terragrunt as-is; only meaningful with --all.
	QueueExcludeDirFlag = "--queue-exclude-dir"
	// QueueIncludeDirFlag represents terragrunt's --queue-include-dir flag.
	// Forwarded to terragrunt as-is; only meaningful with --all.
	QueueIncludeDirFlag = "--queue-include-dir"
)

// IsStateManipulationCommand checks if the given arguments represent a state manipulation command.
func IsStateManipulationCommand(arguments []string) bool {
	if len(arguments) == 0 {
		return false
	}

	// State manipulation commands
	stateCommands := []string{
		"import", StateCommand,
	}

	firstArg := arguments[0]
	if slices.Contains(stateCommands, firstArg) {
		return true
	}

	// Check for state subcommands (e.g., "state rm", "state mv").
	if len(arguments) >= 2 && firstArg == StateCommand {
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

// HasAllFlag checks if the --all flag is present in arguments (including --all=true form).
func HasAllFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == AllFlag || strings.HasPrefix(arg, AllFlag+"=") {
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
	for _, arg := range arguments {
		if arg == DeprecatedNoParallelBypassFlag || strings.HasPrefix(arg, DeprecatedNoParallelBypassFlag+"=") {
			return true
		}
	}
	return false
}

// HasDeprecatedIncludeFlag checks if the renamed --include= flag is present.
func HasDeprecatedIncludeFlag(arguments []string) bool {
	return hasFlagWithPrefix(arguments, DeprecatedIncludeFlagPrefix)
}

// HasFilterFlag checks if terragrunt's --filter flag is present (either --filter=<query>
// or the space form --filter <query>). Used only by validation to warn when terragrunt
// queue flags are combined with terra's --parallel=N (where they have no effect).
func HasFilterFlag(arguments []string) bool {
	for i, arg := range arguments {
		if strings.HasPrefix(arg, FilterFlagPrefix) {
			return true
		}
		if arg == "--filter" && i+1 < len(arguments) {
			return true
		}
	}
	return false
}

// HasQueueExcludeDirFlag checks if terragrunt's --queue-exclude-dir flag is present
// (either --queue-exclude-dir=<dir> or the space form --queue-exclude-dir <dir>).
func HasQueueExcludeDirFlag(arguments []string) bool {
	for i, arg := range arguments {
		if strings.HasPrefix(arg, QueueExcludeDirFlag+"=") {
			return true
		}
		if arg == QueueExcludeDirFlag && i+1 < len(arguments) {
			return true
		}
	}
	return false
}

// HasQueueIncludeDirFlag checks if terragrunt's --queue-include-dir flag is present
// (either --queue-include-dir=<dir> or the space form --queue-include-dir <dir>).
func HasQueueIncludeDirFlag(arguments []string) bool {
	for i, arg := range arguments {
		if strings.HasPrefix(arg, QueueIncludeDirFlag+"=") {
			return true
		}
		if arg == QueueIncludeDirFlag && i+1 < len(arguments) {
			return true
		}
	}
	return false
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

// IsInteractiveCommand checks if the command triggers yes/no prompts in terragrunt.
// Skips leading flags (arguments starting with "-") to find the actual command.
func IsInteractiveCommand(arguments []string) bool {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return arg == "apply" || arg == "destroy"
	}
	return false
}

// HasReplyFlag checks if --reply, -r, --reply=<value>, or -r=<value> is present.
func HasReplyFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == ReplyFlag || arg == ReplyShortFlag ||
			strings.HasPrefix(arg, ReplyFlag+"=") ||
			strings.HasPrefix(arg, ReplyShortFlag+"=") {
			return true
		}
	}
	return false
}

// GetReplyValue extracts the value from --reply=<value> or -r=<value>.
// Returns "n" as default when the boolean form (--reply or -r) is used.
func GetReplyValue(arguments []string) string {
	for _, arg := range arguments {
		if arg == ReplyFlag || arg == ReplyShortFlag {
			return "n"
		}
		if strings.HasPrefix(arg, ReplyFlag+"=") {
			return arg[len(ReplyFlag+"="):]
		}
		if strings.HasPrefix(arg, ReplyShortFlag+"=") {
			return arg[len(ReplyShortFlag+"="):]
		}
	}
	return ""
}

// HasExplicitReplyValue returns true when --reply=<value> or -r=<value> form is used
// with a non-empty value suffix. Returns false for the boolean form (--reply / -r)
// and for empty-value forms like --reply= or -r=.
func HasExplicitReplyValue(arguments []string) bool {
	for _, arg := range arguments {
		if strings.HasPrefix(arg, ReplyFlag+"=") {
			value := arg[len(ReplyFlag+"="):]
			if value != "" {
				return true
			}
			continue
		}
		if strings.HasPrefix(arg, ReplyShortFlag+"=") {
			value := arg[len(ReplyShortFlag+"="):]
			if value != "" {
				return true
			}
		}
	}
	return false
}

// RemoveReplyFlag removes --reply, -r, --reply=<value>, and -r=<value> from arguments.
func RemoveReplyFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != ReplyFlag && arg != ReplyShortFlag &&
			!strings.HasPrefix(arg, ReplyFlag+"=") &&
			!strings.HasPrefix(arg, ReplyShortFlag+"=") {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// HasYesFlag checks if --yes or -y is present in arguments.
func HasYesFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == YesFlag || arg == YesShortFlag {
			return true
		}
	}
	return false
}

// HasNoFlag checks if --no or -n is present in arguments.
func HasNoFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == NoFlag || arg == NoShortFlag {
			return true
		}
	}
	return false
}

// HasConfirmationFlag returns true when any terra-level confirmation flag is present:
// --yes/-y, --no/-n, or the deprecated --reply/-r forms.
func HasConfirmationFlag(arguments []string) bool {
	return HasYesFlag(arguments) || HasNoFlag(arguments) || HasReplyFlag(arguments)
}

// RemoveConfirmationFlags removes every terra-level confirmation flag from arguments
// (--yes/-y, --no/-n, and all forms of the deprecated --reply/-r, including the
// space-separated form --reply y / -r y where the value is a separate argument).
//
// The space-form heuristic only consumes the next argument when it is a recognised
// reply value (y/yes/n/no/true/false/0/1). This avoids swallowing the subcommand
// in inputs like "--reply plan" (bare --reply followed by the Terragrunt command).
func RemoveConfirmationFlags(arguments []string) []string {
	filtered := make([]string, 0, len(arguments))
	skipNext := false
	for index, arg := range arguments {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == YesFlag || arg == YesShortFlag ||
			arg == NoFlag || arg == NoShortFlag {
			continue
		}
		if strings.HasPrefix(arg, ReplyFlag+"=") ||
			strings.HasPrefix(arg, ReplyShortFlag+"=") {
			continue
		}
		if arg == ReplyFlag || arg == ReplyShortFlag {
			if index+1 < len(arguments) && isRecognisedReplyValue(arguments[index+1]) {
				skipNext = true
			}
			continue
		}
		filtered = append(filtered, arg)
	}
	return filtered
}

// isRecognisedReplyValue returns true when the token is one of the literal reply
// values that could follow --reply / -r in the space-separated form. Anything
// else (subcommands, paths, targets) is not a reply value and must not be consumed.
func isRecognisedReplyValue(token string) bool {
	switch strings.ToLower(token) {
	case "y", "yes", "true", "1", "n", "no", "false", "0":
		return true
	default:
		return false
	}
}

// ResolveConfirmation returns which terra-level confirmation intent is present
// as a pair of booleans (isYes, isNo). Legacy --reply/-r forms are mapped:
// --reply=y and --reply y map to yes, --reply=n and --reply n map to no, and
// bare --reply (or -r) with no following value defaults to yes. The bare default
// matches the old --parallel path, which injected --non-interactive -auto-approve
// for bare --reply; the PTY path's "answer n" default was an accidental side
// effect that often blocked CI pipelines.
func ResolveConfirmation(arguments []string) (bool, bool) {
	if HasYesFlag(arguments) {
		return true, false
	}
	if HasNoFlag(arguments) {
		return false, true
	}

	if !HasReplyFlag(arguments) {
		return false, false
	}

	for index, arg := range arguments {
		if after, ok := strings.CutPrefix(arg, ReplyFlag+"="); ok {
			return mapReplyValueToConfirmation(after)
		}
		if after, ok := strings.CutPrefix(arg, ReplyShortFlag+"="); ok {
			return mapReplyValueToConfirmation(after)
		}
		if arg == ReplyFlag || arg == ReplyShortFlag {
			// Space-separated form: --reply y / -r n. Only consume the next
			// argument as the value when it is a recognised reply value, so
			// we don't misinterpret "--reply plan" (bare --reply + subcommand).
			if index+1 < len(arguments) && isRecognisedReplyValue(arguments[index+1]) {
				return mapReplyValueToConfirmation(arguments[index+1])
			}
			// Bare --reply / -r with no following value: default to yes.
			return true, false
		}
	}
	return true, false
}

// mapReplyValueToConfirmation translates a legacy --reply=<value> suffix into
// the equivalent (isYes, isNo) intent. Anything that looks like a negative
// answer maps to no; everything else (including unknown values) maps to yes.
func mapReplyValueToConfirmation(value string) (bool, bool) {
	switch strings.ToLower(value) {
	case "n", "no", "false", "0":
		return false, true
	default:
		return true, false
	}
}

// BuildConfirmationInjection returns the native Terraform/Terragrunt flags that
// should be appended to the forwarded argument list based on the presence of
// --yes/--no (or the deprecated --reply forms).
//
//   - --yes injects --non-interactive (so Terragrunt assumes "yes" for its own
//     prompts) and, for apply/destroy, -auto-approve (so Terraform skips its own
//     confirmation, which --non-interactive alone does not suppress).
//   - --no injects only --non-interactive. Without -auto-approve, Terraform's
//     apply/destroy prompt aborts instead of proceeding, which matches "no".
//
// Returns nil when no confirmation flag is present.
func BuildConfirmationInjection(arguments []string) []string {
	yes, no := ResolveConfirmation(arguments)
	if !yes && !no {
		return nil
	}

	injected := []string{"--non-interactive"}
	if yes && IsInteractiveCommand(arguments) {
		injected = append(injected, "-auto-approve")
	}
	return injected
}
