package commands

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
)

const (
	notInstalledVersion    = "not installed"
	latestAvailableVersion = "latest available"
	commandTimeout         = 10 * time.Second
)

// TerraVersion will be set at build time via ldflags
//
//nolint:gochecknoglobals // Version set at build time via ldflags
var TerraVersion = "1.5.0"

type VersionCommand struct {
	dependencies []entities.Dependency
}

func NewVersionCommand(dependencies []entities.Dependency) *VersionCommand {
	return &VersionCommand{
		dependencies: dependencies,
	}
}

func (it *VersionCommand) Execute() {
	logger.Infof("Terra version: %s", TerraVersion)

	// Get Terraform version
	terraformVersion := it.getTerraformVersion()
	logger.Infof("Terraform version: %s", terraformVersion)

	// Get Terragrunt version
	terragruntVersion := it.getTerragruntVersion()
	logger.Infof("Terragrunt version: %s", terragruntVersion)
}

func (it *VersionCommand) getTerraformVersion() string {
	// Try to get version from terraform CLI first
	if version := it.getVersionFromCLI("terraform"); version != "" {
		return version
	}

	// If terraform CLI is not available, return "not installed"
	return notInstalledVersion
}

func (it *VersionCommand) getTerragruntVersion() string {
	// Try to get version from terragrunt CLI first
	if version := it.getVersionFromCLI("terragrunt"); version != "" {
		return version
	}

	// If terragrunt CLI is not available, return "not installed"
	return notInstalledVersion
}

func (it *VersionCommand) getVersionFromCLI(tool string) string {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, tool, "--version")
	output, err := cmd.Output()
	if err != nil {
		logger.Debugf("Failed to get %s version from CLI: %s", tool, err)
		return ""
	}

	version := strings.TrimSpace(string(output))

	// Extract version number from output
	switch tool {
	case "terraform":
		// Terraform output: "Terraform v1.5.7"
		re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
		matches := re.FindStringSubmatch(version)
		if len(matches) > 1 {
			return matches[1]
		}
	case "terragrunt":
		// Terragrunt output: "terragrunt version v0.50.17"
		re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
		matches := re.FindStringSubmatch(version)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return version
}
