package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
)

type DeleteCacheCommand struct {
	settings *entities.Settings
}

func NewDeleteCacheCommand(settings *entities.Settings) *DeleteCacheCommand {
	return &DeleteCacheCommand{settings: settings}
}

func (it *DeleteCacheCommand) Execute(toBeDeleted []string, global bool) {
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

	if global {
		it.clearGlobalCache()
	}
}

// clearGlobalCache removes the centralized module and provider cache directories.
func (it *DeleteCacheCommand) clearGlobalCache() {
	moduleDir, moduleDirErr := it.settings.GetModuleCacheDir()
	if moduleDirErr != nil {
		logger.Errorf("Failed to determine module cache directory: %s", moduleDirErr)
		return
	}

	providerDir, providerDirErr := it.settings.GetProviderCacheDir()
	if providerDirErr != nil {
		logger.Errorf("Failed to determine provider cache directory: %s", providerDirErr)
		return
	}

	for _, dir := range []string{moduleDir, providerDir} {
		if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
			logger.Infof("Global cache directory does not exist, skipping: %s", dir)
			continue
		}
		logger.Infof("Removing global cache directory: %s", dir)
		if removeErr := os.RemoveAll(dir); removeErr != nil {
			logger.Errorf("Failed to remove global cache directory: %s, error: %v", dir, removeErr)
		}
	}
}
