package commands

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type DeleteCacheCommand struct {
}

func NewDeleteCacheCommand() *DeleteCacheCommand {
	return &DeleteCacheCommand{}
}

func (it DeleteCacheCommand) Execute(toBeDeleted []string) {
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
			}
		}
	}
}
