//go:build unit

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteCacheCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: The NewDeleteCacheCommand constructor is available

		// WHEN: Creating a new delete cache command
		cmd := commands.NewDeleteCacheCommand(&entities.Settings{})

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestDeleteCacheCommand_Execute(t *testing.T) {
	// Note: Cannot use t.Parallel() when using t.Chdir()
	t.Run("should delete target directories when called with valid paths", func(t *testing.T) {
		// GIVEN: A delete cache command and test directory structure with cache directories
		cmd := commands.NewDeleteCacheCommand(&entities.Settings{})
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		// Create terraform cache directories
		terraformDir := ".terraform"
		nestedTerraformDir := "module1/.terraform"
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(terraformDir, 0755))
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(nestedTerraformDir, 0755))
		require.NoError(
			t,
			os.WriteFile(filepath.Join(terraformDir, "test.txt"), []byte("test"), 0644),
		)

		// Create terragrunt cache directories
		terragruntDir := ".terragrunt-cache"
		nestedTerragruntDir := "module2/.terragrunt-cache"
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(terragruntDir, 0755))
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(nestedTerragruntDir, 0755))

		// Create directories that should NOT be deleted
		keepDir := "src"
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(keepDir, 0755))

		// WHEN: Executing the delete command
		cmd.Execute([]string{".terraform", ".terragrunt-cache"}, false)

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
		// GIVEN: A delete cache command and temporary directory with test directories
		cmd := commands.NewDeleteCacheCommand(&entities.Settings{})
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		testDir := ".terraform"
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(testDir, 0755))

		// WHEN: Executing with empty directory list
		cmd.Execute([]string{}, false)

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
			// GIVEN: A delete cache command and temporary directory
			cmd := commands.NewDeleteCacheCommand(&entities.Settings{})
			tempDir := t.TempDir()
			t.Chdir(tempDir)

			// WHEN: Executing with non-existent directory names
			// THEN: Should complete without crashing or erroring
			cmd.Execute([]string{".nonexistent", ".alsononexistent"}, false)
		},
	)

	t.Run(
		"should delete only specified directories when selective execution requested",
		func(t *testing.T) {
			// GIVEN: A delete cache command and temporary directory with multiple cache types
			cmd := commands.NewDeleteCacheCommand(&entities.Settings{})
			tempDir := t.TempDir()
			t.Chdir(tempDir)

			terraformDir := "project/.terraform"
			terragruntDir := "project/.terragrunt-cache"

			// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
			require.NoError(t, os.MkdirAll(terraformDir, 0755))
			// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
			require.NoError(t, os.MkdirAll(terragruntDir, 0755))

			// WHEN: Executing command to delete only terraform directories
			cmd.Execute([]string{".terraform"}, false)

			// THEN: Should delete only terraform directories and preserve terragrunt
			_, err := os.Stat(terraformDir)
			assert.True(t, os.IsNotExist(err), "Terraform directory should be deleted")

			_, err = os.Stat(terragruntDir)
			assert.False(t, os.IsNotExist(err), "Terragrunt directory should be preserved")
		},
	)

	t.Run("should remove global cache directories when global flag is true", func(t *testing.T) {
		// GIVEN: A delete cache command with custom cache directories that exist
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		moduleCache := filepath.Join(tempDir, "custom-modules")
		providerCache := filepath.Join(tempDir, "custom-providers")

		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(moduleCache, 0755))
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(providerCache, 0755))

		settings := &entities.Settings{
			TerraModuleCacheDir:   moduleCache,
			TerraProviderCacheDir: providerCache,
		}
		cmd := commands.NewDeleteCacheCommand(settings)

		// WHEN: Executing with global=true
		cmd.Execute([]string{}, true)

		// THEN: Should remove both global cache directories
		_, err := os.Stat(moduleCache)
		assert.True(t, os.IsNotExist(err), "Module cache directory should be deleted")

		_, err = os.Stat(providerCache)
		assert.True(t, os.IsNotExist(err), "Provider cache directory should be deleted")
	})

	t.Run("should skip non-existent global cache directories gracefully when global flag is true", func(t *testing.T) {
		// GIVEN: A delete cache command with cache directories that do not exist
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		settings := &entities.Settings{
			TerraModuleCacheDir:   filepath.Join(tempDir, "nonexistent-modules"),
			TerraProviderCacheDir: filepath.Join(tempDir, "nonexistent-providers"),
		}
		cmd := commands.NewDeleteCacheCommand(settings)

		// WHEN: Executing with global=true
		// THEN: Should complete without error
		cmd.Execute([]string{}, true)
	})

	t.Run("should respect custom cache paths from settings when global flag is true", func(t *testing.T) {
		// GIVEN: A delete cache command with custom settings and only module cache exists
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		moduleCache := filepath.Join(tempDir, "my-modules")
		providerCache := filepath.Join(tempDir, "my-providers")

		// Only create the module cache, not the provider cache
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(moduleCache, 0755))

		settings := &entities.Settings{
			TerraModuleCacheDir:   moduleCache,
			TerraProviderCacheDir: providerCache,
		}
		cmd := commands.NewDeleteCacheCommand(settings)

		// WHEN: Executing with global=true
		cmd.Execute([]string{}, true)

		// THEN: Should delete existing module cache and skip non-existent provider cache
		_, err := os.Stat(moduleCache)
		assert.True(t, os.IsNotExist(err), "Module cache directory should be deleted")
	})

	t.Run("should not remove global cache when global flag is false", func(t *testing.T) {
		// GIVEN: A delete cache command with custom cache directories that exist
		tempDir := t.TempDir()
		t.Chdir(tempDir)

		moduleCache := filepath.Join(tempDir, "keep-modules")
		providerCache := filepath.Join(tempDir, "keep-providers")

		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(moduleCache, 0755))
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(providerCache, 0755))

		settings := &entities.Settings{
			TerraModuleCacheDir:   moduleCache,
			TerraProviderCacheDir: providerCache,
		}
		cmd := commands.NewDeleteCacheCommand(settings)

		// WHEN: Executing with global=false
		cmd.Execute([]string{}, false)

		// THEN: Should preserve both global cache directories
		_, err := os.Stat(moduleCache)
		assert.False(t, os.IsNotExist(err), "Module cache directory should be preserved")

		_, err = os.Stat(providerCache)
		assert.False(t, os.IsNotExist(err), "Provider cache directory should be preserved")
	})
}
