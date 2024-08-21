package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type InstallDependenciesController struct {
	command      interfaces.InstallDependencies
	dependencies []entities.Dependency
}

func NewInstallDependenciesController(
	command interfaces.InstallDependencies,
	dependencies []entities.Dependency,
) *InstallDependenciesController {
	return &InstallDependenciesController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it *InstallDependenciesController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "install",
		Short: "Install Terraform and Terragrunt (they are pre-requisites)",
		Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	}
}

func (it *InstallDependenciesController) Execute(_ *cobra.Command, _ []string) {
	listeners := interfaces.InstallDependenciesListeners{
		OnSuccess: func() {
			logger.Info("Dependencies installed successfully")
		},
		OnError: func(err error) {
			logger.Errorf("Failed to install dependencies: %v", err)
		},
	}
	it.command.Execute(it.dependencies, listeners)
}
