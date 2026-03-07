//go:build unit

package entities_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
)

func TestGetPlatformInfo(t *testing.T) {
	t.Parallel()

	t.Run("should return valid platform info when called", func(t *testing.T) {
		t.Parallel()
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
	})
}

func TestPlatformInfo_GetPlatformString(t *testing.T) {
	t.Parallel()

	t.Run("should return formatted string when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with specific OS and architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithOS("linux").WithArch("amd64").BuildPlatformInfo()

		// WHEN: Getting the platform string
		result := platform.GetPlatformString()

		// THEN: Should return formatted OS_Arch string
		expectedResult := "linux_amd64"
		assert.Equal(t, expectedResult, result,
			"Platform string should be formatted as OS_Arch")
	})
}

func TestPlatformInfo_GetTerraformArchString(t *testing.T) {
	t.Parallel()

	t.Run("should return amd64 when called with amd64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with amd64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("amd64").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return amd64
		assert.Equal(t, "amd64", result,
			"Should return amd64 for amd64 architecture")
	})

	t.Run("should return arm64 when called with arm64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with arm64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("arm64").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return arm64
		assert.Equal(t, "arm64", result,
			"Should return arm64 for arm64 architecture")
	})

	t.Run("should return 386 when called with 386", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with 386 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("386").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return 386
		assert.Equal(t, "386", result,
			"Should return 386 for 386 architecture")
	})

	t.Run("should return arm64 when called with android arm64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_arm64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_arm64").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return arm64 (android converted)
		assert.Equal(t, "arm64", result,
			"Should return arm64 for android_arm64 architecture")
	})

	t.Run("should return amd64 when called with android amd64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_amd64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_amd64").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return amd64 (android converted)
		assert.Equal(t, "amd64", result,
			"Should return amd64 for android_amd64 architecture")
	})

	t.Run("should return 386 when called with android 386", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_386 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_386").BuildPlatformInfo()

		// WHEN: Getting Terraform architecture string
		result := platform.GetTerraformArchString()

		// THEN: Should return 386 (android converted)
		assert.Equal(t, "386", result,
			"Should return 386 for android_386 architecture")
	})
}

func TestPlatformInfo_GetTerragruntArchString(t *testing.T) {
	t.Parallel()

	t.Run("should return amd64 when called with amd64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with amd64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("amd64").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return amd64
		assert.Equal(t, "amd64", result,
			"Should return amd64 for amd64 architecture")
	})

	t.Run("should return arm64 when called with arm64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with arm64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("arm64").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return arm64
		assert.Equal(t, "arm64", result,
			"Should return arm64 for arm64 architecture")
	})

	t.Run("should return 386 when called with 386", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with 386 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("386").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return 386
		assert.Equal(t, "386", result,
			"Should return 386 for 386 architecture")
	})

	t.Run("should return arm64 when called with android arm64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_arm64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_arm64").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return arm64 (android converted)
		assert.Equal(t, "arm64", result,
			"Should return arm64 for android_arm64 architecture")
	})

	t.Run("should return amd64 when called with android amd64", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_amd64 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_amd64").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return amd64 (android converted)
		assert.Equal(t, "amd64", result,
			"Should return amd64 for android_amd64 architecture")
	})

	t.Run("should return 386 when called with android 386", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android_386 architecture
		platform := entitybuilders.NewPlatformInfoBuilder().WithArch("android_386").BuildPlatformInfo()

		// WHEN: Getting Terragrunt architecture string
		result := platform.GetTerragruntArchString()

		// THEN: Should return 386 (android converted)
		assert.Equal(t, "386", result,
			"Should return 386 for android_386 architecture")
	})
}

func TestPlatformInfo_GetOSString(t *testing.T) {
	t.Parallel()

	t.Run("should return linux when called with linux", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with linux OS
		platform := entitybuilders.NewPlatformInfoBuilder().WithOS("linux").BuildPlatformInfo()

		// WHEN: Getting OS string
		result := platform.GetOSString()

		// THEN: Should return linux
		assert.Equal(t, "linux", result,
			"Should return linux for linux OS")
	})

	t.Run("should return windows when called with windows", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with windows OS
		platform := entitybuilders.NewPlatformInfoBuilder().WithOS("windows").BuildPlatformInfo()

		// WHEN: Getting OS string
		result := platform.GetOSString()

		// THEN: Should return windows
		assert.Equal(t, "windows", result,
			"Should return windows for windows OS")
	})

	t.Run("should return darwin when called with darwin", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with darwin OS
		platform := entitybuilders.NewPlatformInfoBuilder().WithOS("darwin").BuildPlatformInfo()

		// WHEN: Getting OS string
		result := platform.GetOSString()

		// THEN: Should return darwin
		assert.Equal(t, "darwin", result,
			"Should return darwin for darwin OS")
	})

	t.Run("should return linux when called with android", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A platform info with android OS
		platform := entitybuilders.NewPlatformInfoBuilder().WithOS("android").BuildPlatformInfo()

		// WHEN: Getting OS string
		result := platform.GetOSString()

		// THEN: Should return linux (android converted to linux)
		assert.Equal(t, "linux", result,
			"Should return linux for android OS (android maps to linux)")
	})
}
