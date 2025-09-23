package repositories

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
	logger "github.com/sirupsen/logrus"
)

const (
	outputChannelSize   = 100
	outputCheckInterval = 100 * time.Millisecond
	shellTimeout        = 30 * time.Minute // Allow long-running terraform/terragrunt commands
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

	// Read output byte by byte and forward to stdout while monitoring for patterns
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					logger.Debugf("PTY read error: %v", err)
				}
				break
			}

			// Filter ANSI sequences before writing to stdout
			filteredOutput := it.removeANSICodes(string(buf[:n]))
			_, _ = os.Stdout.WriteString(filteredOutput)

			// Skip pattern matching if already in manual mode
			if manualModeActivated {
				continue
			}

			// Add original output to buffer for pattern matching (before ANSI filtering)
			outputBuffer.Write(buf[:n])
			
			// Check recent output for patterns (keep buffer reasonable size)
			output := outputBuffer.String()
			if len(output) > 4096 {
				// Keep only the last 2048 characters to avoid unbounded growth
				output = output[len(output)-2048:]
				outputBuffer.Reset()
				outputBuffer.WriteString(output)
			}

			cleanOutput := it.removeANSICodes(output)
			
			// Pattern 1: External dependency prompt - answer "n"
			externalDepPattern := regexp.MustCompile(
				`(?i)should terragrunt apply the external dependency.*\?`,
			)
			if externalDepPattern.MatchString(cleanOutput) {
				logger.Debug("Detected external dependency prompt, responding with 'n'")
				_, _ = ptmx.Write([]byte("n\r"))
				outputBuffer.Reset() // Clear buffer after response
				continue
			}

			// Pattern 2: "Are you sure you want to run" prompt - switch to manual mode
			confirmationPattern := regexp.MustCompile(`(?i)are you sure you want to run.*`)
			if confirmationPattern.MatchString(cleanOutput) {
				logger.Info("Detected confirmation prompt, switching to manual mode")
				logger.Info("Manual interaction mode activated - user input forwarded to process")
				manualModeActivated = true
				select {
				case manualMode <- true:
				default:
				}
				continue
			}

			// Pattern 3: Any other "yes/no" prompts - answer "n" by default
			yesNoPattern := regexp.MustCompile(`(?i).*\?.*\[y/n\]`)
			if yesNoPattern.MatchString(cleanOutput) {
				logger.Debug("Detected yes/no prompt, responding with 'n'")
				_, _ = ptmx.Write([]byte("n\r"))
				outputBuffer.Reset() // Clear buffer after response
				continue
			}
		}
	}()

	// Handle manual input when needed with more careful input forwarding
	go func() {
		<-manualMode // Wait for signal to switch to manual mode
		
		// Use a more controlled input forwarding approach
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				logger.Debugf("Stdin read error in manual mode: %v", err)
				break
			}
			if n > 0 {
				// Forward input to PTY
				_, err = ptmx.Write(buf[:n])
				if err != nil {
					logger.Debugf("PTY write error in manual mode: %v", err)
					break
				}
			}
		}
	}()

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		err = fmt.Errorf("failed to perform command execution: %w", err)
	}

	return err
}

func (it *InteractiveShellRepository) removeANSICodes(text string) string {
	// Remove ANSI escape sequences - comprehensive filtering
	// Handles color codes, cursor movements, device status reports, etc.
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\([AB]|\^?\[\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(text, "")
}
