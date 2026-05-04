package repositories

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

// getCancellationPatterns returns output patterns that indicate the user cancelled the operation.
// When a cancellation is detected, the upgrade retry must be skipped even if upgrade patterns
// also appear in the output, because the failure was intentional.
func getCancellationPatterns() []string {
	return []string{
		"Apply cancelled",
		"Plan cancelled",
		"Destroy cancelled",
	}
}

// getUpgradePatterns returns error output patterns that indicate terraform/terragrunt
// needs initialization with --upgrade before the command can succeed.
// Only patterns where Terraform/Terragrunt explicitly suggests running init are included;
// broad patterns that match runtime provider errors (e.g., TLS failures) are excluded.
func getUpgradePatterns() []string {
	return []string{
		// Explicit "run init" suggestions from Terraform/Terragrunt output.
		"please run \"terraform init\"",
		"run \"terraform init\"",
		"run \"terragrunt init\"",
		"please run 'terraform init'",
		"terraform init -upgrade",
		"You must run 'terragrunt init --upgrade'",
		"install it automatically by running",

		// Initialization-required diagnostics (exact Terraform source strings).
		"Backend initialization required",
		"State store initialization required",
		"initialization required: please run",

		// Uninitialized working directory.
		"terraform init has not been run",
		"Working directory is not initialized",

		// Backend configuration change detection.
		"Backend configuration changed",
		"backend type changed",
		"backend configuration has changed",

		// Module and plugin initialization.
		"Module not installed",
		"Required plugins are not installed",
		// Module source change detection. The single-line diagnostic title
		// "Module source has changed" (and the surrounding sentence) does not
		// rely on the multi-line "Run \"terraform init\"" hint, which Terragrunt
		// splits across stderr lines with its per-line "<ts> STDERR <cmd>: │ "
		// prefix and therefore breaks the substring match used by needsUpgrade.
		"Module source has changed",
		"source address was changed since this module was installed",

		// Provider version/lock file conflicts requiring init --upgrade.
		"Inconsistent dependency lock file",
		"does not match configured version constraint",
		"Provider doesn't satisfy version constraints",
	}
}

// UpgradeAwareShellRepository wraps command execution with automatic upgrade detection.
// When a command fails and the output matches known upgrade-needed patterns, it
// automatically runs "init --upgrade" and retries the original command.
type UpgradeAwareShellRepository struct{}

// NewUpgradeAwareShellRepository creates a new UpgradeAwareShellRepository.
func NewUpgradeAwareShellRepository() *UpgradeAwareShellRepository {
	return &UpgradeAwareShellRepository{}
}

// ExecuteCommandWithUpgrade runs the command, captures output, and if the command fails
// with patterns indicating an upgrade is needed, runs "init --upgrade" and retries.
func (it *UpgradeAwareShellRepository) ExecuteCommandWithUpgrade(
	command string,
	arguments []string,
	directory string,
) error {
	output, err := it.executeAndCapture(command, arguments, directory)
	if err == nil {
		return nil
	}

	matchedPattern := needsUpgrade(output)
	if matchedPattern == "" {
		return err
	}

	logger.Infof(
		"Detected that %s needs initialization with upgrade (matched pattern: %q), running 'init --upgrade'...",
		command, matchedPattern,
	)

	initErr := it.runInitUpgrade(command, arguments, directory)
	if initErr != nil {
		logger.Errorf("Init --upgrade failed: %s", initErr)
		return fmt.Errorf("auto init --upgrade failed: %w (original error: %w)", initErr, err)
	}

	logger.Infof("Init --upgrade completed successfully, retrying original command...")

	return it.executePassthrough(command, arguments, directory)
}

// executeAndCapture runs a command while streaming output to stdout/stderr AND capturing
// the combined output into a buffer for pattern analysis.
func (it *UpgradeAwareShellRepository) executeAndCapture(
	command string,
	arguments []string,
	directory string,
) (string, error) {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(arguments, " "), directory)

	start := time.Now()
	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdin = os.Stdin

	var outputBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &outputBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &outputBuf)

	err := cmd.Run()
	logCommandDuration(command, arguments, directory, time.Since(start), err)
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}

	return outputBuf.String(), err
}

