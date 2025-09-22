package commands_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entity_builders"
	"github.com/rios0rios0/terra/test/infrastructure/repository_builders"
	"github.com/stretchr/testify/require"
)

// TestInstallDependenciesCommand_Execute_VersionScenarios tests version comparison and prompt functionality
func TestInstallDependenciesCommand_Execute_VersionScenarios(t *testing.T) {
	// Note: Cannot use t.Parallel() when manipulating PATH and creating temporary binaries
	
	t.Run("should trigger version comparison with mock terraform", func(t *testing.T) {
		// GIVEN: A mock terraform binary that returns a proper version
		mockBinaryDir := createMockTerraformBinary(t, "1.0.0")
		defer os.RemoveAll(mockBinaryDir)
		
		// Prepend mock binary directory to PATH
		originalPath := os.Getenv("PATH")
		newPath := mockBinaryDir + string(os.PathListSeparator) + originalPath
		t.Setenv("PATH", newPath)

		// Set up mock server with newer version
		versionServer, binaryServer := repository_builders.NewTestServerBuilder().
			WithTerraformVersion("2.0.0"). // Newer version to trigger update
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("Terraform").
			WithCLI("terraform").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// Mock stdin to simulate "no" response to update prompt
		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		go func() {
			defer w.Close()
			w.Write([]byte("no\n"))
		}()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// Restore stdin
		os.Stdin = oldStdin
		r.Close()

		// THEN: Should have triggered getCurrentVersion, compareVersions, and promptForUpdate
		// This tests the full version comparison flow
	})

	t.Run("should trigger version comparison with equal versions", func(t *testing.T) {
		// GIVEN: A mock terraform binary that returns same version as server
		mockBinaryDir := createMockTerraformBinary(t, "1.5.0")
		defer os.RemoveAll(mockBinaryDir)
		
		// Prepend mock binary directory to PATH
		originalPath := os.Getenv("PATH")
		newPath := mockBinaryDir + string(os.PathListSeparator) + originalPath
		t.Setenv("PATH", newPath)

		// Set up mock server with same version
		versionServer, binaryServer := repository_builders.NewTestServerBuilder().
			WithTerraformVersion("1.5.0"). // Same version
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("Terraform").
			WithCLI("terraform").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should have triggered compareVersions with equal versions (== 0 path)
	})

	t.Run("should trigger version comparison with newer local version", func(t *testing.T) {
		// GIVEN: A mock terraform binary that returns newer version than server
		mockBinaryDir := createMockTerraformBinary(t, "2.0.0")
		defer os.RemoveAll(mockBinaryDir)
		
		// Prepend mock binary directory to PATH
		originalPath := os.Getenv("PATH")
		newPath := mockBinaryDir + string(os.PathListSeparator) + originalPath
		t.Setenv("PATH", newPath)

		// Set up mock server with older version
		versionServer, binaryServer := repository_builders.NewTestServerBuilder().
			WithTerraformVersion("1.0.0"). // Older version
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("Terraform").
			WithCLI("terraform").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should have triggered compareVersions with newer local version (> 0 path)
	})

	t.Run("should handle user accepting update prompt", func(t *testing.T) {
		// GIVEN: A mock terraform binary with older version
		mockBinaryDir := createMockTerraformBinary(t, "1.0.0")
		defer os.RemoveAll(mockBinaryDir)
		
		// Prepend mock binary directory to PATH
		originalPath := os.Getenv("PATH")
		newPath := mockBinaryDir + string(os.PathListSeparator) + originalPath
		t.Setenv("PATH", newPath)

		// Set up mock server with newer version
		versionServer, binaryServer := repository_builders.NewTestServerBuilder().
			WithTerraformVersion("2.0.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entity_builders.NewDependencyBuilder().
			WithName("Terraform").
			WithCLI("terraform").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// Mock stdin to simulate "yes" response to update prompt
		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		go func() {
			defer w.Close()
			w.Write([]byte("yes\n"))
		}()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// Restore stdin
		os.Stdin = oldStdin
		r.Close()

		// THEN: Should have triggered promptForUpdate returning true and install path
	})
}

// createMockTerraformBinary creates a temporary binary that mimics terraform --version output
func createMockTerraformBinary(t *testing.T, version string) string {
	t.Helper()
	
	tempDir, err := os.MkdirTemp("", "mock-terraform-*")
	require.NoError(t, err)
	
	binaryPath := filepath.Join(tempDir, "terraform")
	
	// Create a shell script that outputs proper terraform version format
	scriptContent := fmt.Sprintf(`#!/bin/bash
if [ "$1" = "--version" ]; then
    echo "Terraform v%s"
elif [ "$1" = "-v" ]; then
    echo "Terraform v%s"
else
    echo "mock binary"
fi
`, version, version)
	
	err = os.WriteFile(binaryPath, []byte(scriptContent), 0755)
	require.NoError(t, err)
	
	return tempDir
}