//go:build unit

package commands_test

import (
	"os"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	infrastructure_repositories "github.com/rios0rios0/terra/internal/infrastructure/repositories"

	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRunFromRootCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when valid dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Valid dependencies for creating the command
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()

		// WHEN: Creating a new RunFromRootCommand
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			&repositorydoubles.StubUpgradeShellRepository{},
			interactiveRepository,
		)

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestRunFromRootCommand_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute all steps when normal execution", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with all dependencies and normal arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode"}
		dependencies := []entities.Dependency{
			entitybuilders.NewDependencyBuilder().
				WithName("terraform").
				WithCLI("terraform").
				BuildDependency(),
			entitybuilders.NewDependencyBuilder().
				WithName("terragrunt").
				WithCLI("terragrunt").
				BuildDependency(),
		}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should execute preparation steps but not install
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.Nil(t, installCommand.LastDependencies)

		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.Equal(t, dependencies, formatCommand.LastDependencies)

		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, targetPath, additionalBefore.LastTargetPath)
		assert.Equal(t, arguments, additionalBefore.LastArguments)

		// Should execute terragrunt with upgrade-aware repository (not interactive)
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount)
		assert.Equal(t, "terragrunt", upgradeRepository.LastCommand)
		assert.Equal(t, arguments, upgradeRepository.LastArguments)
		assert.Equal(t, targetPath, upgradeRepository.LastDirectory)
	})

	t.Run("should handle empty arguments when no arguments provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with empty arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty arguments gracefully
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.True(t, additionalBefore.ExecuteCalled, "Should execute additional before command")
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should execute terragrunt command")
		assert.Len(
			t,
			upgradeRepository.LastArguments,
			len(arguments),
			"Should pass arguments with same length",
		)
	})

	t.Run("should handle empty dependencies when no dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with empty dependencies
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should handle empty dependencies gracefully
		assert.False(t, installCommand.ExecuteCalled, "Should not execute install command automatically")
		assert.Nil(t, installCommand.LastDependencies)

		assert.True(t, formatCommand.ExecuteCalled, "Should execute format command")
		assert.Equal(t, dependencies, formatCommand.LastDependencies)

		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should execute terragrunt command")
	})

	t.Run("should pass correct target path when different paths used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with specific target path
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/custom/terraform/modules/vpc"
		arguments := []string{"validate"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should pass correct target path to all components
		assert.Equal(t, targetPath, additionalBefore.LastTargetPath)
		assert.Equal(t, targetPath, upgradeRepository.LastDirectory)
	})

	t.Run("should not use interactive mode when no reply flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command without reply flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := infrastructure_repositories.NewInteractiveShellRepository()
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode", "--out=plan.out"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use normal repository (indirectly tests hasReplyFlag)
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use normal repository")
		assert.Equal(t, arguments, upgradeRepository.LastArguments, "Should pass arguments unchanged")
	})

	t.Run("should inject --non-interactive when boolean --reply maps to yes on plan", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with the deprecated boolean --reply flag on a non-interactive command
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"--reply", "plan", "--detailed-exitcode"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Routes through the upgrade repository with --non-interactive injected
		// (no -auto-approve because plan is not an interactive command).
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use upgrade repository")
		assert.Equal(t,
			[]string{"plan", "--detailed-exitcode", "--non-interactive"},
			upgradeRepository.LastArguments,
			"Should strip --reply and append --non-interactive",
		)
		assert.Equal(t, 0, interactiveRepository.ExecuteWithAnswerCallCount, "PTY path is no longer used")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})

	t.Run("should inject --non-interactive and -auto-approve when --reply=y maps to yes on apply", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with the deprecated --reply=y flag on apply
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"--reply=y", "apply"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Routes through the upgrade repository with --non-interactive AND
		// -auto-approve injected (terraform apply needs -auto-approve to skip its
		// "Enter a value:" prompt; --non-interactive alone does not cover that).
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use upgrade repository")
		assert.Equal(t,
			[]string{"apply", "--non-interactive", "-auto-approve"},
			upgradeRepository.LastArguments,
			"Should strip --reply and append --non-interactive and -auto-approve",
		)
		assert.Equal(t, 0, interactiveRepository.ExecuteWithAnswerCallCount, "PTY path is no longer used")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})

	t.Run("should inject --non-interactive and -auto-approve when --yes is used on apply", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with the new --yes flag on apply
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"apply", "--yes"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Routes through the upgrade repository with native flags injected
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use upgrade repository")
		assert.Equal(t,
			[]string{"apply", "--non-interactive", "-auto-approve"},
			upgradeRepository.LastArguments,
			"Should strip --yes and append --non-interactive and -auto-approve",
		)
		assert.Equal(t, 0, interactiveRepository.ExecuteWithAnswerCallCount, "PTY path is no longer used")
	})

	t.Run("should inject only --non-interactive when --no is used on apply", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with the new --no flag on apply
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		upgradeRepository := &repositorydoubles.StubUpgradeShellRepository{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
		cmd := commands.NewRunFromRootCommand(
			entitybuilders.NewSettingsBuilder().BuildSettings(),
			installCommand,
			formatCommand,
			additionalBefore,
			&commanddoubles.StubParallelState{},
			repository,
			upgradeRepository,
			interactiveRepository,
		)

		targetPath := "/test/path"
		arguments := []string{"apply", "--no"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: --non-interactive is injected but not -auto-approve, so terraform's
		// apply prompt aborts instead of proceeding, matching "no" semantics.
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use upgrade repository")
		assert.Equal(t,
			[]string{"apply", "--non-interactive"},
			upgradeRepository.LastArguments,
			"Should strip --no and append only --non-interactive",
		)
		assert.Equal(t, 0, interactiveRepository.ExecuteWithAnswerCallCount, "PTY path is no longer used")
	})
}

