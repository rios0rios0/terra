//go:build integration || unit || test

package entitybuilders //nolint:revive,staticcheck // Test package naming follows established project structure

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// PlatformInfoBuilder helps create test PlatformInfo instances with a fluent interface.
type PlatformInfoBuilder struct {
	*testkit.BaseBuilder
	os   string
	arch string
}

// NewPlatformInfoBuilder creates a new PlatformInfo builder with sensible defaults.
func NewPlatformInfoBuilder() *PlatformInfoBuilder {
	return &PlatformInfoBuilder{
		BaseBuilder: testkit.NewBaseBuilder(),
		os:          "linux",
		arch:        "amd64",
	}
}

// WithOS sets the operating system.
func (b *PlatformInfoBuilder) WithOS(os string) *PlatformInfoBuilder {
	b.os = os
	return b
}

// WithArch sets the architecture.
func (b *PlatformInfoBuilder) WithArch(arch string) *PlatformInfoBuilder {
	b.arch = arch
	return b
}

// Build creates the PlatformInfo (satisfies testkit.Builder interface).
func (b *PlatformInfoBuilder) Build() interface{} {
	return b.BuildPlatformInfo()
}

// BuildPlatformInfo creates the PlatformInfo with a concrete return type for convenience.
func (b *PlatformInfoBuilder) BuildPlatformInfo() entities.PlatformInfo {
	return entities.PlatformInfo{
		OS:   b.os,
		Arch: b.arch,
	}
}

// Reset clears the builder state, allowing it to be reused.
func (b *PlatformInfoBuilder) Reset() testkit.Builder {
	b.BaseBuilder.Reset()
	b.os = "linux"
	b.arch = "amd64"
	return b
}

// Clone creates a deep copy of the PlatformInfoBuilder.
func (b *PlatformInfoBuilder) Clone() testkit.Builder {
	return &PlatformInfoBuilder{
		BaseBuilder: b.BaseBuilder.Clone().(*testkit.BaseBuilder),
		os:          b.os,
		arch:        b.arch,
	}
}
