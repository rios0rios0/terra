package commands

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

type DeleteCacheCommand struct{}

func NewDeleteCacheCommand() *DeleteCacheCommand {
	return &DeleteCacheCommand{}
}

func (it *DeleteCacheCommand) Execute(toBeDeleted []string, listeners interfaces.DeleteCacheListeners) {
	var foundDirectories []string
	for _, dir := range toBeDeleted {
		logger.Infof("Clearing all %s directories...", dir)
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && strings.HasSuffix(path, dir) {
				foundDirectories = append(foundDirectories, path)
			}
			return nil
		})

		for _, dirPath := range foundDirectories {
			logger.Infof("Removing directory: %s", dirPath)
			err := os.RemoveAll(dirPath)
			if err != nil {
				logger.Errorf("Failed to remove directory: %s, error: %v", dirPath, err)
				listeners.OnError(err)
				return
			}
		}
	}

	listeners.OnSuccess()
}
