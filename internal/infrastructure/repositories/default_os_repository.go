package repositories

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func (it *OsRepository) IsSuperUser() bool {
	return os.Geteuid() == 0
}

func (it *OsRepository) ExecuteCommand(command string, arguments []string, directory string) error {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(arguments, " "), directory)
	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}
	return err
}

func (it *OsRepository) InstallExecutable(sourcePath, destinationPath string, currentOS entities.OS) error {
	if isZipFile(sourcePath) {
		if err := currentOS.Extract(sourcePath, destinationPath); err != nil {
			return fmt.Errorf("failed to extract it: %s", err)
		}
		if err := currentOS.Remove(sourcePath); err != nil {
			return fmt.Errorf("failed to remove it: %s", err)
		}
	} else {
		if err := currentOS.Move(sourcePath, destinationPath); err != nil {
			return fmt.Errorf("failed to move it: %s", err)
		}
	}

	if err := currentOS.MakeExecutable(destinationPath); err != nil {
		return fmt.Errorf("failed to make it executable: %s", err)
	}

	return nil
}
