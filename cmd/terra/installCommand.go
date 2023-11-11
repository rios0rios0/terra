package main

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Terraform and Terragrunt (they are pre-requisites)",
	Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	Run: func(cmd *cobra.Command, args []string) {
		ensureToolsInstalled()
	},
}
