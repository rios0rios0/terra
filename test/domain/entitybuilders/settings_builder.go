//go:build integration || unit || test

package entitybuilders //nolint:revive,staticcheck // Test package naming follows established project structure

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// SettingsBuilder helps create test Settings instances with a fluent interface.
type SettingsBuilder struct {
	*testkit.BaseBuilder
	terraCloud               string
	terraTerraformWorkspace  string
	terraAwsRoleArn          string
	terraAzureSubscriptionID string
	terraModuleCacheDir      string
	terraProviderCacheDir    string
	terraNoCAS               bool
	terraNoProviderCache     bool
	terraNoPartialParseCache bool
	terraNoWorkspace         bool
}

// NewSettingsBuilder creates a new Settings builder with empty defaults.
func NewSettingsBuilder() *SettingsBuilder {
	return &SettingsBuilder{
		BaseBuilder: testkit.NewBaseBuilder(),
	}
}

// WithTerraCloud sets the cloud provider.
func (b *SettingsBuilder) WithTerraCloud(cloud string) *SettingsBuilder {
	b.terraCloud = cloud
	return b
}

// WithTerraTerraformWorkspace sets the Terraform workspace.
func (b *SettingsBuilder) WithTerraTerraformWorkspace(workspace string) *SettingsBuilder {
	b.terraTerraformWorkspace = workspace
	return b
}

// WithTerraAwsRoleArn sets the AWS role ARN.
func (b *SettingsBuilder) WithTerraAwsRoleArn(roleArn string) *SettingsBuilder {
	b.terraAwsRoleArn = roleArn
	return b
}

// WithTerraAzureSubscriptionID sets the Azure subscription ID.
func (b *SettingsBuilder) WithTerraAzureSubscriptionID(subscriptionID string) *SettingsBuilder {
	b.terraAzureSubscriptionID = subscriptionID
	return b
}

// WithTerraModuleCacheDir sets the module cache directory.
func (b *SettingsBuilder) WithTerraModuleCacheDir(dir string) *SettingsBuilder {
	b.terraModuleCacheDir = dir
	return b
}

// WithTerraProviderCacheDir sets the provider cache directory.
func (b *SettingsBuilder) WithTerraProviderCacheDir(dir string) *SettingsBuilder {
	b.terraProviderCacheDir = dir
	return b
}

// WithTerraNoCAS sets the no-CAS flag.
func (b *SettingsBuilder) WithTerraNoCAS(noCAS bool) *SettingsBuilder {
	b.terraNoCAS = noCAS
	return b
}

// WithTerraNoProviderCache sets the no-provider-cache flag.
func (b *SettingsBuilder) WithTerraNoProviderCache(noProviderCache bool) *SettingsBuilder {
	b.terraNoProviderCache = noProviderCache
	return b
}

// WithTerraNoPartialParseCache sets the no-partial-parse-cache flag.
func (b *SettingsBuilder) WithTerraNoPartialParseCache(noPartialParseCache bool) *SettingsBuilder {
	b.terraNoPartialParseCache = noPartialParseCache
	return b
}

// WithTerraNoWorkspace sets the no-workspace flag.
func (b *SettingsBuilder) WithTerraNoWorkspace(noWorkspace bool) *SettingsBuilder {
	b.terraNoWorkspace = noWorkspace
	return b
}

// Build creates the Settings (satisfies testkit.Builder interface).
func (b *SettingsBuilder) Build() interface{} {
	return b.BuildSettings()
}

// BuildSettings creates the Settings with a concrete return type for convenience.
func (b *SettingsBuilder) BuildSettings() *entities.Settings {
	return &entities.Settings{
		TerraCloud:               b.terraCloud,
		TerraTerraformWorkspace:  b.terraTerraformWorkspace,
		TerraAwsRoleArn:          b.terraAwsRoleArn,
		TerraAzureSubscriptionID: b.terraAzureSubscriptionID,
		TerraModuleCacheDir:      b.terraModuleCacheDir,
		TerraProviderCacheDir:    b.terraProviderCacheDir,
		TerraNoCAS:               b.terraNoCAS,
		TerraNoProviderCache:     b.terraNoProviderCache,
		TerraNoPartialParseCache: b.terraNoPartialParseCache,
		TerraNoWorkspace:         b.terraNoWorkspace,
	}
}

// Reset clears the builder state, allowing it to be reused.
func (b *SettingsBuilder) Reset() testkit.Builder {
	b.BaseBuilder.Reset()
	b.terraCloud = ""
	b.terraTerraformWorkspace = ""
	b.terraAwsRoleArn = ""
	b.terraAzureSubscriptionID = ""
	b.terraModuleCacheDir = ""
	b.terraProviderCacheDir = ""
	b.terraNoCAS = false
	b.terraNoProviderCache = false
	b.terraNoPartialParseCache = false
	b.terraNoWorkspace = false
	return b
}

// Clone creates a deep copy of the SettingsBuilder.
func (b *SettingsBuilder) Clone() testkit.Builder {
	return &SettingsBuilder{
		BaseBuilder:              b.BaseBuilder.Clone().(*testkit.BaseBuilder),
		terraCloud:               b.terraCloud,
		terraTerraformWorkspace:  b.terraTerraformWorkspace,
		terraAwsRoleArn:          b.terraAwsRoleArn,
		terraAzureSubscriptionID: b.terraAzureSubscriptionID,
		terraModuleCacheDir:      b.terraModuleCacheDir,
		terraProviderCacheDir:    b.terraProviderCacheDir,
		terraNoCAS:               b.terraNoCAS,
		terraNoProviderCache:     b.terraNoProviderCache,
		terraNoPartialParseCache: b.terraNoPartialParseCache,
		terraNoWorkspace:         b.terraNoWorkspace,
	}
}
