package entities_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

// TestDependency_ShouldGeneratePlatformSpecificURL_WhenPlatformPlaceholdersProvided demonstrates
// the new BDD testing guidelines with testify framework
func TestDependency_ShouldGeneratePlatformSpecificURL_WhenPlatformPlaceholdersProvided(t *testing.T) {
	// GIVEN: A dependency with platform-specific URL template
	dependency := entities.Dependency{
		Name:      "Terraform",
		CLI:       "terraform",
		BinaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
	}
	version := "1.5.0"

	// WHEN: GetBinaryURL is called with a version
	result := dependency.GetBinaryURL(version)

	// THEN: URL should contain version, OS, and architecture
	require.NotEmpty(t, result, "Binary URL should not be empty")
	
	expectedSubstrings := []string{
		"https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0",
		runtime.GOOS,
		runtime.GOARCH,
		".zip",
	}

	for _, substring := range expectedSubstrings {
		assert.Contains(t, result, substring, 
			"URL should contain platform-specific information: %s", substring)
	}

	// Verify placeholders are replaced
	placeholders := []string{"%[1]s", "%[2]s", "%[3]s"}
	for _, placeholder := range placeholders {
		assert.NotContains(t, result, placeholder,
			"URL should not contain unreplaced placeholder: %s", placeholder)
	}
}

// TestDependency_ShouldUseFallbackFormat_WhenNoPlatformPlaceholdersFound demonstrates
// testing backward compatibility behavior
func TestDependency_ShouldUseFallbackFormat_WhenNoPlatformPlaceholdersFound(t *testing.T) {
	// GIVEN: A dependency with simple version-only URL template (backward compatibility)
	dependency := entities.Dependency{
		Name:      "SimpleTool",
		CLI:       "simple",
		BinaryURL: "https://example.com/tool_%s.tar.gz",
	}
	version := "2.1.0"

	// WHEN: GetBinaryURL is called
	result := dependency.GetBinaryURL(version)

	// THEN: URL should use simple version formatting
	expectedURL := "https://example.com/tool_2.1.0.tar.gz"
	assert.Equal(t, expectedURL, result,
		"Should use simple version formatting for backward compatibility")
	
	// And should not contain platform-specific information
	assert.NotContains(t, result, runtime.GOOS,
		"Backward compatible URLs should not include OS information")
	assert.NotContains(t, result, runtime.GOARCH,
		"Backward compatible URLs should not include architecture information")
}

// TestDependency_ShouldHandleEmptyVersion_WhenCalledWithEmptyString demonstrates
// edge case testing following BDD guidelines
func TestDependency_ShouldHandleEmptyVersion_WhenCalledWithEmptyString(t *testing.T) {
	// GIVEN: A dependency with URL template and empty version
	dependency := entities.Dependency{
		Name:      "TestTool",
		BinaryURL: "https://example.com/tool_%[1]s_%[2]s_%[3]s",
	}
	emptyVersion := ""

	// WHEN: GetBinaryURL is called with empty version
	result := dependency.GetBinaryURL(emptyVersion)

	// THEN: URL should still be generated with platform info but empty version
	require.NotEmpty(t, result, "Should still generate URL even with empty version")
	
	platform := entities.GetPlatformInfo()
	expectedURL := "https://example.com/tool__" + platform.GetOSString() + "_" + platform.GetTerraformArchString()
	
	assert.Equal(t, expectedURL, result,
		"Should handle empty version gracefully while preserving platform information")
}

// TestDependency_ShouldDetectPlatformPlaceholders_WhenOnlyPartialPlaceholdersPresent demonstrates
// testing the platform detection logic
func TestDependency_ShouldDetectPlatformPlaceholders_WhenOnlyPartialPlaceholdersPresent(t *testing.T) {
	testCases := []struct {
		name        string
		binaryURL   string
		description string
	}{
		{
			name:        "OS placeholder only",
			binaryURL:   "https://example.com/tool_%[1]s_%[2]s",
			description: "should detect platform format when only OS placeholder present",
		},
		{
			name:        "Architecture placeholder only", 
			binaryURL:   "https://example.com/tool_%[1]s_%[3]s",
			description: "should detect platform format when only arch placeholder present",
		},
	}

	version := "1.0.0"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// GIVEN: A dependency with partial platform placeholders
			dependency := entities.Dependency{
				Name:      "TestTool",
				BinaryURL: tc.binaryURL,
			}

			// WHEN: GetBinaryURL is called
			result := dependency.GetBinaryURL(version)

			// THEN: Should use platform-specific formatting
			platform := entities.GetPlatformInfo()
			assert.Contains(t, result, version, "Should contain version")
			
			if strings.Contains(tc.binaryURL, "%[2]s") {
				assert.Contains(t, result, platform.GetOSString(), 
					"Should contain OS when OS placeholder present")
			}
			
			if strings.Contains(tc.binaryURL, "%[3]s") {
				assert.Contains(t, result, platform.GetTerraformArchString(),
					"Should contain architecture when arch placeholder present")
			}
		})
	}
}