package main

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var dirsToDelete = []string{".terraform", ".terragrunt-cache"}

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
					logger.Infof("Removing directory: %s", path)
					return os.RemoveAll(path)
				}
				return nil
			})
		}
	},
}
