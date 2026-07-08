package repositories

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
)

// Allow long-running terraform/terragrunt commands

// StdShellRepository is not totally necessary, but it is rather a good example for other applications.
type StdShellRepository struct {
	// consoleMu serializes writes to the shared console (os.Stdout/os.Stderr) across
	// concurrent prefixed executions so lines from different modules never interleave
	// mid-line. It lives on the struct (not as a package global) because DIG provides a
	// single StdShellRepository instance that every parallel worker shares.
	consoleMu sync.Mutex
}

func NewStdShellRepository() *StdShellRepository {
	return &StdShellRepository{}
}

// ExecuteCommand runs a command with its stdio connected directly to the terminal. Use it
// for sequential, possibly interactive invocations.
func (it *StdShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	return it.run(command, arguments, directory, os.Stdout, os.Stderr, os.Stdin)
}

// ExecuteCommandWithPrefix runs a command while streaming its stdout and stderr through
// per-line prefix writers, so concurrent module executions stay attributable in the
// combined console output. Stdin is left disconnected because parallel workers cannot
// share a single terminal; callers must run non-interactively (e.g. via --yes/--no).
func (it *StdShellRepository) ExecuteCommandWithPrefix(
	command string,
	arguments []string,
	directory string,
	prefix string,
) error {
	stdout := NewLinePrefixWriter(os.Stdout, prefix, &it.consoleMu)
	stderr := NewLinePrefixWriter(os.Stderr, prefix, &it.consoleMu)

	err := it.run(command, arguments, directory, stdout, stderr, nil)

	// Emit any trailing output that did not end with a newline.
	stdout.Flush()
	stderr.Flush()

	return err
}

// run executes the command with the given stdio wiring and logs its duration.
func (it *StdShellRepository) run(
	command string,
	arguments []string,
	directory string,
	stdout, stderr io.Writer,
	stdin io.Reader,
) error {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(arguments, " "), directory)
	start := time.Now()

	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin

	err := cmd.Run()
	logCommandDuration(command, arguments, directory, time.Since(start), err)
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}
	return err
}

// logCommandDuration logs the elapsed time for a command execution, using Warn level for
// failures and Info level for successes.
func logCommandDuration(command string, arguments []string, directory string, elapsed time.Duration, err error) {
	args := strings.Join(arguments, " ")
	if err != nil {
		logger.Warnf("Failed [%s %s] in %s (took %.2fs)", command, args, directory, elapsed.Seconds())
	} else {
		logger.Infof("Completed [%s %s] in %s (took %.2fs)", command, args, directory, elapsed.Seconds())
	}
}
