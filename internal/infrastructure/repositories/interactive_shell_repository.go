package repositories

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/creack/pty"
	logger "github.com/sirupsen/logrus"
)

const (
	bufferSize          = 1024
	maxOutputBufferSize = 4096
	outputTrimSize      = 2048
)

// InteractiveShellRepository handles interactive commands with auto-answering capabilities.
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

	cmd := exec.CommandContext(context.Background(), command, arguments...)
	cmd.Dir = directory

	// Set environment to reduce ANSI sequences, similar to expect script
	cmd.Env = append(os.Environ(), "TERM=dumb")

	// Start the command with a pseudo-terminal to preserve interactivity
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start command with PTY: %w", err)
	}
	defer func() {
		_ = ptmx.Close()
	}()

	// Channel to signal when we should switch to manual mode
	manualMode := make(chan bool, 1)
	manualModeActivated := false

	// Buffer to accumulate output for pattern matching
	var outputBuffer strings.Builder

	// Start output processing goroutine
	go it.handleOutput(ptmx, &outputBuffer, &manualModeActivated, manualMode)

	// Start manual input handling goroutine
	go it.handleManualInput(ptmx, manualMode)

	// Wait for the command to complete
	waitErr := cmd.Wait()
	if waitErr != nil {
		waitErr = fmt.Errorf("failed to perform command execution: %w", waitErr)
	}

	return waitErr
}

func (it *InteractiveShellRepository) handleOutput(
	ptmx *os.File,
	outputBuffer *strings.Builder,
	manualModeActivated *bool,
	manualMode chan bool,
) {
	buf := make([]byte, bufferSize)
	for {
		n, readErr := ptmx.Read(buf)
		if readErr != nil {
			if readErr != io.EOF {
				logger.Debugf("PTY read error: %v", readErr)
			}
			break
		}

		// Filter ANSI sequences before writing to stdout
		filteredOutput := it.removeANSICodes(string(buf[:n]))
		_, _ = os.Stdout.WriteString(filteredOutput)

		// Skip pattern matching if already in manual mode
		if *manualModeActivated {
			continue
		}

		it.processOutputForPatterns(string(buf[:n]), outputBuffer, ptmx, manualModeActivated, manualMode)
	}
}

func (it *InteractiveShellRepository) processOutputForPatterns(
	output string,
	outputBuffer *strings.Builder,
	ptmx *os.File,
	manualModeActivated *bool,
	manualMode chan bool,
) {
	// Add original output to buffer for pattern matching (before ANSI filtering)
	outputBuffer.WriteString(output)

	// Check recent output for patterns (keep buffer reasonable size)
	bufferContent := outputBuffer.String()
	if len(bufferContent) > maxOutputBufferSize {
		// Keep only the last outputTrimSize characters to avoid unbounded growth
		bufferContent = bufferContent[len(bufferContent)-outputTrimSize:]
		outputBuffer.Reset()
		outputBuffer.WriteString(bufferContent)
	}

	cleanOutput := it.removeANSICodes(bufferContent)

	// Pattern 1: External dependency prompt - answer "n"
	externalDepPattern := regexp.MustCompile(
		`(?i)should terragrunt apply the external dependency.*\?`,
	)
	if externalDepPattern.MatchString(cleanOutput) {
		logger.Debug("Detected external dependency prompt, responding with 'n'")
		_, _ = ptmx.WriteString("n\r")
		outputBuffer.Reset() // Clear buffer after response
		return
	}

	// Pattern 2: "Are you sure you want to run" prompt - switch to manual mode
	confirmationPattern := regexp.MustCompile(`(?i)are you sure you want to run.*`)
	if confirmationPattern.MatchString(cleanOutput) {
		// Add newline before log messages for better formatting
		fmt.Fprintln(os.Stdout)
		logger.Info("Detected confirmation prompt, switching to manual mode")
		logger.Info("Manual interaction mode activated - user input forwarded to process")
		*manualModeActivated = true
		select {
		case manualMode <- true:
		default:
		}
		return
	}

	// Pattern 3: Any other "yes/no" prompts - answer "n" by default
	yesNoPattern := regexp.MustCompile(`(?i).*\?.*\[y/n\]`)
	if yesNoPattern.MatchString(cleanOutput) {
		logger.Debug("Detected yes/no prompt, responding with 'n'")
		_, _ = ptmx.WriteString("n\r")
		outputBuffer.Reset() // Clear buffer after response
		return
	}
}

func (it *InteractiveShellRepository) handleManualInput(ptmx *os.File, manualMode chan bool) {
	<-manualMode // Wait for signal to switch to manual mode

	// Use a more controlled input forwarding approach
	buf := make([]byte, 1)
	for {
		n, readErr := os.Stdin.Read(buf)
		if readErr != nil {
			logger.Debugf("Stdin read error in manual mode: %v", readErr)
			break
		}
		if n > 0 {
			// Forward input to PTY
			_, writeErr := ptmx.Write(buf[:n])
			if writeErr != nil {
				logger.Debugf("PTY write error in manual mode: %v", writeErr)
				break
			}
		}
	}
}

func (it *InteractiveShellRepository) removeANSICodes(text string) string {
	// Remove ANSI escape sequences - comprehensive filtering
	// Handles color codes, cursor movements, device status reports, etc.
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\([AB]|\^?\[\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(text, "")
}
