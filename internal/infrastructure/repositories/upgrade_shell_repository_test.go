//go:build unit

package repositories_test

import (
	"os"
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
			"should detect please run terraform init suggestion",
			`Backend initialization required, please run "terraform init"`,
		},
		{
			"should detect please run terraform init with single quotes",
			"Backend initialization required, please run 'terraform init'",
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
			"should detect backend initialization required",
			"Backend initialization required, please run \"terraform init\"",
		},
		{
			"should detect state store initialization required",
			"State store initialization required, please run \"terraform init\"",
		},
		{
			"should detect initialization required please run",
			"HCP Terraform or Terraform Enterprise initialization required: please run \"terraform init\"",
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
			"should detect does not match configured version constraint",
			"locked provider registry.terraform.io/hashicorp/aws 4.0.0 does not match configured version constraint ~> 5.0; must use terraform init -upgrade",
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
			"should detect install it automatically by running",
			"You may be able to install it automatically by running:\n  terraform init",
		},
	}

	for _, tt := range upgradeOutputs {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			result := repositories.NeedsUpgradePublic(tt.output)

			// then
			assert.NotEmpty(t, result, "Should detect upgrade need for: %s", tt.output)
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
		{
			"should not detect upgrade for TLS handshake failure from provider",
			"Error: Post \"https://kc.example.io/realms/master/protocol/openid-connect/token\": " +
				"remote error: tls: handshake failure\n\n" +
				"  with module.keycloak_customer_realm.keycloak_realm.customer_realm,\n" +
				"  on .terraform/modules/keycloak_customer_realm/main.tf line 24, " +
				"in resource \"keycloak_realm\" \"customer_realm\":\n" +
				"    24: resource \"keycloak_realm\" \"customer_realm\" {\n",
		},
		{
			"should not detect upgrade for provider API error with registry path",
			"Error: error creating resource: provider registry.terraform.io/hashicorp/aws returned 403",
		},
		{
			"should not detect upgrade for state loading network error",
			"Error loading state: AccessDenied: Access Denied\n\tstatus code: 403",
		},
		{
			"should not detect upgrade for provider timeout",
			"Error: Provider produced inconsistent result after apply\n\n" +
				"  with provider[\"registry.terraform.io/hashicorp/azurerm\"],\n" +
				"  context deadline exceeded (Client.Timeout)",
		},
		{
			"should not detect upgrade for failed to query available provider packages",
			"Error: Failed to query available provider packages\n\n" +
				"Could not retrieve the list of available versions for provider\n" +
				"hashicorp/aws: could not connect to registry.terraform.io",
		},
		{
			"should not detect upgrade for IAM validation error",
			"Error: creating IAM User (dev-1021-bedrock): operation error IAM: CreateUser, " +
				"https response error StatusCode: 400, api error ValidationError: " +
				"1 validation error detected: Value at 'tags.3.member.value' failed to satisfy constraint: " +
				"Member must satisfy regular expression pattern: [\\p{L}\\p{Z}\\p{N}_.:/=+\\-@]*\n\n" +
				"  with module.bedrock[0].aws_iam_user.bedrock,\n" +
				"  on ../../modules/aws_bedrock_claude/main.tf line 12, in resource \"aws_iam_user\" \"bedrock\":\n" +
				"    12: resource \"aws_iam_user\" \"bedrock\" {\n",
		},
		{
			"should not detect upgrade for state lock error with init success in buffer",
			"Terraform has been successfully initialized!\n" +
				"rerun this command to reinitialize your working directory.\n\n" +
				"Error: Error acquiring the state lock\n\n" +
				"Error message: state blob is already locked\n" +
				"Terraform acquires a state lock to protect the state from being written\n" +
				"by multiple users at the same time.",
		},
	}

	for _, tt := range noUpgradeOutputs {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// when
			result := repositories.NeedsUpgradePublic(tt.output)

			// then
			assert.Empty(t, result, "Should NOT detect upgrade need for: %s", tt.output)
		})
	}
}

