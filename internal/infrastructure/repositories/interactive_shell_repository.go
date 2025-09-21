package repositories

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
)

const (
	outputChannelSize   = 100
	outputCheckInterval = 100 * time.Millisecond
)

// InteractiveShellRepository handles interactive commands with auto-answering capabilities
type InteractiveShellRepository struct{}

func NewInteractiveShellRepository() *InteractiveShellRepository {
	return &InteractiveShellRepository{}
}

func (it *InteractiveShellRepository) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	logger.Infof(
		"Running [%s %s] in %s with auto-answering",
		command,
		strings.Join(arguments, " "),
		directory,
	)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory

	// Set up pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if startErr := cmd.Start(); startErr != nil {
		return fmt.Errorf("failed to start command: %w", startErr)
	}

	// Handle output and input in separate goroutines
	go it.handleOutput(stdout, stderr, stdin)

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}

	return err
}

func (it *InteractiveShellRepository) handleOutput(
	stdout, stderr io.ReadCloser,
	stdin io.WriteCloser,
) {
	// Create scanners for both stdout and stderr
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	// Channel to coordinate between output handling and input sending
	outputChan := make(chan string, outputChannelSize)

	// Read stdout
	go func() {
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			logger.Info(line) // Print to user
			outputChan <- line
		}
	}()

	// Read stderr
	go func() {
		for stderrScanner.Scan() {
			line := stderrScanner.Text()
			fmt.Fprintln(os.Stderr, line) // Print to stderr
			outputChan <- line
		}
	}()

	// Process output and send responses
	go func() {
		for {
			select {
			case line := <-outputChan:
				it.processLineAndRespond(line, stdin)
			case <-time.After(outputCheckInterval):
				// Continue checking for output
			}
		}
	}()
}

func (it *InteractiveShellRepository) processLineAndRespond(line string, stdin io.WriteCloser) {
	// Remove ANSI escape codes for pattern matching
	cleanLine := it.removeANSICodes(line)

	// Pattern 1: External dependency prompt - answer "n"
	externalDepPattern := regexp.MustCompile(
		`(?i)should terragrunt apply the external dependency.*\?`,
	)
	if externalDepPattern.MatchString(cleanLine) {
		logger.Debug("Detected external dependency prompt, responding with 'n'")
		fmt.Fprintln(stdin, "n")
		return
	}

	// Pattern 2: "Are you sure you want to run" prompt - drop to manual mode
	confirmationPattern := regexp.MustCompile(`(?i)are you sure you want to run.*`)
	if confirmationPattern.MatchString(cleanLine) {
		logger.Info("Detected confirmation prompt, switching to manual mode")
		// For confirmation prompts, we let the user handle it manually
		// by copying stdin from the terminal
		go it.copyStdinToProcess(stdin)
		return
	}

	// Pattern 3: Any other "yes/no" prompts - answer "n" by default
	yesNoPattern := regexp.MustCompile(`(?i).*\?.*\[y/n\]`)
	if yesNoPattern.MatchString(cleanLine) {
		logger.Debug("Detected yes/no prompt, responding with 'n'")
		fmt.Fprintln(stdin, "n")
		return
	}
}

func (it *InteractiveShellRepository) copyStdinToProcess(stdin io.WriteCloser) {
	// Copy user input from terminal to the process
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(stdin, line)
	}
}

func (it *InteractiveShellRepository) removeANSICodes(text string) string {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(text, "")
}
