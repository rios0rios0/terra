package main

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Main command logic
var rootCmd = &cobra.Command{
	Use:                "terra [flags] [terragrunt command] [directory]",
	Short:              "Terra is a CLI wrapper for Terragrunt",
	Long:               "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing the command.",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		ensureToolsInstalled()
		terraArgs, absDir := findAbsDirectory(args)
		err := godotenv.Load()
		if err != nil {
			logger.Warnf("Error loading .env file: %s", err)
		}
		format()

		changeSubscription(absDir)

		// TODO: it should be a mapping with all undesired commands
		if terraArgs[0] != "init" && terraArgs[0] != "run-all" {
			_ = runInDir("terragrunt", []string{"init"}, absDir)
			// TODO: this is not working with "tfvars" file
			changeWorkspace(absDir)
		}

		err = runInDir("terragrunt", terraArgs, absDir)
		if err != nil {
			logger.Fatalf("Terragrunt command failed: %s", err)
		}
	},
}

func changeWorkspace(dir string) {
	workspace := ""
	acceptedEnvs := []string{"TERRA_WORKSPACE"}
	for _, env := range acceptedEnvs {
		workspace = os.Getenv(env)
		if workspace != "" {
			break
		}
	}

	if workspace != "" {
		err := runInDir("terragrunt", []string{"workspace", "select", "-or-create", workspace}, dir)
		if err != nil {
			logger.Fatalf("Error changing workspace: %s", err)
		}
	}
}

func changeSubscription(dir string) {
	subscriptionId := ""
	acceptedEnvs := []string{"TERRA_AZURE_SUBSCRIPTION_ID"}
	for _, env := range acceptedEnvs {
		subscriptionId = os.Getenv(env)
		if subscriptionId != "" {
			break
		}
	}

	if subscriptionId != "" {
		err := runInDir("az", []string{"account", "set", "--subscription", subscriptionId}, dir)
		if err != nil {
			logger.Fatalf("Error changing subscription: %s", err)
		}
	}
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
