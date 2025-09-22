package entities_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestDependency_GetBinaryURL(t *testing.T) {
	t.Parallel()
	
	t.Run("should generate terraform URL when terraform dependency provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Terraform dependency with platform-specific URL template
		dependency := entities.Dependency{
			Name:      "Terraform",
			CLI:       "terraform",
			BinaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
		}
		version := "1.5.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: URL should contain all expected components
		require.NotEmpty(t, result, "Binary URL should not be empty")

		platform := entities.GetPlatformInfo()
		expectedSubstrings := []string{
			"https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0",
			platform.GetOSString(),            // Use GetOSString() instead of runtime.GOOS for Android compatibility
			platform.GetTerraformArchString(), // Use GetTerraformArchString() for consistent arch handling
			".zip",
		}

		for _, substring := range expectedSubstrings {
			assert.Contains(t, result, substring,
				"URL should contain expected substring: %s", substring)
		}

		// Verify the URL doesn't contain the placeholders anymore
		placeholders := []string{"%[1]s", "%[2]s", "%[3]s"}
		for _, placeholder := range placeholders {
			assert.NotContains(t, result, placeholder,
				"URL should not contain unreplaced placeholder: %s", placeholder)
		}
	})
	
	t.Run("should generate terragrunt URL when terragrunt dependency provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Terragrunt dependency with platform-specific URL template
		dependency := entities.Dependency{
			Name:      "Terragrunt",
			CLI:       "terragrunt",
			BinaryURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
		}
		version := "0.50.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: URL should contain all expected components
		require.NotEmpty(t, result, "Binary URL should not be empty")

		platform := entities.GetPlatformInfo()
		expectedSubstrings := []string{
			"https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt",
			platform.GetOSString(),             // Use GetOSString() instead of runtime.GOOS for Android compatibility
			platform.GetTerragruntArchString(), // Use GetTerragruntArchString() for consistent arch handling
		}

		for _, substring := range expectedSubstrings {
			assert.Contains(t, result, substring,
				"URL should contain expected substring: %s", substring)
		}

		// Verify the URL doesn't contain the placeholders anymore
		placeholders := []string{"%[2]s", "%[3]s"}
		for _, placeholder := range placeholders {
			assert.NotContains(t, result, placeholder,
				"URL should not contain unreplaced placeholder: %s", placeholder)
		}
	})
	
	t.Run("should generate platform specific URL when platform variations used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency with platform placeholders
		dependency := entities.Dependency{
			BinaryURL: "https://example.com/%[1]s/%[2]s_%[3]s",
		}
		version := "1.0.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: URL should match expected platform-specific format
		platform := entities.GetPlatformInfo()
		expectedURL := "https://example.com/1.0.0/" + platform.GetOSString() + "_" + platform.GetTerraformArchString()

		assert.Equal(t, expectedURL, result,
			"URL should match expected platform-specific format")
	})
	
	t.Run("should use simple version format when backward compatibility required", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency with simple version-only URL template (backward compatibility)
		dependency := entities.Dependency{
			BinaryURL: "https://example.com/tool_%s",
		}
		version := "1.0.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: URL should use simple version formatting
		expectedURL := "https://example.com/tool_1.0.0"
		assert.Equal(t, expectedURL, result,
			"URL should use simple version formatting for backward compatibility")
	})
	
	t.Run("should use simple version format when simple version format used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency with simple version format
		dependency := entities.Dependency{BinaryURL: "https://example.com/tool_%s"}
		version := "1.0.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: Should return URL with simple version format
		expectedURL := "https://example.com/tool_1.0.0"
		assert.Equal(t, expectedURL, result,
			"Should generate URL with simple version format")
	})
	
	t.Run("should include OS information when platform format with OS placeholder used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency with platform format containing OS placeholder
		dependency := entities.Dependency{BinaryURL: "https://example.com/tool_%[1]s_%[2]s"}
		version := "1.0.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: Should return URL with OS information
		platform := entities.GetPlatformInfo()
		expectedURL := "https://example.com/tool_1.0.0_" + platform.GetOSString()
		assert.Equal(t, expectedURL, result,
			"Should generate URL with OS information")
	})
	
	t.Run("should include arch information when platform format with arch placeholder used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A dependency with platform format containing architecture placeholder
		dependency := entities.Dependency{BinaryURL: "https://example.com/tool_%[1]s_%[3]s"}
		version := "1.0.0"

		// WHEN: Getting the binary URL
		result := dependency.GetBinaryURL(version)

		// THEN: Should return URL with architecture information
		platform := entities.GetPlatformInfo()
		expectedURL := "https://example.com/tool_1.0.0_" + platform.GetTerraformArchString()
		assert.Equal(t, expectedURL, result,
			"Should generate URL with architecture information")
	})
	
	t.Run("should generate linux arm64 URL when terraform with android arm64 used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Terraform dependency and android_arm64 platform
		dependency := entities.Dependency{BinaryURL: "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip"}
		testPlatform := entities.PlatformInfo{OS: "android", Arch: "android_arm64"}
		version := "1.5.0"

		// WHEN: Getting the binary URL (simulating android platform conversion to linux)
		var result string
		if strings.Contains(dependency.BinaryURL, "%[2]s") || strings.Contains(dependency.BinaryURL, "%[3]s") {
			archString := testPlatform.GetTerraformArchString()
			result = fmt.Sprintf(dependency.BinaryURL, version, testPlatform.GetOSString(), archString)
		} else {
			result = fmt.Sprintf(dependency.BinaryURL, version)
		}

		// THEN: Should generate URL with linux_arm64 (android converted to linux)
		expectedURL := "https://releases.hashicorp.com/terraform/1.5.0/terraform_1.5.0_linux_arm64.zip"
		assert.Contains(t, result, expectedURL,
			"Should generate URL with linux_arm64 for android_arm64 platform")
	})
	
	t.Run("should generate linux arm64 URL when terragrunt with android arm64 used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Terragrunt dependency and android_arm64 platform
		dependency := entities.Dependency{BinaryURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s"}
		testPlatform := entities.PlatformInfo{OS: "android", Arch: "android_arm64"}
		version := "0.50.0"

		// WHEN: Getting the binary URL (simulating android platform conversion to linux)
		var result string
		if strings.Contains(dependency.BinaryURL, "%[2]s") || strings.Contains(dependency.BinaryURL, "%[3]s") {
			archString := testPlatform.GetTerragruntArchString()
			result = fmt.Sprintf(dependency.BinaryURL, version, testPlatform.GetOSString(), archString)
		} else {
			result = fmt.Sprintf(dependency.BinaryURL, version)
		}

		// THEN: Should generate URL with linux_arm64 (android converted to linux)
		expectedURL := "https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt_linux_arm64"
		assert.Contains(t, result, expectedURL,
			"Should generate URL with linux_arm64 for android_arm64 platform")
	})
}
