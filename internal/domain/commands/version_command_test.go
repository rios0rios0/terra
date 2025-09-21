package commands

import (
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestNewVersionCommand(t *testing.T) {
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
	}

	cmd := NewVersionCommand(dependencies)

	if cmd == nil {
		t.Fatal("NewVersionCommand returned nil")
	}

	if len(cmd.dependencies) != len(dependencies) {
		t.Errorf("Expected %d dependencies, got %d", len(dependencies), len(cmd.dependencies))
	}

	if cmd.dependencies[0].Name != dependencies[0].Name {
		t.Errorf(
			"Expected dependency name %s, got %s",
			dependencies[0].Name,
			cmd.dependencies[0].Name,
		)
	}
}

func TestVersionCommand_Execute(t *testing.T) {
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
		{
			Name:         "Terragrunt",
			VersionURL:   "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
			RegexVersion: `"tag_name":"v([^"]+)"`,
		},
	}

	cmd := NewVersionCommand(dependencies)

	// Test that Execute method runs without panicking
	// The actual output is logged via logrus to stderr, which is tested in other test functions
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	cmd.Execute()
}

func TestVersionCommand_getTerraformVersion(t *testing.T) {
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
	}

	cmd := NewVersionCommand(dependencies)
	version := cmd.getTerraformVersion()

	// Since terraform is likely not installed in the test environment,
	// we expect it to return "latest available" from the API fallback
	// or "not installed" if dependency lookup fails
	if version != "latest available" && version != "not installed" {
		// If terraform is actually installed, the version should be a valid version string
		if !strings.Contains(version, ".") {
			t.Errorf(
				"Expected version to be 'latest available', 'not installed', or a valid version, got: %s",
				version,
			)
		}
	}
}

func TestVersionCommand_getTerragruntVersion(t *testing.T) {
	dependencies := []entities.Dependency{
		{
			Name:         "Terragrunt",
			VersionURL:   "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
			RegexVersion: `"tag_name":"v([^"]+)"`,
		},
	}

	cmd := NewVersionCommand(dependencies)
	version := cmd.getTerragruntVersion()

	// Since terragrunt is likely not installed in the test environment,
	// we expect it to return "latest available" from the API fallback
	// or "not installed" if dependency lookup fails
	if version != "latest available" && version != "not installed" {
		// If terragrunt is actually installed, the version should be a valid version string
		if !strings.Contains(version, ".") {
			t.Errorf(
				"Expected version to be 'latest available', 'not installed', or a valid version, got: %s",
				version,
			)
		}
	}
}

func TestVersionCommand_getVersionFromCLI_Terraform(t *testing.T) {
	dependencies := []entities.Dependency{}
	cmd := NewVersionCommand(dependencies)

	// Test with terraform (likely not installed)
	version := cmd.getVersionFromCLI("terraform")

	// If terraform is not installed, should return empty string
	// If it is installed, should return a version string
	if version != "" && !strings.Contains(version, ".") {
		t.Errorf("Expected empty string or valid version, got: %s", version)
	}
}

func TestVersionCommand_getVersionFromCLI_Terragrunt(t *testing.T) {
	dependencies := []entities.Dependency{}
	cmd := NewVersionCommand(dependencies)

	// Test with terragrunt (likely not installed)
	version := cmd.getVersionFromCLI("terragrunt")

	// If terragrunt is not installed, should return empty string
	// If it is installed, should return a version string
	if version != "" && !strings.Contains(version, ".") {
		t.Errorf("Expected empty string or valid version, got: %s", version)
	}
}

func TestVersionCommand_getVersionFromCLI_InvalidTool(t *testing.T) {
	dependencies := []entities.Dependency{}
	cmd := NewVersionCommand(dependencies)

	// Test with invalid tool
	version := cmd.getVersionFromCLI("nonexistent_tool_12345")

	// Should return empty string for non-existent tools
	if version != "" {
		t.Errorf("Expected empty string for non-existent tool, got: %s", version)
	}
}

func TestVersionCommand_getLatestVersionFromAPI(t *testing.T) {
	dependencies := []entities.Dependency{}
	cmd := NewVersionCommand(dependencies)

	// Test the API method
	version := cmd.getLatestVersionFromAPI("https://example.com", `"version":"([^"]+)"`)

	// Currently the method returns "latest available" without making actual network calls
	expectedVersion := "latest available"
	if version != expectedVersion {
		t.Errorf("Expected version %s, got: %s", expectedVersion, version)
	}
}

func TestVersionCommand_WithNoDependencies(t *testing.T) {
	cmd := NewVersionCommand([]entities.Dependency{})

	// Test getTerraformVersion with no dependencies
	terraformVersion := cmd.getTerraformVersion()
	if terraformVersion != "not installed" {
		// If terraform CLI is available, it might return a version
		if terraformVersion != "" && !strings.Contains(terraformVersion, ".") {
			t.Errorf("Expected 'not installed' or valid version, got: %s", terraformVersion)
		}
	}

	// Test getTerragruntVersion with no dependencies
	terragruntVersion := cmd.getTerragruntVersion()
	if terragruntVersion != "not installed" {
		// If terragrunt CLI is available, it might return a version
		if terragruntVersion != "" && !strings.Contains(terragruntVersion, ".") {
			t.Errorf("Expected 'not installed' or valid version, got: %s", terragruntVersion)
		}
	}
}

func TestVersionCommand_NotInstalledBehavior(t *testing.T) {
	// Test that when dependencies are provided but CLI tools are not installed,
	// the version methods return "not installed" instead of "latest available"
	dependencies := []entities.Dependency{
		{
			Name:         "Terraform",
			VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
			RegexVersion: `"current_version":"([^"]+)"`,
		},
		{
			Name:         "Terragrunt",
			VersionURL:   "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
			RegexVersion: `"tag_name":"v([^"]+)"`,
		},
	}

	cmd := NewVersionCommand(dependencies)

	// Test terraform version - should return "not installed" when CLI is not available
	terraformVersion := cmd.getTerraformVersion()
	if terraformVersion != "not installed" {
		// If terraform is actually installed, the version should be a valid version string
		if !strings.Contains(terraformVersion, ".") {
			t.Errorf("Expected 'not installed' or valid version, got: %s", terraformVersion)
		}
	}

	// Test terragrunt version - should return "not installed" when CLI is not available
	terragruntVersion := cmd.getTerragruntVersion()
	if terragruntVersion != "not installed" {
		// If terragrunt is actually installed, the version should be a valid version string
		if !strings.Contains(terragruntVersion, ".") {
			t.Errorf("Expected 'not installed' or valid version, got: %s", terragruntVersion)
		}
	}
}
