package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type FormatFilesController struct {
	command      interfaces.FormatFiles
	dependencies []entities.Dependency
}

func NewFormatFilesController(
	command interfaces.FormatFiles,
	dependencies []entities.Dependency,
) *FormatFilesController {
	return &FormatFilesController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it *FormatFilesController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "format",
		Short: "Format all files in the current directory",
		Long:  "Format all the Terraform and Terragrunt files in the current directory.",
	}
}

func (it *FormatFilesController) Execute(_ *cobra.Command, _ []string) {
	listeners := interfaces.FormatFilesListeners{
		OnSuccess: func() {
			logger.Info("Files formatted successfully")
		},
		OnError: func(err error) {
			logger.Errorf("Failed to format files: %v", err)
		},
	}
	it.command.Execute(it.dependencies, listeners)
}
