package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DeleteCacheController struct {
	command interfaces.DeleteCache
}

func NewDeleteCacheController(command interfaces.DeleteCache) *DeleteCacheController {
	return &DeleteCacheController{command: command}
}

func (it *DeleteCacheController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "clear",
		Short: "Clear all cache and modules directories",
		Long:  "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.",
	}
}

func (it *DeleteCacheController) Execute(_ *cobra.Command, _ []string) {
	listeners := interfaces.DeleteCacheListeners{
		OnSuccess: func() {
			logger.Info("Cache cleared successfully")
		},
		OnError: func(err error) {
			logger.Errorf("Failed to clear cache: %v", err)
		},
	}
	it.command.Execute([]string{".terraform", ".terragrunt-cache"}, listeners)
}
