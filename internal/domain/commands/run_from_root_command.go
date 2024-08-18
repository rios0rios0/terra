package commands

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"slices"
)

type RunFromRootCommand struct {
	formatCommand FormatFilesCommand
	repository    repositories.ShellRepository
}

func (it RunFromRootCommand) Execute(dependencies []entities.Dependency) {
	ensureDependenciesInstallation()

	terraArgs, absDir := findAbsDirectory(args)

	it.formatCommand.Execute(dependencies)
	changeSubscription(absDir)

	undesiredCommands := []string{"init", "run-all"}
	if !slices.Contains(undesiredCommands, terraArgs[0]) {
		_ = it.repository.ExecuteCommand("terragrunt", []string{"init"}, absDir)

		changeWorkspace(absDir)
	}

	err := it.repository.ExecuteCommand("terragrunt", terraArgs, absDir)
	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
}

func findAbsDirectory(args []string) ([]string, string) {
	dir := "."
	terraArgs := args

	// check if the first or last argument is a directory
	if _, err := os.Stat(args[0]); err == nil {
		dir = args[0]
		terraArgs = args[1:]
	} else if _, err := os.Stat(args[len(args)-1]); err == nil {
		dir = args[len(args)-1]
		terraArgs = args[:len(args)-1]
	}

	// convert to the absolute path TODO: it might not be necessary
	absDir, err := filepath.Abs(dir)
	if err != nil {
		logger.Fatalf("Error resolving directory path: %s", err)
	}
	return terraArgs, absDir
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
