package entities_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestGetPlatformInfo_ShouldReturnValidPlatformInfo_WhenCalled(t *testing.T) {
	// GIVEN: No preconditions needed

	// WHEN: Getting platform information
	platform := entities.GetPlatformInfo()

	// THEN: Should return valid OS and Arch values matching runtime
	require.NotEmpty(t, platform.OS, "OS should not be empty")
	require.NotEmpty(t, platform.Arch, "Arch should not be empty")
	
	assert.Equal(t, runtime.GOOS, platform.OS,
		"OS should match runtime.GOOS")
	assert.Equal(t, runtime.GOARCH, platform.Arch,
		"Arch should match runtime.GOARCH")
}

func TestPlatformInfo_ShouldReturnFormattedString_WhenGetPlatformStringCalled(t *testing.T) {
	// GIVEN: A platform info with specific OS and architecture
	platform := entities.PlatformInfo{
		OS:   "linux",
		Arch: "amd64",
	}

	// WHEN: Getting the platform string
	result := platform.GetPlatformString()

	// THEN: Should return formatted OS_Arch string
	expectedResult := "linux_amd64"
	assert.Equal(t, expectedResult, result,
		"Platform string should be formatted as OS_Arch")
}

func TestPlatformInfo_ShouldReturnAmd64_WhenTerraformArchStringCalledWithAmd64(t *testing.T) {
	// GIVEN: A platform info with amd64 architecture
	platform := entities.PlatformInfo{Arch: "amd64"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return amd64
	assert.Equal(t, "amd64", result,
		"Should return amd64 for amd64 architecture")
}

func TestPlatformInfo_ShouldReturnArm64_WhenTerraformArchStringCalledWithArm64(t *testing.T) {
	// GIVEN: A platform info with arm64 architecture
	platform := entities.PlatformInfo{Arch: "arm64"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return arm64
	assert.Equal(t, "arm64", result,
		"Should return arm64 for arm64 architecture")
}

func TestPlatformInfo_ShouldReturn386_WhenTerraformArchStringCalledWith386(t *testing.T) {
	// GIVEN: A platform info with 386 architecture
	platform := entities.PlatformInfo{Arch: "386"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return 386
	assert.Equal(t, "386", result,
		"Should return 386 for 386 architecture")
}

func TestPlatformInfo_ShouldReturnArm64_WhenTerraformArchStringCalledWithAndroidArm64(t *testing.T) {
	// GIVEN: A platform info with android_arm64 architecture
	platform := entities.PlatformInfo{Arch: "android_arm64"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return arm64 (android converted)
	assert.Equal(t, "arm64", result,
		"Should return arm64 for android_arm64 architecture")
}

func TestPlatformInfo_ShouldReturnAmd64_WhenTerraformArchStringCalledWithAndroidAmd64(t *testing.T) {
	// GIVEN: A platform info with android_amd64 architecture
	platform := entities.PlatformInfo{Arch: "android_amd64"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return amd64 (android converted)
	assert.Equal(t, "amd64", result,
		"Should return amd64 for android_amd64 architecture")
}

func TestPlatformInfo_ShouldReturn386_WhenTerraformArchStringCalledWithAndroid386(t *testing.T) {
	// GIVEN: A platform info with android_386 architecture
	platform := entities.PlatformInfo{Arch: "android_386"}

	// WHEN: Getting Terraform architecture string
	result := platform.GetTerraformArchString()

	// THEN: Should return 386 (android converted)
	assert.Equal(t, "386", result,
		"Should return 386 for android_386 architecture")
}

func TestPlatformInfo_ShouldReturnAmd64_WhenTerragruntArchStringCalledWithAmd64(t *testing.T) {
	// GIVEN: A platform info with amd64 architecture
	platform := entities.PlatformInfo{Arch: "amd64"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return amd64
	assert.Equal(t, "amd64", result,
		"Should return amd64 for amd64 architecture")
}

func TestPlatformInfo_ShouldReturnArm64_WhenTerragruntArchStringCalledWithArm64(t *testing.T) {
	// GIVEN: A platform info with arm64 architecture
	platform := entities.PlatformInfo{Arch: "arm64"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return arm64
	assert.Equal(t, "arm64", result,
		"Should return arm64 for arm64 architecture")
}

func TestPlatformInfo_ShouldReturn386_WhenTerragruntArchStringCalledWith386(t *testing.T) {
	// GIVEN: A platform info with 386 architecture
	platform := entities.PlatformInfo{Arch: "386"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return 386
	assert.Equal(t, "386", result,
		"Should return 386 for 386 architecture")
}

func TestPlatformInfo_ShouldReturnArm64_WhenTerragruntArchStringCalledWithAndroidArm64(t *testing.T) {
	// GIVEN: A platform info with android_arm64 architecture
	platform := entities.PlatformInfo{Arch: "android_arm64"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return arm64 (android converted)
	assert.Equal(t, "arm64", result,
		"Should return arm64 for android_arm64 architecture")
}

func TestPlatformInfo_ShouldReturnAmd64_WhenTerragruntArchStringCalledWithAndroidAmd64(t *testing.T) {
	// GIVEN: A platform info with android_amd64 architecture
	platform := entities.PlatformInfo{Arch: "android_amd64"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return amd64 (android converted)
	assert.Equal(t, "amd64", result,
		"Should return amd64 for android_amd64 architecture")
}

func TestPlatformInfo_ShouldReturn386_WhenTerragruntArchStringCalledWithAndroid386(t *testing.T) {
	// GIVEN: A platform info with android_386 architecture
	platform := entities.PlatformInfo{Arch: "android_386"}

	// WHEN: Getting Terragrunt architecture string
	result := platform.GetTerragruntArchString()

	// THEN: Should return 386 (android converted)
	assert.Equal(t, "386", result,
		"Should return 386 for android_386 architecture")
}

func TestPlatformInfo_ShouldReturnLinux_WhenGetOSStringCalledWithLinux(t *testing.T) {
	// GIVEN: A platform info with linux OS
	platform := entities.PlatformInfo{OS: "linux"}

	// WHEN: Getting OS string
	result := platform.GetOSString()

	// THEN: Should return linux
	assert.Equal(t, "linux", result,
		"Should return linux for linux OS")
}

func TestPlatformInfo_ShouldReturnWindows_WhenGetOSStringCalledWithWindows(t *testing.T) {
	// GIVEN: A platform info with windows OS
	platform := entities.PlatformInfo{OS: "windows"}

	// WHEN: Getting OS string
	result := platform.GetOSString()

	// THEN: Should return windows
	assert.Equal(t, "windows", result,
		"Should return windows for windows OS")
}

func TestPlatformInfo_ShouldReturnDarwin_WhenGetOSStringCalledWithDarwin(t *testing.T) {
	// GIVEN: A platform info with darwin OS
	platform := entities.PlatformInfo{OS: "darwin"}

	// WHEN: Getting OS string
	result := platform.GetOSString()

	// THEN: Should return darwin
	assert.Equal(t, "darwin", result,
		"Should return darwin for darwin OS")
}

func TestPlatformInfo_ShouldReturnLinux_WhenGetOSStringCalledWithAndroid(t *testing.T) {
	// GIVEN: A platform info with android OS
	platform := entities.PlatformInfo{OS: "android"}

	// WHEN: Getting OS string
	result := platform.GetOSString()

	// THEN: Should return linux (android converted to linux)
	assert.Equal(t, "linux", result,
		"Should return linux for android OS (android maps to linux)")
}
