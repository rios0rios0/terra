package commands

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
	"os"
	"slices"
)

type RunFromRootCommand struct {
	installCommand InstallDependenciesCommand
	formatCommand  FormatFilesCommand
	repository     repositories.ShellRepository
}

func (it RunFromRootCommand) Execute(targetDirectory string, dependencies []entities.Dependency) {
	// ensure that all dependencies are installed
	it.installCommand.Execute(dependencies)

	terraArgs, _ := findAbsDirectory(args)

	it.formatCommand.Execute(dependencies)
	changeSubscription(targetDirectory)

	undesiredCommands := []string{"init", "run-all"}
	if !slices.Contains(undesiredCommands, terraArgs[0]) {
		_ = it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetDirectory)

		changeWorkspace(targetDirectory)
	}

	err := it.repository.ExecuteCommand("terragrunt", terraArgs, targetDirectory)
	if err != nil {
		logger.Fatalf("Terragrunt command failed: %s", err)
	}
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
