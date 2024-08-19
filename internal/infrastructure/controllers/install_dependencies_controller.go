package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type InstallDependenciesController struct {
	command      commands.InstallDependencies
	dependencies []entities.Dependency
}

func NewInstallDependenciesController(
	command commands.InstallDependencies,
	dependencies []entities.Dependency,
) *InstallDependenciesController {
	return &InstallDependenciesController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it InstallDependenciesController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "install",
		Short: "Install Terraform and Terragrunt (they are pre-requisites)",
		Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	}
}

func (it InstallDependenciesController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute(it.dependencies)
}
