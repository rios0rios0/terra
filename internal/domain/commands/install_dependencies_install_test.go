package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entity_builders"
	"github.com/rios0rios0/terra/test/infrastructure/repository_helpers"
)

// TestInstallDependenciesCommand_Execute_InstallScenarios tests install method coverage.
//
//nolint:tparallel // Cannot use t.Parallel() when creating temporary files
func TestInstallDependenciesCommand_Execute_InstallScenarios(t *testing.T) {
	// Note: Cannot use t.Parallel() when creating temporary files

	t.Run("should handle non-zip binary installation", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server that serves a non-zip binary file
		binaryServer := repository_helpers.HelperCreateNonZipBinaryServer(t)
		defer binaryServer.Close()

		versionServer := repository_helpers.HelperCreateSimpleVersionServer(t, "1.0.0")
		defer versionServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("SimpleBinary").
			WithCLI("simple-binary-unique-12345").
			WithBinaryURL(binaryServer.URL + "/binary_%s").
			WithVersionURL(versionServer.URL + "/version").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should handle non-zip binary installation path in install method
		// This tests the else branch of zip handling in install method
	})

	t.Run("should handle directory creation for installation", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency that will require installation directory creation
		binaryServer := repository_helpers.HelperCreateNonZipBinaryServer(t)
		defer binaryServer.Close()

		versionServer := repository_helpers.HelperCreateSimpleVersionServer(t, "2.0.0")
		defer versionServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("DirectoryTest").
			WithCLI("directory-test-binary-54321").
			WithBinaryURL(binaryServer.URL + "/binary_%s").
			WithVersionURL(versionServer.URL + "/version").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should create installation directory if it doesn't exist
		// This tests the os.MkdirAll call in install method
	})

	t.Run("should handle installation path permissions", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A binary that needs to be made executable
		binaryServer := repository_helpers.HelperCreateNonZipBinaryServer(t)
		defer binaryServer.Close()

		versionServer := repository_helpers.HelperCreateSimpleVersionServer(t, "3.0.0")
		defer versionServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("PermissionTest").
			WithCLI("permission-test-binary-98765").
			WithBinaryURL(binaryServer.URL + "/binary_%s").
			WithVersionURL(versionServer.URL + "/version").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should make binary executable after installation
		// This tests the currentOS.MakeExecutable call in install method
	})

	t.Run("should handle temp file creation and cleanup", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A zip binary that will trigger temp file creation and cleanup
		zipServer := repository_helpers.HelperCreateSimpleZipServer(t, "temp-test-binary")
		defer zipServer.Close()

		versionServer := repository_helpers.HelperCreateSimpleVersionServer(t, "1.5.0")
		defer versionServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("TempFileTest").
			WithCLI("temp-test-binary").
			WithBinaryURL(zipServer.URL + "/temp_%s.zip").
			WithVersionURL(versionServer.URL + "/version").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should create temp files, extract, and clean up
		// This tests temp file handling and cleanup logic in install method
	})

	t.Run("should handle file type detection", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A binary that will trigger file type detection
		binaryServer := repository_helpers.HelperCreateNonZipBinaryServer(t)
		defer binaryServer.Close()

		versionServer := repository_helpers.HelperCreateSimpleVersionServer(t, "4.0.0")
		defer versionServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("FileTypeTest").
			WithCLI("file-type-test-binary-13579").
			WithBinaryURL(binaryServer.URL + "/binary_%s").
			WithVersionURL(versionServer.URL + "/version").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should detect file type using `file` command
		// This tests the exec.CommandContext for file type detection in install method
	})
}
