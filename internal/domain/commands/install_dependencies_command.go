package commands

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
	"os"
	"path"
)

type InstallDependenciesCommand struct {
	osRepository  repositories.OSRepository
	webRepository repositories.WebStringsRepository
}

func NewInstallDependenciesCommand() *InstallDependenciesCommand {
	return &InstallDependenciesCommand{}
}

func (it *InstallDependenciesCommand) Execute(dependencies []entities.Dependency, listeners interfaces.InstallDependenciesListeners) {
	if os.Geteuid() != 0 {
		listeners.OnError(fmt.Errorf("run this command with root privileges to install the dependencies"))
		return
	}

	for _, dependency := range dependencies {
		latestVersion, err := it.webRepository.FindStringMatchInURL(dependency.VersionURL, dependency.RegexVersion)
		if err != nil {
			listeners.OnError(fmt.Errorf("failed to fetch latest version for %s: %w", dependency.Name, err))
			return
		}

		if !dependency.IsAvailable() {
			logger.Infof("'%s' is not installed, installing now...", dependency.Name)

			downloadURL := dependency.GetDownloadURL(latestVersion)
			logger.Infof("Downloading '%s' from '%s'...", dependency.Name, downloadURL)

			currentOS := entities.GetOS()
			temporaryPath := path.Join(currentOS.GetTempDir(), dependency.Name)
			destinationPath := path.Join(currentOS.GetInstallationPath(), dependency.Name)
			if err = currentOS.Download(downloadURL, temporaryPath); err != nil {
				listeners.OnError(fmt.Errorf("failed to download %s: %s", dependency.Name, err))
				return
			}

			if err = it.osRepository.InstallExecutable(temporaryPath, destinationPath, currentOS); err != nil {
				listeners.OnError(fmt.Errorf("failed to install dependency %s: %w", dependency.Name, err))
				return
			}
		}
	}

	listeners.OnSuccess()
}
