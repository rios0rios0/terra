package main

import (
	"os"

	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	logger "github.com/sirupsen/logrus"
)

//nolint:unparam
func injectApp() entities.App {
	wire.Build(
		internal.Container,
		internal.NewAppCLI,
		newCloudCLI,
		newDependencies,
	)
	return nil
}

//nolint:unparam
func injectRootController() entities.Controller {
	wire.Build(
		internal.Container,
		controllers.NewRunFromRootController,
		newCloudCLI,
		newDependencies,
	)
	return nil
}

func newCloudCLI() entities.CloudCLI {
	mapping := map[string]entities.CloudCLI{
		"aws": entities.CLIAws{},
		"az":  entities.CLIAzm{},
	}

	value, ok := mapping[os.Getenv("TERRA_CLOUD")]
	if !ok {
		value = nil
		logger.Warnf("No cloud CLI found, avoiding to execute customized commands...")
	}
	return value
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
