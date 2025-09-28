package repositories

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// UpgradeAwareShellRepository extends StdShellRepository with auto-upgrade detection capabilities.
type UpgradeAwareShellRepository struct {
	*StdShellRepository
}

func NewUpgradeAwareShellRepository() *UpgradeAwareShellRepository {
	return &UpgradeAwareShellRepository{
		StdShellRepository: NewStdShellRepository(),
	}
}

// ExecuteCommandWithUpgradeDetection executes a command and automatically runs init --upgrade if needed.
func (it *UpgradeAwareShellRepository) ExecuteCommandWithUpgradeDetection(
	command string,
	arguments []string,
	directory string,
) error {
	logger.Infof("Running [%s %s] in %s with upgrade detection", command, strings.Join(arguments, " "), directory)

	// Capture command output for analysis
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	// Forward output to console
	if stdout.Len() > 0 {
		_, _ = fmt.Fprint(os.Stdout, stdout.String())
	}
	if stderr.Len() > 0 {
		_, _ = fmt.Fprint(os.Stderr, stderr.String())
	}

	if err != nil {
		// Check if the error is due to initialization or upgrade needs
		combinedOutput := stdout.String() + stderr.String()
		if it.needsUpgrade(combinedOutput) {
			logger.Info("Detected that terraform/terragrunt needs initialization with upgrade, " +
				"running 'init --upgrade'...")

			// Run init --upgrade
			initErr := it.ExecuteCommand(command, []string{"init", "--upgrade"}, directory)
			if initErr != nil {
				return fmt.Errorf("failed to run init --upgrade: %w", initErr)
			}

			logger.Info("Init --upgrade completed successfully, retrying original command...")

			// Retry the original command
			return it.ExecuteCommand(command, arguments, directory)
		}

		return fmt.Errorf("failed to perform command execution: %w", err)
	}

	return nil
}

// needsUpgrade analyzes command output to determine if init --upgrade is needed.
func (it *UpgradeAwareShellRepository) needsUpgrade(output string) bool {
	// Common patterns that indicate initialization or upgrade is needed
	upgradePatterns := []string{
		// Terraform/Terragrunt not initialized
		`terraform init.*has not been run`,
		`Working directory is not initialized`,
		`This working directory has been initialized`,
		`Backend initialization required`,

		// Backend configuration changes
		`Backend configuration changed`,
		`backend type changed`,
		`backend configuration has changed`,
		`Terraform detected that the backend type changed`,

		// Provider version conflicts
		`provider version constraint`,
		`provider.*version constraint`,
		`required provider.*not installed`,
		`provider configuration has changed`,
		`Provider .* doesn't satisfy version constraints`,
		`incompatible provider version`,

		// Module upgrade needs
		`terraform init -upgrade`,
		`run.*terraform init.*upgrade`,
		`You must run.*init.*upgrade`,
		`terraform init --upgrade`,
		`terragrunt init --upgrade`,

		// Lock file issues
		`dependency lock file`,
		`lock file is read-only`,
		`provider registry.*lock file`,

		// General initialization errors
		`run "terraform init"`,
		`run "terragrunt init"`,
		`Please run.*init`,
		`initialization required`,
	}

	for _, pattern := range upgradePatterns {
		matched, err := regexp.MatchString(`(?i)`+pattern, output)
		if err != nil {
			logger.Debugf("Error matching pattern '%s': %v", pattern, err)
			continue
		}
		if matched {
			logger.Debugf("Detected upgrade need with pattern: %s", pattern)
			return true
		}
	}

	return false
}