func TestUpgradeAwareShellRepository_ExecuteCommandWithUpgrade(t *testing.T) {
	t.Parallel()

	t.Run("should succeed when command executes successfully", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN
		err := repo.ExecuteCommandWithUpgrade("echo", []string{"hello"}, dir)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should return error when command fails without upgrade pattern", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN
		err := repo.ExecuteCommandWithUpgrade("false", []string{}, dir)

		// THEN
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform command execution")
	})

	t.Run("should return error when command does not exist", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN
		err := repo.ExecuteCommandWithUpgrade("nonexistent-command-xyz", []string{}, dir)

		// THEN
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform command execution")
	})

	t.Run("should return error when directory does not exist", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()

		// WHEN
		err := repo.ExecuteCommandWithUpgrade("echo", []string{"hello"}, "/nonexistent/directory/path")

		// THEN
		assert.Error(t, err)
	})

	t.Run("should return error when command fails with upgrade pattern but init also fails", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A script that outputs an upgrade pattern and exits with error
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN: Running a command that outputs an upgrade pattern to stderr and fails
		// The "init --upgrade" retry will also fail because the command is not terraform/terragrunt
		err := repo.ExecuteCommandWithUpgrade(
			"sh",
			[]string{"-c", "echo 'Error: terraform init has not been run' >&2; exit 1"},
			dir,
		)

		// THEN: Should return an error with the auto init upgrade failure message
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "auto init --upgrade failed")
	})

	t.Run("should succeed when command succeeds with output", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN: Running a command that produces both stdout and stderr but succeeds
		err := repo.ExecuteCommandWithUpgrade(
			"sh",
			[]string{"-c", "echo 'output to stdout'; echo 'output to stderr' >&2; exit 0"},
			dir,
		)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should pass arguments correctly to the command", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// WHEN: Running a command that uses its arguments
		err := repo.ExecuteCommandWithUpgrade(
			"sh",
			[]string{"-c", "test \"$0\" = 'arg1' && test \"$1\" = 'arg2'", "arg1", "arg2"},
			dir,
		)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should retry and succeed when command fails with upgrade pattern and init succeeds", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A script that uses a marker file to fail on first run with an upgrade
		// pattern, succeed on "init --upgrade", and succeed on retry
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// Create a wrapper script that:
		// - On first call (plan): outputs upgrade pattern and fails
		// - On "init --upgrade" call: creates marker file and succeeds
		// - On retry (plan): sees marker file and succeeds
		scriptContent := `#!/bin/bash
if [ "$1" = "init" ] && [ "$2" = "--upgrade" ]; then
    touch "$PWD/.init_done"
    exit 0
fi
if [ -f "$PWD/.init_done" ]; then
    exit 0
fi
echo "Error: terraform init has not been run" >&2
exit 1
`
		scriptPath := dir + "/fake_terraform.sh"
		require.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o755)) //nolint:gosec // Test file with intentional permissions

		// WHEN: Running the command that will trigger upgrade detection, init, and retry
		err := repo.ExecuteCommandWithUpgrade(scriptPath, []string{"plan"}, dir)

		// THEN: Should succeed after the automatic init --upgrade and retry
		assert.NoError(t, err)

		// Verify the init marker file was created (proving init --upgrade ran)
		_, statErr := os.Stat(dir + "/.init_done")
		assert.NoError(t, statErr, "Init marker file should exist, proving init --upgrade was executed")
	})

	t.Run("should return error when retry fails after successful init upgrade", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A script where init succeeds but the retry also fails
		repo := repositories.NewUpgradeAwareShellRepository()
		dir := t.TempDir()

		// Create a wrapper script that:
		// - On first call (plan): outputs upgrade pattern and fails
		// - On "init --upgrade" call: succeeds
		// - On retry (plan): still fails (simulating a persistent issue)
		scriptContent := `#!/bin/bash
if [ "$1" = "init" ] && [ "$2" = "--upgrade" ]; then
    exit 0
fi
echo "Error: terraform init has not been run" >&2
exit 1
`
		scriptPath := dir + "/fake_terraform.sh"
		require.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o755)) //nolint:gosec // Test file with intentional permissions

		// WHEN: Running the command where both initial and retry commands fail
		err := repo.ExecuteCommandWithUpgrade(scriptPath, []string{"plan"}, dir)

		// THEN: Should return an error from the retry (executePassthrough)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform command execution")
	})
}
