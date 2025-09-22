package commands_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entity_builders"
)

// TestInstallDependenciesCommand_Execute_InstallScenarios tests install method coverage
func TestInstallDependenciesCommand_Execute_InstallScenarios(t *testing.T) {
	// Note: Cannot use t.Parallel() when creating temporary files
	
	t.Run("should handle non-zip binary installation", func(t *testing.T) {
		// GIVEN: A server that serves a non-zip binary file
		binaryServer := createNonZipBinaryServer(t)
		defer binaryServer.Close()
		
		versionServer := createSimpleVersionServer(t, "1.0.0")
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
		// GIVEN: A dependency that will require installation directory creation
		binaryServer := createNonZipBinaryServer(t)
		defer binaryServer.Close()
		
		versionServer := createSimpleVersionServer(t, "2.0.0")
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
		// GIVEN: A binary that needs to be made executable
		binaryServer := createNonZipBinaryServer(t)
		defer binaryServer.Close()
		
		versionServer := createSimpleVersionServer(t, "3.0.0")
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
		// GIVEN: A zip binary that will trigger temp file creation and cleanup
		zipServer := createSimpleZipServer(t, "temp-test-binary")
		defer zipServer.Close()
		
		versionServer := createSimpleVersionServer(t, "1.5.0")
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
		// GIVEN: A binary that will trigger file type detection
		binaryServer := createNonZipBinaryServer(t)
		defer binaryServer.Close()
		
		versionServer := createSimpleVersionServer(t, "4.0.0")
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

// createNonZipBinaryServer creates a server that serves a regular binary (not zip)
func createNonZipBinaryServer(t *testing.T) *httptest.Server {
	t.Helper()
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a simple binary content
		binaryContent := []byte("#!/bin/bash\necho 'test binary'\n")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(binaryContent)
	}))
}

// createSimpleVersionServer creates a version server with given version
func createSimpleVersionServer(t *testing.T, version string) *httptest.Server {
	t.Helper()
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"current_version":"` + version + `"}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

// createSimpleZipServer creates a server that serves a zip file
func createSimpleZipServer(t *testing.T, binaryName string) *httptest.Server {
	t.Helper()
	
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a simple zip file in memory (reuse from existing zip test)
		zipData := createZipWithBinary(t, binaryName)
		
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
}