// executePassthrough runs a command with direct stdout/stderr (no capture), used for the retry.
func (it *UpgradeAwareShellRepository) executePassthrough(
	command string,
	arguments []string,
	directory string,
) error {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(arguments, " "), directory)

	start := time.Now()
	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	logCommandDuration(command, arguments, directory, time.Since(start), err)
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}

	return err
}

// runInitUpgrade runs "command init --upgrade" in the given directory, propagating
// any queue-scoping flags (--all, --filter, --queue-*) from the original command so
// the init walks the same set of units. Without this, an "--all" command that fails
// in a parent directory would retry init in that parent directory — which has no
// terragrunt.hcl — and fail with "You attempted to run terragrunt in a folder that
// does not contain a terragrunt.hcl file".
func (it *UpgradeAwareShellRepository) runInitUpgrade(
	command string,
	originalArguments []string,
	directory string,
) error {
	initArgs := append([]string{"init", "--upgrade"}, extractQueueScopingFlags(originalArguments)...)
	logger.Infof("Running [%s %s] in %s", command, strings.Join(initArgs, " "), directory)

	start := time.Now()
	cmd := exec.CommandContext(context.Background(), command, initArgs...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	logCommandDuration(command, initArgs, directory, time.Since(start), err)
	if err != nil {
		return fmt.Errorf("failed to perform init --upgrade: %w", err)
	}

	return nil
}

// queueScopingFlagSpec describes a flag that scopes a terragrunt run to a queue
// of units. When the original command included one of these, the auto-init
// --upgrade retry must include it too.
type queueScopingFlagSpec struct {
	name       string
	takesValue bool
}

// queueScopingFlags lists every terragrunt flag that controls which units a
// queued run touches. Boolean flags are propagated as-is; valued flags are
// propagated together with the next argument when the space form is used.
func queueScopingFlags() []queueScopingFlagSpec {
	return []queueScopingFlagSpec{
		{"--all", false},
		{"--queue-strict-include", false},
		{"--queue-include-external", false},
		{"--queue-exclude-external", false},
		{"--filter", true},
		{"--queue-include-dir", true},
		{"--queue-exclude-dir", true},
		{"--queue-include-units-reading", true},
	}
}

// extractQueueScopingFlags returns the subset of arguments that scope a terragrunt
// run to a queue of units. Both "--flag value" and "--flag=value" forms are
// recognised. The returned slice preserves the original relative order.
func extractQueueScopingFlags(arguments []string) []string {
	specs := queueScopingFlags()
	bySpec := make(map[string]queueScopingFlagSpec, len(specs))
	for _, spec := range specs {
		bySpec[spec.name] = spec
	}

	var forwarded []string
	for index := 0; index < len(arguments); index++ {
		arg := arguments[index]

		if eq := strings.IndexByte(arg, '='); eq > 0 {
			if _, ok := bySpec[arg[:eq]]; ok {
				forwarded = append(forwarded, arg)
			}
			continue
		}

		spec, ok := bySpec[arg]
		if !ok {
			continue
		}
		forwarded = append(forwarded, arg)
		if spec.takesValue && index+1 < len(arguments) {
			forwarded = append(forwarded, arguments[index+1])
			index++
		}
	}
	return forwarded
}

// needsUpgrade checks if the command output contains patterns indicating
// that terraform/terragrunt needs initialization with --upgrade.
// Returns the matched pattern (empty string if no match).
// If the output indicates the user cancelled the operation, upgrade is skipped.
func needsUpgrade(output string) string {
	lowerOutput := strings.ToLower(output)

	for _, pattern := range getCancellationPatterns() {
		if strings.Contains(lowerOutput, strings.ToLower(pattern)) {
			return ""
		}
	}

	for _, pattern := range getUpgradePatterns() {
		if strings.Contains(lowerOutput, strings.ToLower(pattern)) {
			return pattern
		}
	}

	return ""
}

// NeedsUpgradePublic is a public wrapper for testing the private needsUpgrade function.
func NeedsUpgradePublic(output string) string {
	return needsUpgrade(output)
}
