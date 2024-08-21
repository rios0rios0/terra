package commands

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"slices"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
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
			listeners.OnError(fmt.Errorf("error changing account: %w", err))
			return
		}
	}

	// init environment if necessary
	if shouldInitEnvironment(arguments) {
		err := it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetPath)
		if err != nil {
			listeners.OnError(fmt.Errorf("error initializing the environment: %w", err))
			return
		}
	}

	// change workspace if necessary
	if it.settings.TerraTerraformWorkspace != "" {
		err := it.repository.ExecuteCommand("terragrunt", []string{"workspace", "select", "-or-create", it.settings.TerraTerraformWorkspace}, targetPath)
		if err != nil {
			listeners.OnError(fmt.Errorf("error changing workspace: %w", err))
			return
		}
	}

	listeners.OnSuccess()
}

func shouldInitEnvironment(arguments []string) bool {
	undesiredCommands := []string{"init", "run-all"}
	return !slices.Contains(undesiredCommands, arguments[0])
}
