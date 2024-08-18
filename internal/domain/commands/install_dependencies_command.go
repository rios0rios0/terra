package commands

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type InstallDependenciesCommand struct {
}

func NewInstallDependenciesCommand() InstallDependenciesCommand {
	return InstallDependenciesCommand{}
}

func (it InstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	for _, dependency := range dependencies {
		latestVersion := fetchLatestVersion(dependency.VersionURL, dependency.RegexVersion)

		if !isDependencyCLIAvailable(dependency.CLI) {
			ensureRootPrivileges()
			logger.Warnf("%s is not installed, installing now...", dependency.Name)
			install(fmt.Sprintf(dependency.BinaryURL, latestVersion), dependency.CLI)
		}
	}
}

// fetch the latest version of software from a URL
func fetchLatestVersion(url, regexPattern string) string {
	resp, err := http.Get(url) //nolint:gosec // no security issue here all URL are HTTPS
	if err != nil {
		logger.Fatalf("Error fetching version info: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("Error reading response body: %s", err)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1]
	}

	// TODO: it should be better
	logger.Fatalf("No version match found")
	return ""
}

// checking if a dependency is available
func isDependencyCLIAvailable(name string) bool {
	cmd := exec.Command(name, "-v")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// check if the "terra" has root privileges to install dependencies
func ensureRootPrivileges() {
	if os.Geteuid() != 0 {
		logger.Fatalf("Run this command with root privileges to install the dependencies")
		return
	}
}

// installing dependencies doesn't matter the operating system
func install(url, name string) {
	currentOS := entities.GetOS()
	tempFilePath := path.Join(currentOS.GetTempDir(), name)
	destPath := path.Join(currentOS.GetInstallationPath(), name)

	logger.Infof("Downloading %s from %s...", name, url)
	if err := currentOS.Download(url, tempFilePath); err != nil {
		logger.Fatalf("Failed to download %s: %s", name, err)
	}

	fileTypeCmd := exec.Command("file", tempFilePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
	}

	if strings.Contains(string(fileTypeOutput), "Zip archive data") {
		logger.Infof("%s is a zip file, extracting...", name)
		if err := currentOS.Extract(tempFilePath, destPath); err != nil {
			logger.Fatalf("Failed to extract %s: %s", name, err)
		}
		if err := currentOS.Remove(tempFilePath); err != nil {
			logger.Fatalf("Failed to remove %s: %s", name, err)
		}
	} else {
		if err := currentOS.Move(tempFilePath, destPath); err != nil {
			logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
		}
	}

	if err := currentOS.MakeExecutable(destPath); err != nil {
		logger.Fatalf("Failed to make %s executable: %s", name, err)
	}
}
