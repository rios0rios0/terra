package entities

import (
	"runtime"
	"strings"
)

// PlatformInfo holds OS and architecture information.
type PlatformInfo struct {
	OS   string
	Arch string
}

// GetPlatformInfo returns the current operating system and architecture.
func GetPlatformInfo() PlatformInfo {
	return PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// GetPlatformString returns a formatted platform string in the format OS_ARCH.
func (p PlatformInfo) GetPlatformString() string {
	return p.OS + "_" + p.Arch
}

// GetTerraformArchString returns the architecture string as expected by Terraform releases.
func (p PlatformInfo) GetTerraformArchString() string {
	// Handle Android architecture which includes "android_" prefix
	// Terraform expects standard arch names without the prefix
	if strings.HasPrefix(p.Arch, "android_") {
		return strings.TrimPrefix(p.Arch, "android_")
	}
	// Terraform uses standard Go architecture names
	return p.Arch
}

// GetTerragruntArchString returns the architecture string as expected by Terragrunt releases.
func (p PlatformInfo) GetTerragruntArchString() string {
	// Handle Android architecture which includes "android_" prefix
	// Terragrunt expects standard arch names without the prefix
	if strings.HasPrefix(p.Arch, "android_") {
		return strings.TrimPrefix(p.Arch, "android_")
	}
	// Terragrunt also uses standard Go architecture names
	return p.Arch
}

// GetOSString returns the OS string as expected by dependency releases.
func (p PlatformInfo) GetOSString() string {
	// Android uses Linux binaries for Terraform/Terragrunt
	if p.OS == "android" {
		return "linux"
	}
	// Most tools use the standard Go OS names
	return p.OS
}
