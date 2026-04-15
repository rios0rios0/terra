package controllers

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers/helpers"
	"github.com/spf13/cobra"
)

type RunFromRootController struct {
	command      commands.RunFromRoot
	dependencies []entities.Dependency
}

func NewRunFromRootController(
	command commands.RunFromRoot,
	dependencies []entities.Dependency,
) *RunFromRootController {
	return &RunFromRootController{
		command:      command,
		dependencies: dependencies,
	}
}

func (it *RunFromRootController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{
		Use:   "terra [flags] [terragrunt command] [directory]",
		Short: "Terra is a CLI wrapper for Terragrunt",
		Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before " +
			"executing commands. It also switches the account/subscription and workspace for " +
			"AWS and Azure automatically based on the .env configuration.\n" +
			"\n" +
			"Parallel execution strategies:\n" +
			"\n" +
			"  --parallel=N   Terra-managed worker pool. Supports --only=mod1,mod2 and\n" +
			"                 --skip=mod1,mod2 for basename-matched module selection.\n" +
			"                 Required for state operations (import, state rm, state mv).\n" +
			"\n" +
			"  --all          Terragrunt-managed run-all. Uses terragrunt's DAG. Filter\n" +
			"                 with --filter='!mod' (recommended; supports globs, graph,\n" +
			"                 and git-diff expressions) or --queue-exclude-dir=mod.\n" +
			"\n" +
			"These two strategies cannot be combined. See docs/parallel-execution.md.",
	}
}

func (it *RunFromRootController) Execute(_ *cobra.Command, arguments []string) {
	absolutePath := helpers.ArgumentsHelper{}.FindAbsolutePath(arguments)
	filteredArguments := helpers.ArgumentsHelper{}.RemovePathFromArguments(arguments)
	it.command.Execute(absolutePath, filteredArguments, it.dependencies)
}
