package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dirsFound    []string
	dirsToDelete = []string{".terraform", ".terragrunt-cache"}
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cache and modules directories",
	Long:  fmt.Sprintf("Clear all the [%s] directories.", strings.Join(dirsToDelete, ", ")),
	Run: func(cmd *cobra.Command, args []string) {
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
