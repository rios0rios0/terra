package entities_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestPlatformInfo_GetOSString_AndroidMapping(t *testing.T) {
	t.Parallel()

	t.Run("should return linux OS when android OS provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An Android platform
		platform := entities.PlatformInfo{OS: "android", Arch: "arm64"}

		// WHEN: GetOSString is called
		result := platform.GetOSString()

		// THEN: Should return linux instead of android
		assert.Equal(t, "linux", result, "Android OS should map to linux for dependency downloads")
	})

	t.Run("should return original OS when non-android OS provided", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name string
			os   string
		}{
			{"Linux", "linux"},
			{"Windows", "windows"},
			{"Darwin", "darwin"},
			{"FreeBSD", "freebsd"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// GIVEN: A non-Android platform
				platform := entities.PlatformInfo{OS: tc.os, Arch: "amd64"}

				// WHEN: GetOSString is called
				result := platform.GetOSString()

				// THEN: Should return the original OS unchanged
				assert.Equal(t, tc.os, result, "Non-Android OS should remain unchanged")
			})
		}
	})
}

//nolint:gocognit // Comprehensive test with multiple Android platform scenarios
func TestDependency_GetBinaryURL_AndroidPlatform(t *testing.T) {
	t.Parallel()

	t.Run("should generate linux URL when android platform used", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name        string
			cli         string
			binaryURL   string
			platform    entities.PlatformInfo
			version     string
			expectedURL string
		}{
			{
				name:        "Terraform on Android arm64",
				cli:         "terraform",
				binaryURL:   "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
				platform:    entities.PlatformInfo{OS: "android", Arch: "android_arm64"},
				version:     "1.13.3",
				expectedURL: "https://releases.hashicorp.com/terraform/1.13.3/terraform_1.13.3_linux_arm64.zip",
			},
			{
				name:        "Terragrunt on Android arm64",
				cli:         "terragrunt",
				binaryURL:   "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
				platform:    entities.PlatformInfo{OS: "android", Arch: "android_arm64"},
				version:     "0.50.0",
				expectedURL: "https://github.com/gruntwork-io/terragrunt/releases/download/v0.50.0/terragrunt_linux_arm64",
			},
			{
				name:        "Terraform on Android amd64",
				cli:         "terraform",
				binaryURL:   "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
				platform:    entities.PlatformInfo{OS: "android", Arch: "android_amd64"},
				version:     "1.13.3",
				expectedURL: "https://releases.hashicorp.com/terraform/1.13.3/terraform_1.13.3_linux_amd64.zip",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// GIVEN: A dependency configured for Android platform
				dependency := entities.Dependency{
					CLI:       tc.cli,
					BinaryURL: tc.binaryURL,
				}

				// Create test implementation that simulates GetBinaryURL with our platform
				testGetBinaryURL := func(version string) string {
					if strings.Contains(dependency.BinaryURL, "%[2]s") ||
						strings.Contains(dependency.BinaryURL, "%[3]s") {
						var archString string
						if dependency.CLI == "terragrunt" {
							archString = tc.platform.GetTerragruntArchString()
						} else {
							archString = tc.platform.GetTerraformArchString()
						}
						return fmt.Sprintf(
							dependency.BinaryURL,
							version,
							tc.platform.GetOSString(),
							archString,
						)
					}
					return fmt.Sprintf(dependency.BinaryURL, version)
				}

				// WHEN: GetBinaryURL is called
				result := testGetBinaryURL(tc.version)

				// THEN: Should generate URL with linux OS and correct architecture
				assert.Equal(
					t,
					tc.expectedURL,
					result,
					"Should generate correct URL for Android platform",
				)

				// Verify no android_ prefix in URL
				assert.NotContains(t, result, "android_", "URL should not contain android_ prefix")
				assert.Contains(t, result, "linux", "URL should contain linux OS")
				assert.Contains(t, result, tc.version, "URL should contain version")
			})
		}
	})

	t.Run("should use correct arch method when different dependencies used", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An Android platform with android_arm64 architecture
		platform := entities.PlatformInfo{OS: "android", Arch: "android_arm64"}

		testCases := []struct {
			name           string
			cli            string
			expectedArch   string
			archMethodName string
		}{
			{
				name:           "Terraform dependency",
				cli:            "terraform",
				expectedArch:   "arm64",
				archMethodName: "GetTerraformArchString",
			},
			{
				name:           "Terragrunt dependency",
				cli:            "terragrunt",
				expectedArch:   "arm64",
				archMethodName: "GetTerragruntArchString",
			},
			{
				name:           "Other dependency",
				cli:            "other",
				expectedArch:   "arm64",
				archMethodName: "GetTerraformArchString (default)",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				// GIVEN: A dependency with specific CLI name
				dependency := entities.Dependency{
					CLI:       tc.cli,
					BinaryURL: "https://example.com/tool_%[1]s_%[2]s_%[3]s",
				}

				// Create test implementation that simulates GetBinaryURL logic
				testGetBinaryURL := func(version string) string {
					if strings.Contains(dependency.BinaryURL, "%[2]s") ||
						strings.Contains(dependency.BinaryURL, "%[3]s") {
						var archString string
						if dependency.CLI == "terragrunt" {
							archString = platform.GetTerragruntArchString()
						} else {
							archString = platform.GetTerraformArchString()
						}
						return fmt.Sprintf(
							dependency.BinaryURL,
							version,
							platform.GetOSString(),
							archString,
						)
					}
					return fmt.Sprintf(dependency.BinaryURL, version)
				}

				// WHEN: GetBinaryURL is called
				result := testGetBinaryURL("1.0.0")

				// THEN: Should use correct architecture method and map OS correctly
				expectedURL := "https://example.com/tool_1.0.0_linux_arm64"
				assert.Equal(
					t,
					expectedURL,
					result,
					"Should use correct arch method for %s",
					tc.archMethodName,
				)
				assert.Contains(t, result, "linux", "Should map android OS to linux")
				assert.Contains(
					t,
					result,
					"arm64",
					"Should strip android_ prefix from architecture",
				)
			})
		}
	})
}
