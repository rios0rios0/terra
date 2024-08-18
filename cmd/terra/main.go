package main

import (
	"github.com/joho/godotenv"
	"github.com/rios0rios0/terra/cmd/terra/domain/commands"
	"github.com/rios0rios0/terra/cmd/terra/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var dependencies = []entities.Dependency{
	{
		Name:              "Terraform",
		CLI:               "terraform",
		BinaryURL:         "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip",
		VersionURL:        "https://checkpoint-api.hashicorp.com/v1/check/terraform",
		RegexVersion:      `"current_version":"([^"]+)"`,
		FormattingCommand: []string{"fmt", "-recursive"},
	},
	{
		Name:              "Terragrunt",
		CLI:               "terragrunt",
		BinaryURL:         "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64",
		VersionURL:        "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest",
		RegexVersion:      `"tag_name":"v([^"]+)"`,
		FormattingCommand: []string{"hclfmt", "**/*.hcl"},
	},
}

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install Terraform and Terragrunt (they are pre-requisites)",
	Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	Run: func(_ *cobra.Command, _ []string) {
		commands.InstallDependenciesCommand{}.Execute(dependencies)
	},
}

var formatCommand = &cobra.Command{
	Use:   "format",
	Short: "Format all files in the current directory",
	Long:  "Format all the Terraform and Terragrunt files in the current directory.",
	Run: func(_ *cobra.Command, _ []string) {
		commands.FormatFilesCommand{}.Execute(dependencies)
	},
}

var clearCommand = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cache and modules directories",
	Long:  "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.",
	Run: func(_ *cobra.Command, _ []string) {
		command := commands.ClearCacheCommand{}
		command.Execute([]string{".terraform", ".terragrunt-cache"})
	},
}

var rootCmd = &cobra.Command{
	Use:   "terra [flags] [terragrunt command] [directory]",
	Short: "Terra is a CLI wrapper for Terragrunt",
	Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. " +
		"It also allows changing the account/subscription and workspace for AWS and Azure.",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	Run: func(_ *cobra.Command, args []string) {

	},
}

func changeWorkspace(dir string) {
	// TODO: this is not working with "tfvars" file
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

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Warnf("Error loading .env file: %s", err)
	}

	rootCmd.AddCommand(installCommand)
	rootCmd.AddCommand(formatCommand)
	rootCmd.AddCommand(clearCommand)
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Error executing 'terra': %s", err)
	}
}
