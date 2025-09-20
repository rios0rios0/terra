package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type VersionController struct {
	command commands.Version
}

func NewVersionController(command commands.Version) *VersionController {
	return &VersionController{command: command}
}

func (it *VersionController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "version",
		Short: "Show Terra, Terraform, and Terragrunt versions",
		Long:  "Display the version information for Terra and its dependencies (Terraform and Terragrunt).",
	}
}

func (it *VersionController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute()
}
