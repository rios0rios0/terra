package main

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Run Terragrunt command in the specified directory
func runTerragruntInDir(args []string, dir string) {
	cmd := exec.Command("terragrunt", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

// Main command logic
var rootCmd = &cobra.Command{
	Use:   "terra [flags] [directory] [terragrunt command]",
	Short: "Terra is a CLI wrapper for Terragrunt",
	Long:  `Terra is a CLI wrapper for Terragrunt that allows changing directory before executing the command.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ensureToolsInstalled()

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

		runTerragruntInDir(terraArgs, absDir)
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Terraform and Terragrunt (they are pre-requisites)",
	Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	Run: func(cmd *cobra.Command, args []string) {
		ensureToolsInstalled()
	},
}

var formatCmd = &cobra.Command{
	Use: "fmt",
	Run: func(cmd *cobra.Command, args []string) {
		/* nothing */
	},
}

var clearCmd = &cobra.Command{
	Use: "clear",
	Run: func(cmd *cobra.Command, args []string) {
		/* nothing */
	},
}

func main() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(clearCmd)
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Error executing terra: %s", err)
	}
}
