//nolint:testpackage // Testing private functions and fields requires same package
package entities

import (
	"runtime"
	"strings"
	"testing"
)

func TestDependency_GetBinaryURL_Terraform(t *testing.T) {
	dependency := Dependency{
		Name:      "Terraform",
		CLI:       "terraform",
		BinaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
	}

	version := "1.5.0"
	result := dependency.GetBinaryURL(version)

	expectedSubstrings := []string{
		"https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0",
		runtime.GOOS,
		runtime.GOARCH,
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
	dependency := Dependency{
		Name:      "Terragrunt",
		CLI:       "terragrunt",
		BinaryURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
	}

	version := "0.50.0"
	result := dependency.GetBinaryURL(version)

	expectedSubstrings := []string{
		"https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt",
		runtime.GOOS,
		runtime.GOARCH,
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
	dependency := Dependency{
		BinaryURL: "https://example.com/%[1]s/%[2]s_%[3]s",
	}

	// Test with current platform
	version := "1.0.0"
	result := dependency.GetBinaryURL(version)

	platform := GetPlatformInfo()
	expectedURL := "https://example.com/1.0.0/" + platform.GetOSString() + "_" + platform.GetTerraformArchString()

	if result != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, result)
	}
}

func TestDependency_GetBinaryURL_BackwardCompatibility(t *testing.T) {
	// Test backward compatibility with simple version-only URLs
	dependency := Dependency{
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
				platform := GetPlatformInfo()
				return "https://example.com/tool_1.0.0_" + platform.GetOSString()
			},
		},
		{
			name:      "platform format with arch placeholder",
			binaryURL: "https://example.com/tool_%[1]s_%[3]s",
			version:   "1.0.0",
			expected: func() string {
				platform := GetPlatformInfo()
				return "https://example.com/tool_1.0.0_" + platform.GetTerraformArchString()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dependency := Dependency{BinaryURL: tc.binaryURL}
			result := dependency.GetBinaryURL(tc.version)
			expected := tc.expected()

			if result != expected {
				t.Errorf("Expected %s, got %s", expected, result)
			}
		})
	}
}
