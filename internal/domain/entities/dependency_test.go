package entities_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestDependency_GetBinaryURL_Terraform(t *testing.T) {
	dependency := entities.Dependency{
		Name:      "Terraform",
		CLI:       "terraform",
		BinaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
	}

	version := "1.5.0"
	result := dependency.GetBinaryURL(version)

	// Get the platform-mapped values for proper assertion on all platforms including Android
	platform := entities.GetPlatformInfo()

	expectedSubstrings := []string{
		"https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0",
		platform.GetOSString(),            // Use platform-mapped OS (android -> linux)
		platform.GetTerraformArchString(), // Use platform-mapped arch (android_arm64 -> arm64)
		".zip",
	}

	for _, substring := range expectedSubstrings {
		if !strings.Contains(result, substring) {
			t.Errorf("Expected URL to contain '%s', got: %s", substring, result)
		}
	}

	// Verify the URL doesn't contain the placeholders anymore
	placeholders := []string{"%[1]s", "%[2]s", "%[3]s"}
	for _, placeholder := range placeholders {
		if strings.Contains(result, placeholder) {
			t.Errorf("URL should not contain placeholder '%s', got: %s", placeholder, result)
		}
	}
}

func TestDependency_GetBinaryURL_Terragrunt(t *testing.T) {
	dependency := entities.Dependency{
		Name:      "Terragrunt",
		CLI:       "terragrunt",
		BinaryURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
	}

	version := "0.50.0"
	result := dependency.GetBinaryURL(version)

	// Get the platform-mapped values for proper assertion on all platforms including Android
	platform := entities.GetPlatformInfo()

	expectedSubstrings := []string{
		"https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt",
		platform.GetOSString(),             // Use platform-mapped OS (android -> linux)
		platform.GetTerragruntArchString(), // Use platform-mapped arch (android_arm64 -> arm64)
	}

	for _, substring := range expectedSubstrings {
		if !strings.Contains(result, substring) {
			t.Errorf("Expected URL to contain '%s', got: %s", substring, result)
		}
	}

	// Verify the URL doesn't contain the placeholders anymore
	placeholders := []string{"%[2]s", "%[3]s"}
	for _, placeholder := range placeholders {
		if strings.Contains(result, placeholder) {
			t.Errorf("URL should not contain placeholder '%s', got: %s", placeholder, result)
		}
	}
}

func TestDependency_GetBinaryURL_PlatformVariations(t *testing.T) {
	dependency := entities.Dependency{
		BinaryURL: "https://example.com/%[1]s/%[2]s_%[3]s",
	}

	// Test with current platform
	version := "1.0.0"
	result := dependency.GetBinaryURL(version)

	platform := entities.GetPlatformInfo()
	expectedURL := "https://example.com/1.0.0/" + platform.GetOSString() + "_" + platform.GetTerraformArchString()

	if result != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, result)
	}
}

func TestDependency_GetBinaryURL_BackwardCompatibility(t *testing.T) {
	// Test backward compatibility with simple version-only URLs
	dependency := entities.Dependency{
		BinaryURL: "https://example.com/tool_%s",
	}

	version := "1.0.0"
	result := dependency.GetBinaryURL(version)
	expectedURL := "https://example.com/tool_1.0.0"

	if result != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, result)
	}
}

func TestDependency_GetBinaryURL_MixedFormats(t *testing.T) {
	testCases := []struct {
		name      string
		binaryURL string
		version   string
		expected  func() string // Function to generate expected result
	}{
		{
			name:      "simple version format",
			binaryURL: "https://example.com/tool_%s",
			version:   "1.0.0",
			expected:  func() string { return "https://example.com/tool_1.0.0" },
		},
		{
			name:      "platform format with OS placeholder",
			binaryURL: "https://example.com/tool_%[1]s_%[2]s",
			version:   "1.0.0",
			expected: func() string {
				platform := entities.GetPlatformInfo()
				return "https://example.com/tool_1.0.0_" + platform.GetOSString()
			},
		},
		{
			name:      "platform format with arch placeholder",
			binaryURL: "https://example.com/tool_%[1]s_%[3]s",
			version:   "1.0.0",
			expected: func() string {
				platform := entities.GetPlatformInfo()
				return "https://example.com/tool_1.0.0_" + platform.GetTerraformArchString()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dependency := entities.Dependency{BinaryURL: tc.binaryURL}
			result := dependency.GetBinaryURL(tc.version)
			expected := tc.expected()

			if result != expected {
				t.Errorf("Expected %s, got %s", expected, result)
			}
		})
	}
}

func TestDependency_GetBinaryURL_AndroidArchitecture(t *testing.T) {
	testCases := []struct {
		name             string
		binaryURL        string
		platform         entities.PlatformInfo
		version          string
		expectedContains []string
	}{
		{
			name:      "Terraform with android_arm64",
			binaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
			platform:  entities.PlatformInfo{OS: "android", Arch: "android_arm64"},
			version:   "1.5.0",
			expectedContains: []string{
				"https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0_linux_arm64.zip",
			},
		},
		{
			name:      "Terragrunt with android_arm64",
			binaryURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
			platform:  entities.PlatformInfo{OS: "android", Arch: "android_arm64"},
			version:   "0.50.0",
			expectedContains: []string{
				"https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt_linux_arm64",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dependency := entities.Dependency{BinaryURL: tc.binaryURL}

			// Create a test dependency that uses our test platform
			testGetBinaryURL := func(version string) string {
				// Simulate the logic from GetBinaryURL but with our test platform
				if strings.Contains(dependency.BinaryURL, "%[2]s") || strings.Contains(dependency.BinaryURL, "%[3]s") {
					// Determine which arch method to use based on dependency type
					var archString string
					if strings.Contains(dependency.BinaryURL, "terragrunt") {
						archString = tc.platform.GetTerragruntArchString()
					} else {
						archString = tc.platform.GetTerraformArchString()
					}
					return fmt.Sprintf(dependency.BinaryURL, version, tc.platform.GetOSString(), archString)
				}
				return fmt.Sprintf(dependency.BinaryURL, version)
			}

			result := testGetBinaryURL(tc.version)

			for _, expectedSubstring := range tc.expectedContains {
				if !strings.Contains(result, expectedSubstring) {
					t.Errorf("Expected URL to contain '%s', got: %s", expectedSubstring, result)
				}
			}
		})
	}
}
