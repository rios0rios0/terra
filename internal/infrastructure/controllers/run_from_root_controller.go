package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type RunFromRootController struct {
	command      commands.RunFromRoot
	dependencies []entities.Dependency
}

func NewRunFromRootController(command commands.RunFromRoot) *RunFromRootController {
	return &RunFromRootController{command: command}
}

func (it RunFromRootController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "terra [flags] [terragrunt command] [directory]",
		Short: "Terra is a CLI wrapper for Terragrunt",
		Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. " +
			"It also allows changing the account/subscription and workspace for AWS and Azure.",
	}
}

func (it RunFromRootController) Execute(_ *cobra.Command, arguments []string) {
	absolutePath := findAbsolutePath(arguments)
	it.command.Execute(absolutePath, it.dependencies)
}

//func findAbsDirectory(arguments []string) ([]string, string) {
//	relativePath := "."
//	filteredArgs := arguments
//
//	// check if the first or last argument is a directory
//	if _, err := os.Stat(arguments[0]); err == nil {
//		relativePath = arguments[0]
//		filteredArgs = arguments[1:]
//	} else if _, err := os.Stat(arguments[len(arguments)-1]); err == nil {
//		relativePath = arguments[len(arguments)-1]
//		filteredArgs = arguments[:len(arguments)-1]
//	}
//
//	// convert to the absolute path TODO: it might not be necessary
//	absolutePath, err := filepath.Abs(relativePath)
//	if err != nil {
//		logger.Fatalf("Error resolving directory path: %s", err)
//	}
//	return filteredArgs, absolutePath
//}

func findRelativePath(arguments []string) (string, int) {
	position := -1
	relativePath := "."

	// check if the first or last argument is a directory
	if _, err := os.Stat(arguments[0]); err == nil {
		position = 0
		relativePath = arguments[position]
	} else if _, err := os.Stat(arguments[len(arguments)-1]); err == nil {
		position = len(arguments) - 1
		relativePath = arguments[position]
	}
	return relativePath, position
}

func filterArguments(arguments []string) []string {
	_, position := findRelativePath(arguments)
	return append(arguments[:position], arguments[position+1:]...)
}

func findAbsolutePath(arguments []string) string {
	relativePath, _ := findRelativePath(arguments)
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		logger.Fatalf("Error resolving directory path: %s", err)
	}
	return absolutePath
}
