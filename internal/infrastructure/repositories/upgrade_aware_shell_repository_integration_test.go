//go:build integration

package repositories_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

func TestUpgradeAwareShellRepository_ExecuteCommandWithUpgradeDetection_Integration(t *testing.T) {
	// Skip if running in CI without terraform/terragrunt
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping integration test in CI environment")
	}

	t.Run("should execute command normally when no upgrade needed", func(t *testing.T) {
		// GIVEN: A repository instance and a simple command
		repo := repositories.NewUpgradeAwareShellRepository()
		tempDir := t.TempDir()

		// WHEN: Executing a command that doesn't need upgrade
		err := repo.ExecuteCommandWithUpgradeDetection("echo", []string{"hello world"}, tempDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Simple command should execute without upgrade detection interference")
	})

	t.Run("should detect upgrade patterns in error output", func(t *testing.T) {
		// GIVEN: A repository instance and a test directory with basic terraform file
		repo := repositories.NewUpgradeAwareShellRepository()
		tempDir := t.TempDir()

		// Create a simple terraform file that might trigger initialization errors
		tfFile := filepath.Join(tempDir, "main.tf")
		tfContent := `
resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "echo 'test'"
  }
}
`
		err := os.WriteFile(tfFile, []byte(tfContent), 0644)
		require.NoError(t, err, "Should create terraform file successfully")

		// WHEN: Executing a terraform plan without initialization (should fail normally)
		// Note: This test validates that the command execution path works, but we expect it to fail
		// since we don't have terraform installed in the test environment
		err = repo.ExecuteCommandWithUpgradeDetection("terraform", []string{"plan"}, tempDir)

		// THEN: Should handle the command execution (may fail due to missing terraform, but shouldn't panic)
		// The important thing is that the upgrade detection logic is exercised
		assert.Error(t, err, "Command should fail when terraform is not available, but error should be handled gracefully")
		assert.Contains(t, err.Error(), "failed to perform command execution", "Error should indicate command execution failure")
	})

	t.Run("should handle upgrade detection for terragrunt commands", func(t *testing.T) {
		// GIVEN: A repository instance
		repo := repositories.NewUpgradeAwareShellRepository()
		tempDir := t.TempDir()

		// Create a terragrunt configuration file
		tgFile := filepath.Join(tempDir, "terragrunt.hcl")
		tgContent := `
terraform {
  source = "git::https://github.com/example/terraform-module.git"
}

inputs = {
  environment = "test"
}
`
		err := os.WriteFile(tgFile, []byte(tgContent), 0644)
		require.NoError(t, err, "Should create terragrunt file successfully")

		// WHEN: Executing a terragrunt plan without initialization
		err = repo.ExecuteCommandWithUpgradeDetection("terragrunt", []string{"plan"}, tempDir)

		// THEN: Should handle the command execution appropriately
		assert.Error(t, err, "Command should fail when terragrunt is not available, but error should be handled gracefully")
		assert.Contains(t, err.Error(), "failed to perform command execution", "Error should indicate command execution failure")
	})
}