package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

type DeleteCacheController struct {
	command commands.DeleteCache
}

func NewDeleteCacheController(command commands.DeleteCache) *DeleteCacheController {
	return &DeleteCacheController{command: command}
}

func (it DeleteCacheController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "clear",
		Short: "Clear all cache and modules directories",
		Long:  "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.",
	}
}

func (it DeleteCacheController) Execute(_ *cobra.Command, _ []string) {
	it.command.Execute([]string{".terraform", ".terragrunt-cache"})
}
