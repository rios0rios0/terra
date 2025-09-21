package entities

import (
	"fmt"
	"strings"
)

type Dependency struct {
	Name              string
	CLI               string
	VersionURL        string
	BinaryURL         string
	RegexVersion      string
	FormattingCommand []string
}

// GetBinaryURL returns the binary URL with platform information dynamically inserted
func (d *Dependency) GetBinaryURL(version string) string {
	platform := GetPlatformInfo()
	
	// Check if the URL contains platform placeholders (%[2]s and %[3]s for OS and arch)
	// If not, fall back to simple version-only formatting for backward compatibility
	if strings.Contains(d.BinaryURL, "%[2]s") || strings.Contains(d.BinaryURL, "%[3]s") {
		// Format the URL with version, OS, and architecture
		return fmt.Sprintf(d.BinaryURL, version, platform.GetOSString(), platform.GetTerraformArchString())
	}
	
	// Fall back to simple version formatting for backward compatibility
	return fmt.Sprintf(d.BinaryURL, version)
}