func TestRunFromRootCommand_hasReplyFlag(t *testing.T) {
	t.Parallel()

	// Create command instance for testing
	cmd := commands.NewRunFromRootCommand(
		entitybuilders.NewSettingsBuilder().BuildSettings(),
		&commanddoubles.StubInstallDependencies{},
		&commanddoubles.StubFormatFiles{},
		&commanddoubles.StubRunAdditionalBefore{},
		&commanddoubles.StubParallelState{},
		&repositorydoubles.StubShellRepositoryForRoot{},
		&repositorydoubles.StubUpgradeShellRepository{},
		&repositorydoubles.StubInteractiveShellRepository{},
	)

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{
			name:      "should return true when --reply flag present",
			arguments: []string{"--reply", "plan"},
			expected:  true,
		},
		{
			name:      "should return true when -r flag present",
			arguments: []string{"-r", "apply"},
			expected:  true,
		},
		{
			name:      "should return true when --reply=y flag present",
			arguments: []string{"--reply=y", "destroy"},
			expected:  true,
		},
		{
			name:      "should return true when -r=n flag present",
			arguments: []string{"-r=n", "import"},
			expected:  true,
		},
		{
			name:      "should return false when no reply flag present",
			arguments: []string{"plan", "--detailed-exitcode"},
			expected:  false,
		},
		{
			name:      "should return false when empty arguments",
			arguments: []string{},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cmd.HasReplyFlagPublic(tt.arguments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunFromRootCommand_getReplyValue(t *testing.T) {
	t.Parallel()

	// Create command instance for testing
	cmd := commands.NewRunFromRootCommand(
		entitybuilders.NewSettingsBuilder().BuildSettings(),
		&commanddoubles.StubInstallDependencies{},
		&commanddoubles.StubFormatFiles{},
		&commanddoubles.StubRunAdditionalBefore{},
		&commanddoubles.StubParallelState{},
		&repositorydoubles.StubShellRepositoryForRoot{},
		&repositorydoubles.StubUpgradeShellRepository{},
		&repositorydoubles.StubInteractiveShellRepository{},
	)

	tests := []struct {
		name      string
		arguments []string
		expected  string
	}{
		{
			name:      "should return 'n' for boolean --reply flag",
			arguments: []string{"--reply", "plan"},
			expected:  "n",
		},
		{
			name:      "should return 'n' for boolean -r flag",
			arguments: []string{"-r", "apply"},
			expected:  "n",
		},
		{
			name:      "should return 'y' for --reply=y flag",
			arguments: []string{"--reply=y", "destroy"},
			expected:  "y",
		},
		{
			name:      "should return 'n' for -r=n flag",
			arguments: []string{"-r=n", "import"},
			expected:  "n",
		},
		{
			name:      "should return empty string when no reply flag present",
			arguments: []string{"plan", "--detailed-exitcode"},
			expected:  "",
		},
		{
			name:      "should return custom value for --reply=custom",
			arguments: []string{"--reply=custom", "plan"},
			expected:  "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cmd.GetReplyValuePublic(tt.arguments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunFromRootCommand_removeReplyFlag(t *testing.T) {
	t.Parallel()

	// Create command instance for testing
	cmd := commands.NewRunFromRootCommand(
		entitybuilders.NewSettingsBuilder().BuildSettings(),
		&commanddoubles.StubInstallDependencies{},
		&commanddoubles.StubFormatFiles{},
		&commanddoubles.StubRunAdditionalBefore{},
		&commanddoubles.StubParallelState{},
		&repositorydoubles.StubShellRepositoryForRoot{},
		&repositorydoubles.StubUpgradeShellRepository{},
		&repositorydoubles.StubInteractiveShellRepository{},
	)

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			name:      "should remove --reply flag",
			arguments: []string{"--reply", "plan", "--detailed-exitcode"},
			expected:  []string{"plan", "--detailed-exitcode"},
		},
		{
			name:      "should remove -r flag",
			arguments: []string{"-r", "apply", "--auto-approve"},
			expected:  []string{"apply", "--auto-approve"},
		},
		{
			name:      "should remove --reply=y flag",
			arguments: []string{"--reply=y", "destroy"},
			expected:  []string{"destroy"},
		},
		{
			name:      "should remove -r=n flag",
			arguments: []string{"import", "-r=n", "resource.type.name"},
			expected:  []string{"import", "resource.type.name"},
		},
		{
			name:      "should return unchanged when no reply flag present",
			arguments: []string{"plan", "--detailed-exitcode"},
			expected:  []string{"plan", "--detailed-exitcode"},
		},
		{
			name:      "should return empty slice when only reply flag present",
			arguments: []string{"--reply"},
			expected:  nil, // Go returns nil slice when filtering results in empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cmd.RemoveReplyFlagPublic(tt.arguments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunFromRootCommand_configureCacheEnvironment(t *testing.T) {
	t.Run("should set TG_DOWNLOAD_DIR and TG_PROVIDER_CACHE_DIR when custom paths provided", func(t *testing.T) {
		// given
		t.Setenv("TG_DOWNLOAD_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TF_PLUGIN_CACHE_DIR", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		tempDir := t.TempDir()
		moduleDir := tempDir + "/modules"
		providerDir := tempDir + "/providers"

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(moduleDir).
			WithTerraProviderCacheDir(providerDir).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Equal(t, moduleDir, os.Getenv("TG_DOWNLOAD_DIR"))
		assert.Equal(t, providerDir, os.Getenv("TG_PROVIDER_CACHE_DIR"))
		_, ok := os.LookupEnv("TF_PLUGIN_CACHE_DIR")
		assert.False(t, ok, "TF_PLUGIN_CACHE_DIR should be unset")

		// Verify directories were created
		_, err := os.Stat(moduleDir)
		assert.False(t, os.IsNotExist(err), "Module cache directory should be created")
		_, err = os.Stat(providerDir)
		assert.False(t, os.IsNotExist(err), "Provider cache directory should be created")
	})

	t.Run("should set default cache paths when settings are empty", func(t *testing.T) {
		// given
		t.Setenv("TG_DOWNLOAD_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TF_PLUGIN_CACHE_DIR", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		tgDir := os.Getenv("TG_DOWNLOAD_DIR")
		tgProviderDir := os.Getenv("TG_PROVIDER_CACHE_DIR")
		assert.Contains(t, tgDir, ".cache/terra/modules")
		assert.Contains(t, tgProviderDir, ".cache/terra/providers")
	})

	t.Run("should enable CAS experiment by default when TerraNoCAS is false", func(t *testing.T) {
		// given
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoCAS(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Equal(t, "cas", os.Getenv("TG_EXPERIMENT"))
	})

	t.Run("should not enable CAS experiment when TerraNoCAS is true", func(t *testing.T) {
		// given
		t.Setenv("TG_EXPERIMENT", "cas")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoCAS(true).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Empty(t, os.Getenv("TG_EXPERIMENT"), "TG_EXPERIMENT should not be set when CAS is disabled")
	})

	t.Run("should enable Provider Cache Server by default when TerraNoProviderCache is false", func(t *testing.T) {
		// given
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoProviderCache(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Equal(t, "1", os.Getenv("TG_PROVIDER_CACHE"))
	})

	t.Run("should not enable Provider Cache Server when TerraNoProviderCache is true", func(t *testing.T) {
		// given
		t.Setenv("TG_PROVIDER_CACHE", "1")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoProviderCache(true).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		_, ok := os.LookupEnv("TG_PROVIDER_CACHE")
		assert.False(t, ok, "TG_PROVIDER_CACHE should not be set when Provider Cache is disabled")
	})

	t.Run("should disable auto-provider-cache-dir when Provider Cache Server is enabled", func(t *testing.T) {
		// GIVEN: Provider Cache Server enabled (default) so terra wants terragrunt to
		// honor TG_PROVIDER_CACHE_DIR instead of letting the CAS experiment override it
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TG_NO_AUTO_PROVIDER_CACHE_DIR", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoProviderCache(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// WHEN
		cmd.ConfigureCacheEnvironmentPublic()

		// THEN: terragrunt is told to keep its hands off the provider cache path
		assert.Equal(t, "true", os.Getenv("TG_NO_AUTO_PROVIDER_CACHE_DIR"),
			"TG_NO_AUTO_PROVIDER_CACHE_DIR should be 'true' so TG_PROVIDER_CACHE_DIR is respected")
	})

	t.Run("should unset auto-provider-cache-dir override when Provider Cache Server is disabled", func(t *testing.T) {
		// GIVEN: Pre-existing TG_NO_AUTO_PROVIDER_CACHE_DIR from a prior session plus
		// TerraNoProviderCache=true, so terra must leave the whole provider-cache
		// knob to terragrunt's defaults rather than forcing an override
		t.Setenv("TG_NO_AUTO_PROVIDER_CACHE_DIR", "true")
		t.Setenv("TG_PROVIDER_CACHE", "1")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoProviderCache(true).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// WHEN
		cmd.ConfigureCacheEnvironmentPublic()

		// THEN
		_, ok := os.LookupEnv("TG_NO_AUTO_PROVIDER_CACHE_DIR")
		assert.False(t, ok, "TG_NO_AUTO_PROVIDER_CACHE_DIR should be unset when Provider Cache is disabled")
	})

	t.Run("should enable Partial Parse Config Cache by default when TerraNoPartialParseCache is false", func(t *testing.T) {
		// given
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoPartialParseCache(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Equal(t, "true", os.Getenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE"))
	})

	t.Run("should not enable Partial Parse Config Cache when TerraNoPartialParseCache is true", func(t *testing.T) {
		// given
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "true")
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoPartialParseCache(true).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// when
		cmd.ConfigureCacheEnvironmentPublic()

		// then
		assert.Empty(t, os.Getenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE"),
			"TG_USE_PARTIAL_PARSE_CONFIG_CACHE should not be set when Partial Parse Cache is disabled")
	})
}

func TestRunFromRootCommand_configureCacheEnvironment_allEnabled(t *testing.T) {
	t.Run("should set all cache env vars when CAS and Partial Parse Cache both enabled", func(t *testing.T) {
		// GIVEN: Settings with custom dirs, CAS enabled, and Partial Parse Cache enabled
		t.Setenv("TG_DOWNLOAD_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TG_NO_AUTO_PROVIDER_CACHE_DIR", "")
		t.Setenv("TF_PLUGIN_CACHE_DIR", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		tempDir := t.TempDir()
		moduleDir := tempDir + "/modules"
		providerDir := tempDir + "/providers"

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(moduleDir).
			WithTerraProviderCacheDir(providerDir).
			WithTerraNoCAS(false).
			WithTerraNoPartialParseCache(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// WHEN
		cmd.ConfigureCacheEnvironmentPublic()

		// THEN: All env vars should be set
		assert.Equal(t, moduleDir, os.Getenv("TG_DOWNLOAD_DIR"))
		assert.Equal(t, providerDir, os.Getenv("TG_PROVIDER_CACHE_DIR"))
		assert.Equal(t, "1", os.Getenv("TG_PROVIDER_CACHE"))
		assert.Equal(t, "true", os.Getenv("TG_NO_AUTO_PROVIDER_CACHE_DIR"),
			"TG_NO_AUTO_PROVIDER_CACHE_DIR must be set alongside TG_PROVIDER_CACHE so CAS does not override the cache path")
		_, ok := os.LookupEnv("TF_PLUGIN_CACHE_DIR")
		assert.False(t, ok, "TF_PLUGIN_CACHE_DIR should be unset")
		assert.Equal(t, "cas", os.Getenv("TG_EXPERIMENT"))
		assert.Equal(t, "true", os.Getenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE"))

		// Directories should exist on disk
		_, err := os.Stat(moduleDir)
		require.NoError(t, err, "Module cache directory should be created")
		_, err = os.Stat(providerDir)
		require.NoError(t, err, "Provider cache directory should be created")
	})

	t.Run("should unset all feature env vars when CAS, Provider Cache, and Partial Parse Cache all disabled", func(t *testing.T) {
		// GIVEN: Settings with all features disabled, and pre-existing env vars
		t.Setenv("TG_EXPERIMENT", "cas")
		t.Setenv("TG_PROVIDER_CACHE", "1")
		t.Setenv("TG_NO_AUTO_PROVIDER_CACHE_DIR", "true")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "true")
		t.Setenv("TG_DOWNLOAD_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE_DIR", "")
		t.Setenv("TF_PLUGIN_CACHE_DIR", "")

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(t.TempDir()).
			WithTerraProviderCacheDir(t.TempDir()).
			WithTerraNoCAS(true).
			WithTerraNoProviderCache(true).
			WithTerraNoPartialParseCache(true).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// WHEN
		cmd.ConfigureCacheEnvironmentPublic()

		// THEN: Feature env vars should be unset
		_, ok := os.LookupEnv("TG_EXPERIMENT")
		assert.False(t, ok, "TG_EXPERIMENT should be unset when CAS disabled")

		_, ok = os.LookupEnv("TG_PROVIDER_CACHE")
		assert.False(t, ok, "TG_PROVIDER_CACHE should be unset when Provider Cache disabled")

		_, ok = os.LookupEnv("TG_NO_AUTO_PROVIDER_CACHE_DIR")
		assert.False(t, ok, "TG_NO_AUTO_PROVIDER_CACHE_DIR should be unset when Provider Cache disabled")

		_, ok = os.LookupEnv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE")
		assert.False(t, ok, "TG_USE_PARTIAL_PARSE_CONFIG_CACHE should be unset when Partial Parse Cache disabled")
	})

	t.Run("should set module and provider dirs when directories do not exist yet", func(t *testing.T) {
		// GIVEN: Settings with directories that need to be created
		t.Setenv("TG_DOWNLOAD_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE_DIR", "")
		t.Setenv("TG_PROVIDER_CACHE", "")
		t.Setenv("TF_PLUGIN_CACHE_DIR", "")
		t.Setenv("TG_EXPERIMENT", "")
		t.Setenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE", "")
		tempDir := t.TempDir()
		moduleDir := tempDir + "/deep/nested/modules"
		providerDir := tempDir + "/deep/nested/providers"

		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir(moduleDir).
			WithTerraProviderCacheDir(providerDir).
			WithTerraNoCAS(false).
			WithTerraNoPartialParseCache(false).
			BuildSettings()
		cmd := commands.NewRunFromRootCommand(
			settings,
			&commanddoubles.StubInstallDependencies{},
			&commanddoubles.StubFormatFiles{},
			&commanddoubles.StubRunAdditionalBefore{},
			&commanddoubles.StubParallelState{},
			&repositorydoubles.StubShellRepositoryForRoot{},
			&repositorydoubles.StubUpgradeShellRepository{},
			&repositorydoubles.StubInteractiveShellRepository{},
		)

		// WHEN
		cmd.ConfigureCacheEnvironmentPublic()

		// THEN: Nested directories should be created via MkdirAll
		info, err := os.Stat(moduleDir)
		require.NoError(t, err, "Nested module cache directory should be created")
		assert.True(t, info.IsDir())

		info, err = os.Stat(providerDir)
		require.NoError(t, err, "Nested provider cache directory should be created")
		assert.True(t, info.IsDir())

		assert.Equal(t, moduleDir, os.Getenv("TG_DOWNLOAD_DIR"))
		assert.Equal(t, providerDir, os.Getenv("TG_PROVIDER_CACHE_DIR"))
	})
}

func TestSettings_GetModuleCacheDir(t *testing.T) {
	t.Parallel()

	t.Run("should return custom path when TerraModuleCacheDir is set", func(t *testing.T) {
		t.Parallel()
		// given
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraModuleCacheDir("/custom/modules").
			BuildSettings()

		// when
		dir, err := settings.GetModuleCacheDir()

		// then
		require.NoError(t, err)
		assert.Equal(t, "/custom/modules", dir)
	})

	t.Run("should return default path when TerraModuleCacheDir is empty", func(t *testing.T) {
		t.Parallel()
		// given
		settings := entitybuilders.NewSettingsBuilder().BuildSettings()

		// when
		dir, err := settings.GetModuleCacheDir()

		// then
		require.NoError(t, err)
		assert.Contains(t, dir, ".cache/terra/modules")
	})
}

func TestSettings_GetProviderCacheDir(t *testing.T) {
	t.Parallel()

	t.Run("should return custom path when TerraProviderCacheDir is set", func(t *testing.T) {
		t.Parallel()
		// given
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraProviderCacheDir("/custom/providers").
			BuildSettings()

		// when
		dir, err := settings.GetProviderCacheDir()

		// then
		require.NoError(t, err)
		assert.Equal(t, "/custom/providers", dir)
	})

	t.Run("should return default path when TerraProviderCacheDir is empty", func(t *testing.T) {
		t.Parallel()
		// given
		settings := entitybuilders.NewSettingsBuilder().BuildSettings()

		// when
		dir, err := settings.GetProviderCacheDir()

		// then
		require.NoError(t, err)
		assert.Contains(t, dir, ".cache/terra/providers")
	})
}
