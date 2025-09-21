package entities

import (
	"runtime"
	"testing"
)

func TestGetPlatformInfo(t *testing.T) {
	platform := GetPlatformInfo()

	// Test that we get valid OS and Arch values
	if platform.OS == "" {
		t.Error("Expected non-empty OS, got empty string")
	}

	if platform.Arch == "" {
		t.Error("Expected non-empty Arch, got empty string")
	}

	// Test that OS matches runtime.GOOS
	if platform.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, platform.OS)
	}

	// Test that Arch matches runtime.GOARCH
	if platform.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %s, got %s", runtime.GOARCH, platform.Arch)
	}
}

func TestPlatformInfo_GetPlatformString(t *testing.T) {
	platform := PlatformInfo{
		OS:   "linux",
		Arch: "amd64",
	}

	expected := "linux_amd64"
	result := platform.GetPlatformString()

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestPlatformInfo_GetTerraformArchString(t *testing.T) {
	testCases := []struct {
		name     string
		arch     string
		expected string
	}{
		{"amd64", "amd64", "amd64"},
		{"arm64", "arm64", "arm64"},
		{"386", "386", "386"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			platform := PlatformInfo{Arch: tc.arch}
			result := platform.GetTerraformArchString()

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestPlatformInfo_GetTerragruntArchString(t *testing.T) {
	testCases := []struct {
		name     string
		arch     string
		expected string
	}{
		{"amd64", "amd64", "amd64"},
		{"arm64", "arm64", "arm64"},
		{"386", "386", "386"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			platform := PlatformInfo{Arch: tc.arch}
			result := platform.GetTerragruntArchString()

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestPlatformInfo_GetOSString(t *testing.T) {
	testCases := []struct {
		name     string
		os       string
		expected string
	}{
		{"linux", "linux", "linux"},
		{"windows", "windows", "windows"},
		{"darwin", "darwin", "darwin"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			platform := PlatformInfo{OS: tc.os}
			result := platform.GetOSString()

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}