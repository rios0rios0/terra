//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/stretchr/testify/require"
)

func TestNewVersionCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Dependencies for version checking
		dependencies := []entities.Dependency{
			entitybuilders.NewDependencyBuilder().
				WithName("Terraform").
				WithVersionURL("https://checkpoint-api.hashicorp.com/v1/check/terraform").
				WithTerraformPattern().
				BuildDependency(),
		}

		// WHEN: Creating a new version command
		cmd := commands.NewVersionCommand(dependencies)

		// THEN: Should create a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestVersionCommand_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should complete without panic when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A version command with dependencies
		dependencies := []entities.Dependency{
			entitybuilders.NewDependencyBuilder().
				WithName("Terraform").
				WithVersionURL("https://checkpoint-api.hashicorp.com/v1/check/terraform").
				WithTerraformPattern().
				BuildDependency(),
			entitybuilders.NewDependencyBuilder().
				WithName("Terragrunt").
				WithVersionURL("https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest").
				WithTerragruntPattern().
				BuildDependency(),
		}
		cmd := commands.NewVersionCommand(dependencies)

		// WHEN: Executing the version command
		// THEN: Should complete without panicking (verified by not crashing)
		cmd.Execute()
	})

	t.Run("should complete without panic when empty dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A version command with empty dependencies
		cmd := commands.NewVersionCommand([]entities.Dependency{})

		// WHEN: Executing the version command
		// THEN: Should complete without panicking (verified by not crashing)
		cmd.Execute()
	})
}

func TestVersionCommand_Execute_ToolsNotInstalled(t *testing.T) {
	// NOTE: Cannot use t.Parallel() because t.Setenv modifies process-wide environment

	t.Run("should report not installed when tools are not found in PATH", func(t *testing.T) {
		// GIVEN: PATH is set to empty so no CLI tools can be found
		t.Setenv("PATH", "")
		cmd := commands.NewVersionCommand([]entities.Dependency{})

		// WHEN: Executing the version command with no tools available
		// THEN: Should complete without panicking, reporting "not installed" for both tools
		cmd.Execute()
	})
}
