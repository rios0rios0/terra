package repositories

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

// Allow long-running terraform/terragrunt commands

// StdShellRepository is not totally necessary, but it is rather a good example for other applications.
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
	start := time.Now()
	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	elapsed := time.Since(start)
	if err != nil {
		logger.Warnf(
			"Failed [%s %s] in %s (took %.2fs)",
			command,
			strings.Join(arguments, " "),
			directory,
			elapsed.Seconds(),
		)
		err = fmt.Errorf("failed to perform command execution: %w", err)
	} else {
		logger.Infof(
			"Completed [%s %s] in %s (took %.2fs)",
			command,
			strings.Join(arguments, " "),
			directory,
			elapsed.Seconds(),
		)
	}
	return err
}
