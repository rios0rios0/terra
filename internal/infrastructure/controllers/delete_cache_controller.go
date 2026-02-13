package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DeleteCacheController struct {
	command commands.DeleteCache
}

func NewDeleteCacheController(command commands.DeleteCache) *DeleteCacheController {
	return &DeleteCacheController{command: command}
}

func (it *DeleteCacheController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "clear",
		Short: "Clear all cache and modules directories",
		Long: "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution. " +
			"Use --global to also remove centralized module and provider cache directories.",
	}
}

func (it *DeleteCacheController) Execute(cmd *cobra.Command, _ []string) {
	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		logger.Warnf("Failed to get global flag: %s, defaulting to false", err)
		global = false
	}
	it.command.Execute([]string{".terraform", ".terragrunt-cache"}, global)
}
