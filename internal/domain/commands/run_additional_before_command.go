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
		if err := it.repository.ExecuteCommand("terragrunt", []string{"init"}, targetPath); err != nil {
			logger.Warnf("Proactive terragrunt init failed in %s: %v", targetPath, err)
		}
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
	if it.settings.TerraNoWorkspace {
		return "", false
	}
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

	// Skip init if the environment was already initialized — the reactive
	// UpgradeAwareShellRepository will handle stale state if needed.
	// Check for .terraform (plain Terraform) and Terragrunt cache directories
	// (.terragrunt-cache and legacy terragrunt-cache), since Terragrunt places
	// .terraform inside .terragrunt-cache/<hash>/<hash>/.
	for _, dir := range []string{".terraform", ".terragrunt-cache", "terragrunt-cache"} {
		if info, err := os.Stat(filepath.Join(targetPath, dir)); err == nil && info.IsDir() {
			return false
		}
	}

	// When centralized caching is active (TG_DOWNLOAD_DIR set by terra),
	// terragrunt doesn't create .terragrunt-cache locally. Skip init if the
	// centralized cache already has content from a previous run.
	if downloadDir := os.Getenv("TG_DOWNLOAD_DIR"); downloadDir != "" {
		entries, err := os.ReadDir(downloadDir)
		if err == nil && len(entries) > 0 {
			return false
		}
	}

	return true
}
