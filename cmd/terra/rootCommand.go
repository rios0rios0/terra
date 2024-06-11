package main

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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

		//if !isInitNeeded(absDir) {
		_ = runInDir("terragrunt", []string{"init"}, absDir)
		//}
		// TODO: this is not working with tfvars file
		changeWorkspace(absDir)

		err = runInDir("terragrunt", terraArgs, absDir)
		if err != nil {
			logger.Fatalf("Terragrunt command failed: %s", err)
		}
	},
}

type Terraform struct {
	Source string `hcl:"source,optional"`
}

type Terragrunt struct {
	Terraform Terraform `hcl:"terraform,block"`
}

func isInitNeeded(absDir string) bool {
	// Construct the path to the terragrunt.hcl file
	hclPath := filepath.Join(absDir, "terragrunt.hcl")

	// Parse the terragrunt.hcl file
	var tgConfig Terragrunt
	hclsimple.DecodeFile(hclPath, nil, &tgConfig)

	// Extract the "needToBeFound" path from the "source" field
	needToBeFound := strings.TrimPrefix(tgConfig.Terraform.Source, "../..//")

	// Construct the path to the ".terraform/providers" directory
	providerDir := filepath.Join(absDir, ".terragrunt-cache", "*", "*", needToBeFound, ".terraform", "providers")

	// Check if the ".terraform/providers" directory exists
	_, err := os.Stat(providerDir)
	if os.IsNotExist(err) {
		// The ".terraform/providers" directory does not exist, so "init" is needed
		return true
	}

	// The ".terraform/providers" directory exists, so "init" is not needed
	return false
}

func changeWorkspace(dir string) {
	workspace := ""
	acceptedEnvs := []string{"TF_VAR_env", "TF_VAR_environment"}
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
