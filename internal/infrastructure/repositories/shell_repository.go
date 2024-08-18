package repositories

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type ShellRepository struct {
}

func NewShellRepository() *ShellRepository {
	return &ShellRepository{}
}

func (it ShellRepository) ExecuteCommand(command string, arguments []string, directory string) error {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(arguments, " "), directory)
	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
