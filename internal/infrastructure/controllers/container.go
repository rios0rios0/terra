package controllers

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/entities"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewDependencies,
	// the root controller is not defined here, it is defined in the "wire.go" file
	NewDeleteCacheController,
	NewFormatFilesController,
	NewInstallDependenciesController,
	NewControllers,
)

// NewControllers could be duplicated depending on the structure of the application.
func NewControllers(
	deleteCacheController *DeleteCacheController,
	formatFilesController *FormatFilesController,
	installDependenciesController *InstallDependenciesController,
) *[]entities.Controller {
	return &[]entities.Controller{
		deleteCacheController,
		formatFilesController,
		installDependenciesController,
	}
}

// NewDependencies creates dependencies with URLs that can be overridden via environment variables
func NewDependencies(settings *entities.Settings) []entities.Dependency {
	terraformVersionURL := "https://checkpoint-api.hashicorp.com/v1/check/terraform"
	terraformBinaryURL := "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip"
	terragruntVersionURL := "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest"
	terragruntBinaryURL := "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64"

	// Override with environment variables if provided
	if settings.TerraformVersionURL != "" {
		terraformVersionURL = settings.TerraformVersionURL
	}
	if settings.TerraformBinaryURL != "" {
		terraformBinaryURL = settings.TerraformBinaryURL
	}
	if settings.TerragruntVersionURL != "" {
		terragruntVersionURL = settings.TerragruntVersionURL
	}
	if settings.TerragruntBinaryURL != "" {
		terragruntBinaryURL = settings.TerragruntBinaryURL
	}

	return []entities.Dependency{
		{
			Name:              "Terraform",
			CLI:               "terraform",
			BinaryURL:         terraformBinaryURL,
			VersionURL:        terraformVersionURL,
			RegexVersion:      `"current_version":"([^"]+)"`,
			FormattingCommand: []string{"fmt", "-recursive"},
		},
		{
			Name:              "Terragrunt",
			CLI:               "terragrunt",
			BinaryURL:         terragruntBinaryURL,
			VersionURL:        terragruntVersionURL,
			RegexVersion:      `"tag_name":"v([^"]+)"`,
			FormattingCommand: []string{"hclfmt", "**/*.hcl"},
		},
	}
}

//nolint:gochecknoglobals
var dependencies = wire.Value([]entities.Dependency{
	{
		Name:              "Terraform",
		CLI:               "terraform",
		BinaryURL:         "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip",
		VersionURL:        "https://checkpoint-api.hashicorp.com/v1/check/terraform",
		RegexVersion:      `"current_version":"([^"]+)"`,
		FormattingCommand: []string{"fmt", "-recursive"},
	},
	{
		Name:              "Terragrunt",
		CLI:               "terragrunt",
		BinaryURL:         "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64",
		VersionURL:        "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
		RegexVersion:      `"tag_name":"v([^"]+)"`,
		FormattingCommand: []string{"hclfmt", "**/*.hcl"},
	},
})
