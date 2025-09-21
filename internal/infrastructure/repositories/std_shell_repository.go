package repositories

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	logger "github.com/sirupsen/logrus"
)

// StdShellRepository is not totally necessary, but it is rather a good example for other applications
type StdShellRepository struct{}

func NewStdShellRepository() *StdShellRepository {
	return &StdShellRepository{}
}

func (it *StdShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
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
