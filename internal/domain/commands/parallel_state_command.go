package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	logger "github.com/sirupsen/logrus"
)

const defaultMaxJobs = 5

type ParallelStateCommand struct {
	repository repositories.ShellRepository
}

func NewParallelStateCommand(repository repositories.ShellRepository) *ParallelStateCommand {
	return &ParallelStateCommand{
		repository: repository,
	}
}

// isStateManipulationCommand checks if the command is a state manipulation command
func (it *ParallelStateCommand) isStateManipulationCommand(arguments []string) bool {
	if len(arguments) == 0 {
		return false
	}

	// Check for state manipulation commands
	stateCommands := []string{
		"import", "state",
	}

	firstArg := arguments[0]
	for _, cmd := range stateCommands {
		if firstArg == cmd {
			return true
		}
	}

	// Check for state subcommands (e.g., "state rm", "state mv")
	if len(arguments) >= 2 && firstArg == "state" {
		stateSubcommands := []string{
			"rm", "mv", "pull", "push", "show",
		}
		secondArg := arguments[1]
		for _, subcmd := range stateSubcommands {
			if secondArg == subcmd {
				return true
			}
		}
	}

	return false
}

// hasAllFlag checks if --all flag is present in arguments
func (it *ParallelStateCommand) hasAllFlag(arguments []string) bool {
	for _, arg := range arguments {
		if arg == "--all" {
			return true
		}
	}
	return false
}

// shouldExecuteInParallel determines if the command should be executed in parallel
func (it *ParallelStateCommand) shouldExecuteInParallel(arguments []string) bool {
	return it.isStateManipulationCommand(arguments) && it.hasAllFlag(arguments)
}

// removeAllFlag removes --all flag from arguments since terragrunt doesn't support it for state commands
func (it *ParallelStateCommand) removeAllFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != "--all" {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// findSubdirectories finds all subdirectories that contain terraform/terragrunt files
func (it *ParallelStateCommand) findSubdirectories(rootPath string) ([]string, error) {
	var modules []string

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !d.IsDir() {
			return nil
		}

		// Skip hidden directories and the root directory itself
		if strings.HasPrefix(d.Name(), ".") || path == rootPath {
			return nil
		}

		// Check if directory contains terraform files
		if it.containsTerraformFiles(path) {
			modules = append(modules, path)
			// Don't traverse deeper into this directory
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find subdirectories: %w", err)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("no terraform/terragrunt modules found in %s", rootPath)
	}

	return modules, nil
}

// containsTerraformFiles checks if a directory contains terraform or terragrunt files
func (it *ParallelStateCommand) containsTerraformFiles(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".tf") ||
			strings.HasSuffix(name, ".tfvars") ||
			name == "terragrunt.hcl" {
			return true
		}
	}

	return false
}

// executeInParallel executes the command in parallel across multiple directories
func (it *ParallelStateCommand) executeInParallel(
	targetPath string,
	arguments []string,
	maxJobs int,
) error {
	// Find all subdirectories with terraform files
	modules, err := it.findSubdirectories(targetPath)
	if err != nil {
		return err
	}

	logger.Infof("Found %d modules to process", len(modules))

	// Remove --all flag from arguments for individual module execution
	filteredArguments := it.removeAllFlag(arguments)

	// Create channels for parallel execution
	jobs := make(chan string, len(modules))
	results := make(chan error, len(modules))

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < maxJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for modulePath := range jobs {
				logger.Infof("==> Processing %s", modulePath)
				err := it.repository.ExecuteCommand("terragrunt", filteredArguments, modulePath)
				if err != nil {
					logger.Errorf("✗ %s: %s", modulePath, err)
					results <- fmt.Errorf("module %s failed: %w", modulePath, err)
				} else {
					logger.Infof("✓ %s", modulePath)
					results <- nil
				}
			}
		}()
	}

	// Send jobs to workers
	go func() {
		defer close(jobs)
		for _, module := range modules {
			jobs <- module
		}
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var errors []error
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Report summary
	successful := len(modules) - len(errors)
	logger.Infof("Parallel execution completed: %d successful, %d failed", successful, len(errors))

	if len(errors) > 0 {
		// Return first error for simplicity, but log all errors
		for _, err := range errors {
			logger.Error(err)
		}
		return fmt.Errorf("parallel execution failed with %d errors", len(errors))
	}

	return nil
}

func (it *ParallelStateCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) error {
	// Check if this should be executed in parallel
	if !it.shouldExecuteInParallel(arguments) {
		// Not a parallel state command, should not reach here
		return fmt.Errorf("command is not a parallel state manipulation command")
	}

	// Execute in parallel with default max jobs
	return it.executeInParallel(targetPath, arguments, defaultMaxJobs)
}
