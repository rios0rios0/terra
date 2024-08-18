package main

import (
	"github.com/rios0rios0/terra/cmd/terra/domain/commands"
	"github.com/rios0rios0/terra/cmd/terra/domain/entities"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Terraform and Terragrunt (they are pre-requisites)",
	Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	Run: func(_ *cobra.Command, _ []string) {
		dependencies := []entities.Dependency{
			{
				Name:         "Terraform",
				CLI:          "terraform",
				BinaryURL:    "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip",
				VersionURL:   "https://checkpoint-api.hashicorp.com/v1/check/terraform",
				RegexVersion: `"current_version":"([^"]+)"`,
			},
			{
				Name:         "Terragrunt",
				CLI:          "terragrunt",
				BinaryURL:    "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64",
				VersionURL:   "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
				RegexVersion: `"tag_name":"v([^"]+)"`,
			},
		}

		command := commands.InstallDependenciesCommand{}
		command.Execute(dependencies)
	},
}
