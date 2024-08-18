package main

import (
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cache and modules directories",
	Long:  "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.",
	Run: func(_ *cobra.Command, _ []string) {
		var (
			dirsFound    []string
			dirsToDelete = []string{".terraform", ".terragrunt-cache"}
		)

		for _, dir := range dirsToDelete {
			logger.Infof("Clearing all %s directories...", dir)
			_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() && strings.HasSuffix(path, dir) {
					dirsFound = append(dirsFound, path)
				}
				return nil
			})

			for _, dirPath := range dirsFound {
				logger.Infof("Removing directory: %s", dirPath)
				err := os.RemoveAll(dirPath)
				if err != nil {
					logger.Errorf("Failed to remove directory: %s, error: %v", dirPath, err)
				}
			}
		}
	},
}
