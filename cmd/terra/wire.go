package main

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
)

func injectApp() entities.App {
	wire.Build(
		internal.Container,
		internal.NewAppCLI,
		newDependencies,
	)
	return nil
}

func injectRootController() entities.Controller {
	wire.Build(
		internal.Container,
		controllers.NewRunFromRootController,
		newDependencies,
	)
	return nil
}

func newDependencies() []entities.Dependency {
	return []entities.Dependency{
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
	}
}
