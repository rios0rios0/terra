package main

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format all files in the current directory",
	Long:  "Format all the Terraform and Terragrunt files in the current directory.",
	Run: func(_ *cobra.Command, _ []string) {
		format()
	},
}

func format() {
	logger.Info("Formatting the code...")
	err := runInDir("terragrunt", []string{"hclfmt", "**/*.hcl"}, ".")
	if err != nil {
		logger.Warnf("Failed to format Terragrunt files: %s", err)
	}
	err = runInDir("terraform", []string{"fmt", "-recursive"}, ".")
	if err != nil {
		logger.Warnf("Failed to format Terraform files: %s", err)
	}
}
