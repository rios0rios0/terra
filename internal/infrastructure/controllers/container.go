package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"go.uber.org/dig"
)

// RegisterProviders registers all controller providers with the DIG container.
func RegisterProviders(container *dig.Container) error {
	// Register dependencies value
	if err := container.Provide(func() []entities.Dependency {
		return []entities.Dependency{
			{
				Name:              "Terraform",
				CLI:               "terraform",
				BinaryURL:         "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_%[2]s_%[3]s.zip",
				VersionURL:        "https://checkpoint-api.hashicorp.com/v1/check/terraform",
				RegexVersion:      `"current_version":"([^"]+)"`,
				FormattingCommand: []string{"fmt", "-recursive"},
			},
			{
				Name:              "Terragrunt",
				CLI:               "terragrunt",
				BinaryURL:         "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_%[2]s_%[3]s",
				VersionURL:        "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
				RegexVersion:      `"tag_name":"v([^"]+)"`,
				FormattingCommand: []string{"hcl", "format", "**/*.hcl"},
			},
		}
	}); err != nil {
		return err
	}

	// Register controller constructors
	if err := container.Provide(NewDeleteCacheController); err != nil {
		return err
	}
	if err := container.Provide(NewFormatFilesController); err != nil {
		return err
	}
	if err := container.Provide(NewInstallDependenciesController); err != nil {
		return err
	}
	if err := container.Provide(NewSelfUpdateController); err != nil {
		return err
	}
	if err := container.Provide(NewVersionController); err != nil {
		return err
	}
	if err := container.Provide(NewControllers); err != nil {
		return err
	}
	if err := container.Provide(NewRunFromRootController); err != nil {
		return err
	}

	return nil
}

// NewControllers could be duplicated depending on the structure of the application.
func NewControllers(
	deleteCacheController *DeleteCacheController,
	formatFilesController *FormatFilesController,
	installDependenciesController *InstallDependenciesController,
	selfUpdateController *SelfUpdateController,
	versionController *VersionController,
) *[]entities.Controller {
	return &[]entities.Controller{
		deleteCacheController,
		formatFilesController,
		installDependenciesController,
		selfUpdateController,
		versionController,
	}
}
