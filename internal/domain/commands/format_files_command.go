package commands

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

type FormatFilesCommand struct {
	repository repositories.ShellRepository
}

func NewFormatFilesCommand(repository repositories.ShellRepository) *FormatFilesCommand {
	return &FormatFilesCommand{repository: repository}
}

func (it *FormatFilesCommand) Execute(dependencies []entities.Dependency, listeners interfaces.FormatFilesListeners) {
	logger.Info("Formatting the code...")
	for _, dependency := range dependencies {
		err := it.repository.ExecuteCommand(dependency.CLI, dependency.FormattingCommand, ".")
		if err != nil {
			logger.Warnf("Failed to format '%s' files: %s", dependency.CLI, err)
			listeners.OnError(err)
			return
		}
	}

	listeners.OnSuccess()
}
