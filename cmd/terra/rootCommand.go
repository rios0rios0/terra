package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// Main command logic
var rootCmd = &cobra.Command{
	Use:   "terra [flags] [terragrunt command] [directory]",
	Short: "Terra is a CLI wrapper for Terragrunt",
	Long:  "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing the command.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ensureToolsInstalled()
		terraArgs, absDir := findAbsDirectory(args)
		err := runInDir("terragrunt", terraArgs, absDir)
		if err != nil {
			logger.Fatalf("Terragrunt command failed: %s", err)
		}
	},
}

func findAbsDirectory(args []string) ([]string, string) {
	dir := "."
	terraArgs := args

	// Check if the first or last argument is a directory
	if _, err := os.Stat(args[0]); err == nil {
		dir = args[0]
		terraArgs = args[1:]
	} else if _, err := os.Stat(args[len(args)-1]); err == nil {
		dir = args[len(args)-1]
		terraArgs = args[:len(args)-1]
	}

	// Convert to the absolute path TODO: it might not be necessary
	absDir, err := filepath.Abs(dir)
	if err != nil {
		logger.Fatalf("Error resolving directory path: %s", err)
	}
	return terraArgs, absDir
}
