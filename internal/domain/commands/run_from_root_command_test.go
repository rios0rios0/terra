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

	t.Run("should not use interactive mode when no auto answer flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command without auto-answer flag in arguments
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

		// THEN: Should use normal repository (indirectly tests hasAutoAnswerFlag)
		assert.Equal(t, 1, upgradeRepository.ExecuteCallCount, "Should use normal repository")
		assert.Equal(t, arguments, upgradeRepository.LastArguments, "Should pass arguments unchanged")
	})

	t.Run("should use interactive mode with default 'n' when boolean auto-answer flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with boolean auto-answer flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
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

		targetPath := "/test/path"
		arguments := []string{"--auto-answer", "plan", "--detailed-exitcode"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository with default 'n' answer
		assert.Equal(t, 1, interactiveRepository.ExecuteWithAnswerCallCount, "Should use interactive repository")
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
		assert.Equal(t, "n", interactiveRepository.LastAutoAnswer, "Should default to 'n' for boolean flag")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})

	t.Run("should use interactive mode with specified 'y' when auto-answer=y flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with auto-answer=y flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
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

		targetPath := "/test/path"
		arguments := []string{"--auto-answer=y", "apply", "--auto-approve"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository with 'y' answer
		assert.Equal(t, 1, interactiveRepository.ExecuteWithAnswerCallCount, "Should use interactive repository")
		assert.Equal(t, []string{"apply", "--auto-approve"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
		assert.Equal(t, "y", interactiveRepository.LastAutoAnswer, "Should use specified 'y' value")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})

	t.Run("should use interactive mode with specified 'n' when -a=n flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with -a=n flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
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

		targetPath := "/test/path"
		arguments := []string{"-a=n", "plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository with 'n' answer
		assert.Equal(t, 1, interactiveRepository.ExecuteWithAnswerCallCount, "Should use interactive repository")
		assert.Equal(t, []string{"plan"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
		assert.Equal(t, "n", interactiveRepository.LastAutoAnswer, "Should use specified 'n' value")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})

	t.Run("should use interactive mode with default 'n' when short -a flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with short -a flag in arguments
		installCommand := &commanddoubles.StubInstallDependencies{}
		formatCommand := &commanddoubles.StubFormatFiles{}
		additionalBefore := &commanddoubles.StubRunAdditionalBefore{}
		repository := &repositorydoubles.StubShellRepositoryForRoot{}
		interactiveRepository := &repositorydoubles.StubInteractiveShellRepository{}
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

		targetPath := "/test/path"
		arguments := []string{"-a", "destroy"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should use interactive repository with default 'n' answer
		assert.Equal(t, 1, interactiveRepository.ExecuteWithAnswerCallCount, "Should use interactive repository")
		assert.Equal(t, []string{"destroy"}, interactiveRepository.LastArguments, "Should filter out auto-answer flag")
		assert.Equal(t, "n", interactiveRepository.LastAutoAnswer, "Should default to 'n' for short flag")
		assert.Equal(t, 0, repository.ExecuteCallCount, "Should not use normal repository")
	})
}

func TestRunFromRootCommand_hasAutoAnswerFlag(t *testing.T) {
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
			name:      "should return true when --auto-answer flag present",
			arguments: []string{"--auto-answer", "plan"},
			expected:  true,
		},
		{
			name:      "should return true when -a flag present",
			arguments: []string{"-a", "apply"},
			expected:  true,
		},
		{
			name:      "should return true when --auto-answer=y flag present",
			arguments: []string{"--auto-answer=y", "destroy"},
			expected:  true,
		},
		{
			name:      "should return true when -a=n flag present",
			arguments: []string{"-a=n", "import"},
			expected:  true,
		},
		{
			name:      "should return false when no auto-answer flag present",
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
			result := cmd.HasAutoAnswerFlagPublic(tt.arguments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunFromRootCommand_getAutoAnswerValue(t *testing.T) {
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
			name:      "should return 'n' for boolean --auto-answer flag",
			arguments: []string{"--auto-answer", "plan"},
			expected:  "n",
		},
		{
			name:      "should return 'n' for boolean -a flag",
			arguments: []string{"-a", "apply"},
			expected:  "n",
		},
		{
			name:      "should return 'y' for --auto-answer=y flag",
			arguments: []string{"--auto-answer=y", "destroy"},
			expected:  "y",
		},
		{
			name:      "should return 'n' for -a=n flag",
			arguments: []string{"-a=n", "import"},
			expected:  "n",
		},
		{
			name:      "should return empty string when no auto-answer flag present",
			arguments: []string{"plan", "--detailed-exitcode"},
			expected:  "",
		},
		{
			name:      "should return custom value for --auto-answer=custom",
			arguments: []string{"--auto-answer=custom", "plan"},
			expected:  "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cmd.GetAutoAnswerValuePublic(tt.arguments)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunFromRootCommand_removeAutoAnswerFlag(t *testing.T) {
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
			name:      "should remove --auto-answer flag",
			arguments: []string{"--auto-answer", "plan", "--detailed-exitcode"},
			expected:  []string{"plan", "--detailed-exitcode"},
		},
		{
			name:      "should remove -a flag",
			arguments: []string{"-a", "apply", "--auto-approve"},
			expected:  []string{"apply", "--auto-approve"},
		},
		{
			name:      "should remove --auto-answer=y flag",
			arguments: []string{"--auto-answer=y", "destroy"},
			expected:  []string{"destroy"},
		},
		{
			name:      "should remove -a=n flag",
			arguments: []string{"import", "-a=n", "resource.type.name"},
			expected:  []string{"import", "resource.type.name"},
		},
		{
			name:      "should return unchanged when no auto-answer flag present",
			arguments: []string{"plan", "--detailed-exitcode"},
			expected:  []string{"plan", "--detailed-exitcode"},
		},
		{
			name:      "should return empty slice when only auto-answer flag present",
			arguments: []string{"--auto-answer"},
			expected:  nil, // Go returns nil slice when filtering results in empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cmd.RemoveAutoAnswerFlagPublic(tt.arguments)
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
		assert.Empty(t, os.Getenv("TF_PLUGIN_CACHE_DIR"), "TF_PLUGIN_CACHE_DIR should be unset")

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
		assert.Empty(t, os.Getenv("TG_PROVIDER_CACHE"),
			"TG_PROVIDER_CACHE should not be set when Provider Cache is disabled")
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
		assert.Empty(t, os.Getenv("TF_PLUGIN_CACHE_DIR"), "TF_PLUGIN_CACHE_DIR should be unset")
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
		assert.Empty(t, os.Getenv("TG_EXPERIMENT"), "TG_EXPERIMENT should be unset when CAS disabled")
		assert.Empty(t, os.Getenv("TG_PROVIDER_CACHE"),
			"TG_PROVIDER_CACHE should be unset when Provider Cache disabled")
		assert.Empty(t, os.Getenv("TG_USE_PARTIAL_PARSE_CONFIG_CACHE"),
			"TG_USE_PARTIAL_PARSE_CONFIG_CACHE should be unset when Partial Parse Cache disabled")
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
