package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type FormatFilesController struct {
	command      commands.FormatFiles
	dependencies []entities.Dependency
}

func NewFormatFilesController(command commands.FormatFiles) *FormatFilesController {
	return &FormatFilesController{command: command}
}

func (it FormatFilesController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "format",
		Short: "Format all files in the current directory",
		Long:  "Format all the Terraform and Terragrunt files in the current directory.",
	}
}

func (it FormatFilesController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute(it.dependencies)
}
