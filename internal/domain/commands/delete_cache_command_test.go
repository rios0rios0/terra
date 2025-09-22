//go:build unit

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteCacheCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: The NewDeleteCacheCommand constructor is available

		// WHEN: Creating a new delete cache command
		cmd := commands.NewDeleteCacheCommand()

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

//nolint:tparallel // Cannot use t.Parallel() when using t.Chdir()
func TestDeleteCacheCommand_Execute(t *testing.T) {
	// Note: Cannot use t.Parallel() when using t.Chdir()
	t.Run("should delete target directories when called with valid paths", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache command and test directory structure with cache directories
		cmd := commands.NewDeleteCacheCommand()
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		// Create terraform cache directories
		terraformDir := ".terraform"
		nestedTerraformDir := "module1/.terraform"
		require.NoError(t, os.MkdirAll(terraformDir, 0755))
		require.NoError(t, os.MkdirAll(nestedTerraformDir, 0755))
		require.NoError(
			t,
			os.WriteFile(filepath.Join(terraformDir, "test.txt"), []byte("test"), 0644),
		)

		// Create terragrunt cache directories
		terragruntDir := ".terragrunt-cache"
		nestedTerragruntDir := "module2/.terragrunt-cache"
		require.NoError(t, os.MkdirAll(terragruntDir, 0755))
		require.NoError(t, os.MkdirAll(nestedTerragruntDir, 0755))

		// Create directories that should NOT be deleted
		keepDir := "src"
		require.NoError(t, os.MkdirAll(keepDir, 0755))

		// WHEN: Executing the delete command
		cmd.Execute([]string{".terraform", ".terragrunt-cache"})

		// THEN: Should delete target directories but preserve others
		_, err := os.Stat(terraformDir)
		assert.True(t, os.IsNotExist(err), "Terraform directory should be deleted")

		_, err = os.Stat(nestedTerraformDir)
		assert.True(t, os.IsNotExist(err), "Nested terraform directory should be deleted")

		_, err = os.Stat(terragruntDir)
		assert.True(t, os.IsNotExist(err), "Terragrunt directory should be deleted")

		_, err = os.Stat(nestedTerragruntDir)
		assert.True(t, os.IsNotExist(err), "Nested terragrunt directory should be deleted")

		_, err = os.Stat(keepDir)
		assert.False(t, os.IsNotExist(err), "Source directory should be preserved")
	})

	t.Run("should complete without error when empty list provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache command and temporary directory with test directories
		cmd := commands.NewDeleteCacheCommand()
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		testDir := ".terraform"
		require.NoError(t, os.MkdirAll(testDir, 0755))

		// WHEN: Executing with empty directory list
		cmd.Execute([]string{})

		// THEN: Should complete without error and not delete any directories
		_, err := os.Stat(testDir)
		assert.False(
			t,
			os.IsNotExist(err),
			"Directory should not be deleted when no targets specified",
		)
	})

	t.Run(
		"should complete without error when non-existent directories provided",
		func(t *testing.T) {
			t.Parallel()
			// GIVEN: A delete cache command and temporary directory
			cmd := commands.NewDeleteCacheCommand()
			tempDir := t.TempDir()
			t.Chdir(tempDir)

			// WHEN: Executing with non-existent directory names
			// THEN: Should complete without crashing or erroring
			cmd.Execute([]string{".nonexistent", ".alsononexistent"})
		},
	)

	t.Run(
		"should delete only specified directories when selective execution requested",
		func(t *testing.T) {
			t.Parallel()
			// GIVEN: A delete cache command and temporary directory with multiple cache types
			cmd := commands.NewDeleteCacheCommand()
			tempDir := t.TempDir()
			t.Chdir(tempDir)

			terraformDir := "project/.terraform"
			terragruntDir := "project/.terragrunt-cache"

			require.NoError(t, os.MkdirAll(terraformDir, 0755))
			require.NoError(t, os.MkdirAll(terragruntDir, 0755))

			// WHEN: Executing command to delete only terraform directories
			cmd.Execute([]string{".terraform"})

			// THEN: Should delete only terraform directories and preserve terragrunt
			_, err := os.Stat(terraformDir)
			assert.True(t, os.IsNotExist(err), "Terraform directory should be deleted")

			_, err = os.Stat(terragruntDir)
			assert.False(t, os.IsNotExist(err), "Terragrunt directory should be preserved")
		},
	)
}
