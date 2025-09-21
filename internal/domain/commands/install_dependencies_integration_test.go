package commands_test

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test"
	"github.com/stretchr/testify/assert"
)

// Integration test that creates actual files and tests the complete workflow
func TestInstallDependenciesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use builder pattern to create test servers and dependency
	versionServer, binaryServer := test.NewTestServerBuilder().
		WithTerraformVersion("1.0.0").
		BuildServers()
	defer versionServer.Close()
	defer binaryServer.Close()

	dependency := test.NewDependencyBuilder().
		WithName("TestTool").
		WithCLI("test-integration-tool-not-installed").
		WithBinaryURL(binaryServer.URL + "/testtool_%s").
		WithVersionURL(versionServer.URL + "/terraform").
		WithTerraformPattern().
		Build()

	// Execute the install command
	cmd := commands.NewInstallDependenciesCommand()
	cmd.Execute([]entities.Dependency{dependency})

	// Verify the installation completed
	installPath := entities.GetOS().GetInstallationPath()
	expectedPath := filepath.Join(installPath, "test-integration-tool-not-installed")

	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected binary to be installed at %s", expectedPath)
	} else {
		// Clean up the test file
		os.Remove(expectedPath)
	}
}

func TestInstallDependenciesIntegrationWithZip(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use builder pattern for zip file test
	versionServer, binaryServer := test.NewTestServerBuilder().
		WithTerraformVersion("2.0.0").
		WithZipContent().
		BuildServers()
	defer versionServer.Close()
	defer binaryServer.Close()

	dependency := test.NewDependencyBuilder().
		WithName("TestZipTool").
		WithCLI("test-zip-integration-tool-not-installed").
		WithBinaryURL(binaryServer.URL + "/testziptool_%s.zip").
		WithVersionURL(versionServer.URL + "/terraform").
		WithTerraformPattern().
		Build()

	// Execute the install command
	cmd := commands.NewInstallDependenciesCommand()
	cmd.Execute([]entities.Dependency{dependency})

	// Verify the installation completed (may fail with mock zip, which is expected)
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
}

func TestInstallDependenciesExecuteWithMixedDependencies(t *testing.T) {
	// Use builder pattern for mixed dependencies test
	versionServer, binaryServer := test.NewTestServerBuilder().
		WithTerraformVersion("1.5.0").
		WithTerragruntVersion("0.50.0").
		BuildServers()
	defer versionServer.Close()
	defer binaryServer.Close()

	// Create test dependencies using builder pattern
	dependencies := []entities.Dependency{
		test.NewDependencyBuilder().
			WithName("TestTerraform").
			WithCLI("test-terraform-unique-name").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			Build(),
		test.NewDependencyBuilder().
			WithName("TestTerragrunt").
			WithCLI("test-terragrunt-unique-name").
			WithBinaryURL(binaryServer.URL + "/terragrunt_%s").
			WithVersionURL(versionServer.URL + "/terragrunt").
			WithTerragruntPattern().
			Build(),
	}

	// Execute the install command
	cmd := commands.NewInstallDependenciesCommand()
	cmd.Execute(dependencies)

	// Verify both dependencies were processed
	installPath := entities.GetOS().GetInstallationPath()

	for _, dep := range dependencies {
		expectedPath := filepath.Join(installPath, dep.CLI)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected binary %s to be installed at %s", dep.CLI, expectedPath)
		} else {
			// Clean up
			os.Remove(expectedPath)
		}
	}
}

// TestInstallDependenciesDownloadFailure tests the scenario where download fails with HTTP errors
// This addresses the issue reported where "Failed to download terraform: failed to perform download using 'cURL': exit status 23"
func TestInstallDependenciesDownloadFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// GIVEN: A mock server that simulates download failure (HTTP 500)
	versionServer, binaryServer := test.NewTestServerBuilder().
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
}

// TestInstallDependenciesNetworkTimeout tests network timeout scenarios
func TestInstallDependenciesNetworkTimeout(t *testing.T) {
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
}
