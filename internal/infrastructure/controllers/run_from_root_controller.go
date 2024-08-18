package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type RunFromRootController struct {
	command      commands.RunFromRoot
	dependencies []entities.Dependency
}

func NewRunFromRootController(command commands.RunFromRoot) *RunFromRootController {
	return &RunFromRootController{command: command}
}

func (it RunFromRootController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "terra [flags] [terragrunt command] [directory]",
		Short: "Terra is a CLI wrapper for Terragrunt",
		Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. " +
			"It also allows changing the account/subscription and workspace for AWS and Azure.",
	}
}

func (it RunFromRootController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute(it.dependencies)
}
