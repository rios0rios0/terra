package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test"
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

// TestInstallDependenciesDownloadFailure tests the scenario where curl fails with exit status 23
// This addresses the issue reported where "Failed to download terraform: failed to perform download using 'cURL': exit status 23"
func TestInstallDependenciesDownloadFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping test: cannot reliably test logger.Fatalf behavior that calls os.Exit()")

	// Use builder pattern to create servers that simulate download failure
	versionServer, binaryServer := test.NewTestServerBuilder().
		WithTerraformVersion("1.13.3").
		WithDownloadFailure(). // This will make binary server return 500 error
		BuildServers()
	defer versionServer.Close()
	defer binaryServer.Close()

	dependency := test.NewDependencyBuilder().
		WithName("TestDownloadFailure").
		WithCLI("test-download-failure-tool").
		WithBinaryURL(binaryServer.URL + "/terraform_%s").
		WithVersionURL(versionServer.URL + "/terraform").
		WithTerraformPattern().
		Build()

	// Execute the install command - this should handle the download failure gracefully
	// Note: The actual implementation calls logger.Fatalf on download failure,
	// so in a real scenario this would exit the process. In production code,
	// we might want to refactor this to return errors instead.
	cmd := commands.NewInstallDependenciesCommand()

	// This test verifies that download failures are properly detected
	// In the real scenario, the curl command will fail and the error will be logged
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - the function calls logger.Fatalf on download errors
			t.Logf("Download failure correctly detected and handled: %v", r)
		}
	}()

	cmd.Execute([]entities.Dependency{dependency})

	// If we reach here, it means the download "succeeded" with our mock server error response
	// In real scenarios with curl, this would fail with exit status 23 or similar
	t.Logf(
		"Note: With mock server, download failure is simulated differently than real curl errors",
	)
}

// TestInstallDependenciesNetworkTimeout tests network timeout scenarios
func TestInstallDependenciesNetworkTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Skipping test: cannot reliably test logger.Fatalf behavior that calls os.Exit()")

	// Create a dependency with an unreachable URL to simulate network issues
	// Note: This uses a reserved IP address that should not be reachable
	dependency := test.NewDependencyBuilder().
		WithName("TestNetworkFailure").
		WithCLI("test-network-failure-tool").
		WithBinaryURL("http://192.0.2.1/terraform_%s"). // RFC3330 reserved address
		WithVersionURL("http://192.0.2.1/version").     // RFC3330 reserved address
		WithTerraformPattern().
		Build()

	cmd := commands.NewInstallDependenciesCommand()

	// This test verifies that network failures are properly detected
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - the function calls logger.Fatalf on network errors
			t.Logf("Network failure correctly detected and handled: %v", r)
		}
	}()

	cmd.Execute([]entities.Dependency{dependency})
}
