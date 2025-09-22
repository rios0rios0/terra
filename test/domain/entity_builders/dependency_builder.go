//go:build integration || unit || test

package entity_builders //nolint:revive,staticcheck // Test package naming follows established project structure

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
)

// DependencyBuilder helps create test dependencies with a fluent interface.
type DependencyBuilder struct {
	name              string
	cli               string
	binaryURL         string
	versionURL        string
	regexVersion      string
	formattingCommand []string
}

// NewDependencyBuilder creates a new dependency builder.
func NewDependencyBuilder() *DependencyBuilder {
	return &DependencyBuilder{
		name:              "TestDependency",
		cli:               "test-dependency",
		regexVersion:      `"version":"([^"]+)"`,
		formattingCommand: []string{"format"},
	}
}

// WithName sets the dependency name.
func (b *DependencyBuilder) WithName(name string) *DependencyBuilder {
	b.name = name
	return b
}

// WithCLI sets the CLI name.
func (b *DependencyBuilder) WithCLI(cli string) *DependencyBuilder {
	b.cli = cli
	return b
}

// WithBinaryURL sets the binary URL.
func (b *DependencyBuilder) WithBinaryURL(url string) *DependencyBuilder {
	b.binaryURL = url
	return b
}

// WithVersionURL sets the version URL.
func (b *DependencyBuilder) WithVersionURL(url string) *DependencyBuilder {
	b.versionURL = url
	return b
}

// WithRegexVersion sets the regex version pattern.
func (b *DependencyBuilder) WithRegexVersion(regex string) *DependencyBuilder {
	b.regexVersion = regex
	return b
}

// WithTerraformPattern sets up Terraform-like patterns.
func (b *DependencyBuilder) WithTerraformPattern() *DependencyBuilder {
	return b.WithRegexVersion(`"current_version":"([^"]+)"`)
}

// WithTerragruntPattern sets up Terragrunt-like patterns.
func (b *DependencyBuilder) WithTerragruntPattern() *DependencyBuilder {
	return b.WithRegexVersion(`"tag_name":"v([^"]+)"`)
}

// Build creates the dependency.
func (b *DependencyBuilder) Build() entities.Dependency {
	return entities.Dependency{
		Name:              b.name,
		CLI:               b.cli,
		BinaryURL:         b.binaryURL,
		VersionURL:        b.versionURL,
		RegexVersion:      b.regexVersion,
		FormattingCommand: b.formattingCommand,
	}
}
