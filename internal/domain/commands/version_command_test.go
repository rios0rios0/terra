package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func TestNewVersionCommand_ShouldCreateInstance_WhenDependenciesProvided(t *testing.T) {
	// GIVEN: Dependencies for version checking
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
	}

	// WHEN: Creating a new version command
	cmd := commands.NewVersionCommand(dependencies)

	// THEN: Should create a valid command instance
	require.NotNil(t, cmd)
}

func TestVersionCommand_ShouldCompleteWithoutPanic_WhenExecuteCalled(t *testing.T) {
	// GIVEN: A version command with dependencies
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
		{
			Name:         "Terragrunt",
			VersionURL:   "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
			RegexVersion: `"tag_name":"v([^"]+)"`,
		},
	}
	cmd := commands.NewVersionCommand(dependencies)

	// WHEN: Executing the version command
	// THEN: Should complete without panicking (verified by not crashing)
	cmd.Execute()
}

func TestVersionCommand_ShouldCompleteWithoutPanic_WhenEmptyDependenciesProvided(t *testing.T) {
	// GIVEN: A version command with empty dependencies
	cmd := commands.NewVersionCommand([]entities.Dependency{})

	// WHEN: Executing the version command
	// THEN: Should complete without panicking (verified by not crashing)
	cmd.Execute()
}

// Note: Additional tests that were testing private methods like getTerraformVersion,
// getVersionFromCLI, and getLatestVersionFromAPI have been removed in accordance with
// the contributing guidelines that state:
// "NEVER test private methods directly. Instead test through public interfaces."
//
// The VersionCommand.Execute method performs network operations and system interactions
// that are better verified through integration tests or manual testing.
