package commands

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

type RunAdditionalBeforeCommand struct {
	settings   *entities.Settings
	cli        entities.CLI
	repository repositories.ShellRepository
}

func NewRunAdditionalBeforeCommand(
	settings *entities.Settings,
	cli entities.CLI,
	repository repositories.ShellRepository,
) *RunAdditionalBeforeCommand {
	return &RunAdditionalBeforeCommand{
		settings:   settings,
		cli:        cli,
		repository: repository,
	}
}

func (it *RunAdditionalBeforeCommand) Execute(targetPath string, arguments []string) {
	// change account if necessary
	if it.cli != nil && it.cli.CanChangeAccount() {
		err := it.repository.ExecuteCommand(
			it.cli.GetName(),
			it.cli.GetCommandChangeAccount(),
			targetPath,
		)
		if err != nil {
			logger.Fatalf("Error changing account: %s", err)
		}
	}

	// init environment if necessary
	if shouldInitEnvironment(arguments, targetPath) {
		_ = it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetPath)
	}

	// change workspace if necessary
	if value, ok := it.shouldChangeWorkspace(); ok {
		err := it.repository.ExecuteCommand(
			"terragrunt",
			[]string{"workspace", "select", "-or-create", value},
			targetPath,
		)
		if err != nil {
			logger.Fatalf("Error changing workspace: %s", err)
		}
	}
}

func (it *RunAdditionalBeforeCommand) shouldChangeWorkspace() (string, bool) {
	workspace := it.settings.TerraTerraformWorkspace
	return workspace, workspace != ""
}

func shouldInitEnvironment(arguments []string, targetPath string) bool {
	// Don't init when the first argument is "init"
	if len(arguments) > 0 && arguments[0] == "init" {
		return false
	}

	// Don't init when using --all flag (equivalent to deprecated run-all)
	if slices.Contains(arguments, "--all") {
		return false
	}

	// Don't init for state manipulation commands: terragrunt handles its own
	// initialization for state operations, and an explicit init triggers full
	// dependency resolution (which can fail on unrelated dependency outputs).
	if IsStateManipulationCommand(arguments) {
		return false
	}

	// Skip init if .terraform already exists — the reactive UpgradeAwareShellRepository
	// will handle stale state if needed.
	if _, err := os.Stat(filepath.Join(targetPath, ".terraform")); err == nil {
		return false
	}

	return true
}
