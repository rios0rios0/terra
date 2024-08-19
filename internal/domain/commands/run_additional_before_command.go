package commands

import (
	"os"
	"slices"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

type RunAdditionalBeforeCommand struct {
	cloudCLI   entities.CloudCLI
	repository repositories.ShellRepository
}

func NewRunAdditionalBeforeCommand(
	cloudCLI entities.CloudCLI,
	repository repositories.ShellRepository,
) *RunAdditionalBeforeCommand {
	return &RunAdditionalBeforeCommand{
		cloudCLI:   cloudCLI,
		repository: repository,
	}
}

func (it RunAdditionalBeforeCommand) Execute(targetPath string, arguments []string) {
	// change account if necessary
	if it.cloudCLI != nil && it.cloudCLI.CanChangeAccount() {
		err := it.repository.ExecuteCommand(it.cloudCLI.GetCLIName(), it.cloudCLI.GetCommandChangeAccount(), targetPath)
		if err != nil {
			logger.Fatalf("Error changing account: %s", err)
		}
	}

	// init environment if necessary
	if shouldInitEnvironment(arguments) {
		_ = it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetPath)
	}

	// change workspace if necessary
	if value, ok := shouldChangeWorkspace(); ok {
		err := it.repository.ExecuteCommand("terragrunt", []string{"workspace", "select", "-or-create", value}, targetPath)
		if err != nil {
			logger.Fatalf("Error changing workspace: %s", err)
		}
	}
}

func shouldInitEnvironment(arguments []string) bool {
	undesiredCommands := []string{"init", "run-all"}
	return !slices.Contains(undesiredCommands, arguments[0])
}

func shouldChangeWorkspace() (string, bool) {
	workspace := os.Getenv("TERRA_WORKSPACE")
	return workspace, workspace != ""
}
