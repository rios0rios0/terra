//go:build unit

package repositories_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUpgradeAwareShellRepository(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// given / when
		repo := repositories.NewUpgradeAwareShellRepository()

		// then
		require.NotNil(t, repo)
	})
}

func TestNeedsUpgrade(t *testing.T) {
	t.Parallel()

	upgradeOutputs := []struct {
		name   string
		output string
	}{
		{
			"should detect terraform init not run",
			"Error: terraform init has not been run for this module",
		},
		{
			"should detect working directory not initialized",
			"Error: Working directory is not initialized",
		},
		{
			"should detect run terraform init suggestion",
			`Error: Could not load backend configuration. Please run "terraform init"`,
		},
		{
			"should detect run terragrunt init suggestion",
			`Error: Please run "terragrunt init" to initialize`,
		},
		{
			"should detect module not installed",
			"Error: Module not installed. Run terraform init to install.",
		},
		{
			"should detect required plugins not installed",
			"Error: Required plugins are not installed",
		},
		{
			"should detect backend configuration changed",
			"Error: Backend configuration changed. Please run terraform init",
		},
		{
			"should detect backend type changed",
			"Error: backend type changed from s3 to gcs",
		},
		{
			"should detect backend configuration has changed",
			"Error: backend configuration has changed since last init",
		},
		{
			"should detect error loading state",
			"Error loading state: unable to load backend",
		},
		{
			"should detect provider version constraint",
			"Error: provider version constraint not satisfied",
		},
		{
			"should detect provider does not satisfy version constraints",
			"Error: Provider doesn't satisfy version constraints",
		},
		{
			"should detect inconsistent dependency lock file",
			"Error: Inconsistent dependency lock file",
		},
		{
			"should detect provider registry issue",
			"Error: provider registry.terraform.io/hashicorp/aws not available",
		},
		{
			"should detect failed to query available provider packages",
			"Error: Failed to query available provider packages",
		},
		{
			"should detect terragrunt upgrade suggestion",
			"You must run 'terragrunt init --upgrade' to continue",
		},
		{
			"should detect terraform init -upgrade suggestion",
			"Please run terraform init -upgrade to resolve",
		},
		{
			"should detect rerun init command suggestion",
			"rerun this command to reinitialize your working directory",
		},
	}

	for _, tt := range upgradeOutputs {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			result := repositories.NeedsUpgradePublic(tt.output)

			// then
			assert.True(t, result, "Should detect upgrade need for: %s", tt.output)
		})
	}

	noUpgradeOutputs := []struct {
		name   string
		output string
	}{
		{
			"should not detect upgrade for normal error",
			"Error: resource not found in state",
		},
		{
			"should not detect upgrade for syntax error",
			"Error: Invalid reference: module.vpc.output",
		},
		{
			"should not detect upgrade for empty output",
			"",
		},
		{
			"should not detect upgrade for successful output",
			"Apply complete! Resources: 1 added, 0 changed, 0 destroyed.",
		},
		{
			"should not detect upgrade for permission error",
			"Error: Access Denied. You don't have permission.",
		},
	}

	for _, tt := range noUpgradeOutputs {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			result := repositories.NeedsUpgradePublic(tt.output)

			// then
			assert.False(t, result, "Should NOT detect upgrade need for: %s", tt.output)
		})
	}
}
