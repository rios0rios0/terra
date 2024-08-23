package repositories

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"os"
	"os/exec"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// DefaultOSRepository is not totally necessary, but it is rather a good example for other applications
type DefaultOSRepository struct{}

func NewDefaultOSRepository() *DefaultOSRepository {
	return &DefaultOSRepository{}
}

func (it *DefaultOSRepository) IsSuperUser() bool {
	return os.Geteuid() == 0
}

func (it *DefaultOSRepository) ExecuteCommand(command string, arguments []string, directory string) error {
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

func (it *DefaultOSRepository) InstallExecutable(sourcePath, destinationPath string, currentOS entities.OS) error {
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

func isZipFile(filePath string) bool {
	fileTypeCmd := exec.Command("file", filePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Errorf("Failed to determine file type of %s: %s", filePath, err)
		return false
	}
	return strings.Contains(string(fileTypeOutput), "Zip archive data")
}
