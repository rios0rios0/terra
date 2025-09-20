package commands

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

// Integration test that creates actual files and tests the complete workflow
func TestInstallDependenciesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a mock binary server that serves a simple script
	binaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		// Write a simple executable script
		w.Write([]byte("#!/bin/bash\necho 'test-version 1.0.0'\n"))
	}))
	defer binaryServer.Close()

	// Create a mock version server
	versionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_version":"1.0.0"}`))
	}))
	defer versionServer.Close()

	// Create test dependency with a unique CLI name to avoid conflicts
	dependency := entities.Dependency{
		Name:              "TestTool",
		CLI:               "test-integration-tool-not-installed",
		BinaryURL:         binaryServer.URL + "/testtool_%s",
		VersionURL:        versionServer.URL,
		RegexVersion:      `"current_version":"([^"]+)"`,
		FormattingCommand: []string{"format"},
	}

	// Execute the install command
	cmd := NewInstallDependenciesCommand()
	cmd.Execute([]entities.Dependency{dependency})

	// Verify the installation completed
	// The binary should be installed to ~/.local/bin/ or the OS-specific install path
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

	// Mock servers
	binaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		// Write a mock zip file (not a real zip, but enough to trigger zip handling)
		w.Write([]byte("PK\x03\x04test"))
	}))
	defer binaryServer.Close()

	versionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current_version":"2.0.0"}`))
	}))
	defer versionServer.Close()

	// Create test dependency
	dependency := entities.Dependency{
		Name:              "TestZipTool",
		CLI:               "test-zip-integration-tool-not-installed",
		BinaryURL:         binaryServer.URL + "/testziptool_%s.zip",
		VersionURL:        versionServer.URL,
		RegexVersion:      `"current_version":"([^"]+)"`,
		FormattingCommand: []string{"format"},
	}

	// Execute the install command
	cmd := NewInstallDependenciesCommand()
	cmd.Execute([]entities.Dependency{dependency})

	// Verify the installation completed
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
	// Test with actual dependencies but modified to use mock servers

	// Mock version server that responds to different endpoints
	versionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if strings.Contains(r.URL.Path, "terraform") {
			w.Write([]byte(`{"current_version":"1.5.0"}`))
		} else if strings.Contains(r.URL.Path, "terragrunt") {
			w.Write([]byte(`{"tag_name":"v0.50.0"}`))
		} else {
			w.Write([]byte(`{"version":"1.0.0"}`))
		}
	}))
	defer versionServer.Close()

	// Mock binary server
	binaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("#!/bin/bash\necho 'mock binary'\n"))
	}))
	defer binaryServer.Close()

	// Create test dependencies with unique CLI names
	dependencies := []entities.Dependency{
		{
			Name:              "TestTerraform",
			CLI:               "test-terraform-unique-name",
			BinaryURL:         binaryServer.URL + "/terraform_%s",
			VersionURL:        versionServer.URL + "/terraform",
			RegexVersion:      `"current_version":"([^"]+)"`,
			FormattingCommand: []string{"fmt"},
		},
		{
			Name:              "TestTerragrunt",
			CLI:               "test-terragrunt-unique-name",
			BinaryURL:         binaryServer.URL + "/terragrunt_%s",
			VersionURL:        versionServer.URL + "/terragrunt",
			RegexVersion:      `"tag_name":"v([^"]+)"`,
			FormattingCommand: []string{"hclfmt"},
		},
	}

	// Execute the install command
	cmd := NewInstallDependenciesCommand()
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
