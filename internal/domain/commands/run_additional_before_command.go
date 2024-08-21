package commands

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"slices"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

type RunAdditionalBeforeCommand struct {
	settings   *entities.Settings
	repository repositories.ShellRepository
}

func NewRunAdditionalBeforeCommand(
	settings *entities.Settings,
	repository repositories.ShellRepository,
) *RunAdditionalBeforeCommand {
	return &RunAdditionalBeforeCommand{
		settings:   settings,
		repository: repository,
	}
}

func (it *RunAdditionalBeforeCommand) Execute(targetPath string, arguments []string, listeners interfaces.RunAdditionalBeforeListeners) {
	var cli entities.CLI
	entities.RetrieveCLI(&cli, it.settings)

	// change account if necessary
	if cli != nil && cli.CanChangeAccount() {
		err := it.repository.ExecuteCommand(cli.GetName(), cli.GetCommandChangeAccount(), targetPath)
		if err != nil {
			logger.Errorf("Error changing account: %s", err)
			listeners.OnError(err)
			return
		}
	}

	// init environment if necessary
	if shouldInitEnvironment(arguments) {
		err := it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetPath)
		if err != nil {
			logger.Errorf("Error initializing the environment: %s", err)
			listeners.OnError(err)
			return
		}
	}

	// change workspace if necessary
	if value, ok := it.shouldChangeWorkspace(); ok {
		err := it.repository.ExecuteCommand("terragrunt", []string{"workspace", "select", "-or-create", value}, targetPath)
		if err != nil {
			logger.Errorf("Error changing workspace: %s", err)
			listeners.OnError(err)
			return
		}
	}

	listeners.OnSuccess()
}

func (it *RunAdditionalBeforeCommand) shouldChangeWorkspace() (string, bool) {
	workspace := it.settings.TerraTerraformWorkspace
	return workspace, workspace != ""
}

func shouldInitEnvironment(arguments []string) bool {
	undesiredCommands := []string{"init", "run-all"}
	return !slices.Contains(undesiredCommands, arguments[0])
}
