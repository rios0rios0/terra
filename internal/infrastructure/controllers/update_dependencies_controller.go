package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type UpdateDependenciesController struct {
	command      commands.InstallDependencies
	dependencies []entities.Dependency
}

func NewUpdateDependenciesController(
	command commands.InstallDependencies,
	dependencies []entities.Dependency,
) *UpdateDependenciesController {
	return &UpdateDependenciesController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it *UpdateDependenciesController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "update",
		Short: "Install or update Terraform and Terragrunt to the latest versions",
		Long:  "Install all the dependencies required to run Terra, or update them if newer versions are available. Dependencies are installed to ~/.local/bin on Linux. This is an alias for the 'install' command.",
	}
}

func (it *UpdateDependenciesController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute(it.dependencies)
}