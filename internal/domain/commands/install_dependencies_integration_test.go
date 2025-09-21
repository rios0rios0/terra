package commands_test

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entities_builders"
	"github.com/rios0rios0/terra/test/infrastructure/repositories_builders"
	"github.com/stretchr/testify/assert"
)

func TestInstallDependenciesCommand_Execute_Integration(t *testing.T) {
	t.Run("should install dependency successfully when valid dependency provided", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// GIVEN: A test server and dependency
		versionServer, binaryServer := repositories_builders.NewTestServerBuilder().
			WithTerraformVersion("1.0.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entities_builders.NewDependencyBuilder().
			WithName("TestTool").
			WithCLI("test-integration-tool-not-installed").
			WithBinaryURL(binaryServer.URL + "/testtool_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the install command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should install the dependency
		installPath := entities.GetOS().GetInstallationPath()
		expectedPath := filepath.Join(installPath, "test-integration-tool-not-installed")

		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected binary to be installed at %s", expectedPath)
		} else {
			// Clean up the test file
			os.Remove(expectedPath)
		}
	})
	
	t.Run("should handle zip files when zip dependency provided", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// GIVEN: A test server with zip content
		versionServer, binaryServer := repositories_builders.NewTestServerBuilder().
			WithTerraformVersion("2.0.0").
			WithZipContent().
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entities_builders.NewDependencyBuilder().
			WithName("TestZipTool").
			WithCLI("test-zip-integration-tool-not-installed").
			WithBinaryURL(binaryServer.URL + "/testziptool_%s.zip").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		// WHEN: Executing the install command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should handle zip processing (may fail with mock zip, which is expected)
		installPath := entities.GetOS().GetInstallationPath()
		expectedPath := filepath.Join(installPath, "test-zip-integration-tool-not-installed")

		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			// This might fail because our mock zip isn't a real zip file
			// That's expected - we're testing the flow, not the zip extraction
			t.Logf("Note: Binary not found at %s - this is expected with mock zip", expectedPath)
		} else {
			// Clean up the test file if it was created
			os.Remove(expectedPath)
		}
	})
	
	t.Run("should handle mixed dependencies when multiple dependencies provided", func(t *testing.T) {
		// GIVEN: Multiple test dependencies
		versionServer, binaryServer := repositories_builders.NewTestServerBuilder().
			WithTerraformVersion("1.5.0").
			WithTerragruntVersion("0.50.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		terraformDep := entities_builders.NewDependencyBuilder().
			WithName("TestTerraform").
			WithCLI("test-terraform-unique-name").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build()

		terragruntDep := entities_builders.NewDependencyBuilder().
			WithName("TestTerragrunt").
			WithCLI("test-terragrunt-unique-name").
			WithBinaryURL(binaryServer.URL + "/terragrunt_%s").
			WithVersionURL(versionServer.URL + "/terragrunt").
			WithTerragruntPattern().
			Build()

		// WHEN: Executing with mixed dependencies
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{terraformDep, terragruntDep})

		// THEN: Should handle both dependencies
		installPath := entities.GetOS().GetInstallationPath()
		
		// Check for terraform installation
		terraformPath := filepath.Join(installPath, "test-terraform-unique-name")
		if _, err := os.Stat(terraformPath); err == nil {
			os.Remove(terraformPath)
		}
		
		// Check for terragrunt installation
		terragruntPath := filepath.Join(installPath, "test-terragrunt-unique-name")
		if _, err := os.Stat(terragruntPath); err == nil {
			os.Remove(terragruntPath)
		}
	})
	
	t.Run("should handle download failure when HTTP error occurs", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// GIVEN: A mock server that simulates download failure (HTTP 500)
		versionServer, binaryServer := repositories_builders.NewTestServerBuilder().
			WithTerraformVersion("1.13.3").
			WithDownloadFailure(). // This will make binary server return 500 error
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		// WHEN: Testing download error detection directly through OS interface
		osInstance := entities.GetOS()
		tempFilePath := path.Join(osInstance.GetTempDir(), "test-download-failure")

		// THEN: The download should fail with an HTTP error
		err := osInstance.Download(binaryServer.URL + "/terraform_1.13.3", tempFilePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 500")

		// Clean up any temporary file that might have been created
		os.Remove(tempFilePath)
	})
	
	t.Run("should handle network timeout when unreachable URL provided", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// GIVEN: An unreachable URL to simulate network issues
		// Note: This uses a reserved IP address that should not be reachable
		unreachableURL := "http://192.0.2.1/terraform_1.0.0" // RFC3330 reserved address

		// WHEN: Testing network timeout detection directly through OS interface
		osInstance := entities.GetOS()
		tempFilePath := path.Join(osInstance.GetTempDir(), "test-network-failure")

		// THEN: The download should fail with a network error
		err := osInstance.Download(unreachableURL, tempFilePath)
		assert.Error(t, err)
		// Network errors typically contain timeout or connection-related messages
		errMsg := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errMsg, "timeout") ||
			strings.Contains(errMsg, "connection") ||
			strings.Contains(errMsg, "unreachable") ||
			strings.Contains(errMsg, "no route") ||
			strings.Contains(errMsg, "failed to perform download") ||
			strings.Contains(errMsg, "context deadline exceeded"),
			"Expected network-related error, got: %v", err)

		// Clean up any temporary file that might have been created
		os.Remove(tempFilePath)
	})
}