package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/helpers"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RunFromRootController struct {
	command      interfaces.RunFromRoot
	dependencies []entities.Dependency
}

func NewRunFromRootController(
	command interfaces.RunFromRoot,
	dependencies []entities.Dependency,
) *RunFromRootController {
	return &RunFromRootController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it *RunFromRootController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "terra [flags] [terragrunt command] [directory]",
		Short: "Terra is a CLI wrapper for Terragrunt",
		Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. " +
			"It also allows changing the account/subscription and workspace for AWS and Azure.",
	}
}

func (it *RunFromRootController) Execute(_ *cobra.Command, arguments []string) {
	absolutePath := helpers.ArgumentsHelper{}.FindAbsolutePath(arguments)
	filteredArguments := helpers.ArgumentsHelper{}.RemovePathFromArguments(arguments)
	listeners := interfaces.RunFromRootListeners{
		OnSuccess: func() {
			logger.Info("Command executed successfully")
		},
		OnError: func(err error) {
			logger.Errorf("Failed to execute command: %v", err)
		},
	}
	it.command.Execute(absolutePath, filteredArguments, it.dependencies, listeners)
}
