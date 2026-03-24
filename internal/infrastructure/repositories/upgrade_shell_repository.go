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
		"rerun this command to reinitialize",
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

	initErr := it.runInitUpgrade(command, directory)
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

// runInitUpgrade runs "command init --upgrade" in the given directory.
func (it *UpgradeAwareShellRepository) runInitUpgrade(
	command string,
	directory string,
) error {
	initArgs := []string{"init", "--upgrade"}
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

// needsUpgrade checks if the command output contains patterns indicating
// that terraform/terragrunt needs initialization with --upgrade.
// Returns the matched pattern (empty string if no match).
func needsUpgrade(output string) string {
	lowerOutput := strings.ToLower(output)

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
