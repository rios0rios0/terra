package repositories

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// getUpgradePatterns returns error output patterns that indicate terraform/terragrunt
// needs initialization with --upgrade before the command can succeed.
func getUpgradePatterns() []string {
	return []string{
		// Module not initialized.
		"terraform init has not been run",
		"Working directory is not initialized",
		"run \"terraform init\"",
		"run \"terragrunt init\"",
		"Module not installed",
		"Required plugins are not installed",

		// Backend configuration changes.
		"Backend configuration changed",
		"backend type changed",
		"backend configuration has changed",
		"Error loading state",

		// Provider version conflicts.
		"provider version constraint",
		"Provider doesn't satisfy version constraints",
		"Inconsistent dependency lock file",
		"provider registry.terraform.io",
		"Failed to query available provider packages",

		// Terragrunt-specific patterns.
		"You must run 'terragrunt init --upgrade'",
		"terraform init -upgrade",
		"rerun this command to reinitialize",
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

	if !needsUpgrade(output) {
		return err
	}

	logger.Infof(
		"Detected that %s needs initialization with upgrade, running 'init --upgrade'...",
		command,
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

	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdin = os.Stdin

	var outputBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &outputBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &outputBuf)

	err := cmd.Run()
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

	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
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

	cmd := exec.CommandContext(context.Background(), command, initArgs...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to perform init --upgrade: %w", err)
	}

	return nil
}

// needsUpgrade checks if the command output contains patterns indicating
// that terraform/terragrunt needs initialization with --upgrade.
func needsUpgrade(output string) bool {
	lowerOutput := strings.ToLower(output)

	for _, pattern := range getUpgradePatterns() {
		if strings.Contains(lowerOutput, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// NeedsUpgradePublic is a public wrapper for testing the private needsUpgrade function.
func NeedsUpgradePublic(output string) bool {
	return needsUpgrade(output)
}
