package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

// shouldExecuteInParallel determines if the command should be executed in parallel.
// It returns true if either:
// 1. It's a state command with --all flag (backward compatibility)
// 2. It has --parallel=N flag (new functionality for any command)
func (it *ParallelStateCommand) shouldExecuteInParallel(arguments []string) bool {
	// New: support parallel=N for any command
	if HasParallelFlag(arguments) {
		return true
	}
	// Backward compatibility: state commands with --all flag
	return IsStateManipulationCommand(arguments) && HasAllFlag(arguments)
}

// removeAllFlag removes --all flag from arguments since terragrunt doesn't support it for state commands.
func (it *ParallelStateCommand) removeAllFlag(arguments []string) []string {
	var filtered []string
	for _, arg := range arguments {
		if arg != AllFlag {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// removeParallelFlags removes --all, --parallel=N, --no-parallel-bypass, and --filter= flags from arguments.
func (it *ParallelStateCommand) removeParallelFlags(arguments []string) []string {
	// First remove --all flag
	filtered := it.removeAllFlag(arguments)
	// Then remove --parallel=N flag
	filtered = RemoveParallelFlag(filtered)
	// Remove --no-parallel-bypass flag
	filtered = RemoveNoParallelBypassFlag(filtered)
	// Finally remove --filter= flag
	return RemoveFilterFlag(filtered)
}

// findSubdirectories finds all subdirectories that contain terraform/terragrunt files.
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

// containsTerraformFiles checks if a directory contains terraform or terragrunt files.
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

// buildFilteredPaths builds full paths by concatenating filter values with the target path.
// Handles both inclusions and exclusions (values starting with !).
func (it *ParallelStateCommand) buildFilteredPaths(targetPath string, filterValues []string) []string {
	parsed := ParseFilterValues(filterValues)
	var paths []string

	// If there are inclusions, use only those
	if len(parsed.Inclusions) > 0 {
		for _, inclusion := range parsed.Inclusions {
			fullPath := filepath.Join(targetPath, inclusion)
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				paths = append(paths, fullPath)
			} else {
				logger.Warnf("Filter path does not exist or is not a directory: %s", fullPath)
			}
		}
	} else {
		// No inclusions, find all subdirectories
		allModules, err := it.findSubdirectories(targetPath)
		if err != nil {
			logger.Warnf("Failed to find subdirectories for exclusion filter: %s", err)
			return paths
		}
		paths = allModules
	}

	// Remove exclusions from the paths
	if len(parsed.Exclusions) > 0 {
		var filteredPaths []string
		for _, path := range paths {
			// Get the relative path from targetPath to check if it matches any exclusion
			relPath, err := filepath.Rel(targetPath, path)
			if err != nil {
				// If we can't get relative path, use the last component
				relPath = filepath.Base(path)
			}
			
			// Check if this path should be excluded
			shouldExclude := false
			for _, exclusion := range parsed.Exclusions {
				// Check if the relative path or its base name matches the exclusion
				if relPath == exclusion || filepath.Base(relPath) == exclusion {
					shouldExclude = true
					break
				}
			}
			
			if !shouldExclude {
				filteredPaths = append(filteredPaths, path)
			}
		}
		paths = filteredPaths
	}

	return paths
}

// executeInParallel executes the command in parallel across multiple directories.
func (it *ParallelStateCommand) executeInParallel(
	targetPath string,
	arguments []string,
	maxJobs int,
) error {
	// Record start time for execution duration
	startTime := time.Now()

	var modules []string
	var err error

	// Check if --filter flag is present
	if filterValues, hasFilter := GetFilterValue(arguments); hasFilter {
		// Use filtered directories
		modules = it.buildFilteredPaths(targetPath, filterValues)
		if len(modules) == 0 {
			return fmt.Errorf("no valid filter paths found")
		}
		logger.Infof("Using filter: %d modules to process", len(modules))
	} else {
		// Find all subdirectories with terraform files
		modules, err = it.findSubdirectories(targetPath)
		if err != nil {
			return err
		}
		logger.Infof("Found %d modules to process", len(modules))
	}

	// Adjust maxJobs to not exceed the number of modules
	if maxJobs > len(modules) {
		maxJobs = len(modules)
		logger.Infof("Reducing thread count to %d (number of modules)", maxJobs)
	}

	// Remove --all and --parallel=N flags from arguments for individual module execution
	filteredArguments := it.removeParallelFlags(arguments)

	// Create channels for parallel execution
	jobs := make(chan string, len(modules))
	results := make(chan error, len(modules))

	// Start worker goroutines
	var wg sync.WaitGroup
	for range maxJobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for modulePath := range jobs {
				logger.Infof("==> Processing %s", modulePath)
				executeErr := it.repository.ExecuteCommand("terragrunt", filteredArguments, modulePath)
				if executeErr != nil {
					logger.Errorf("✗ %s: %s", modulePath, executeErr)
					results <- fmt.Errorf("module %s failed: %w", modulePath, executeErr)
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
	var executeErrors []error
	for err := range results {
		if err != nil {
			executeErrors = append(executeErrors, err)
		}
	}

	// Calculate execution duration
	duration := time.Since(startTime)

	// Report summary with threads and duration
	successful := len(modules) - len(executeErrors)
	logger.Infof("Parallel execution completed: %d successful, %d failed (threads: %d, duration: %s)", 
		successful, len(executeErrors), maxJobs, duration.Round(time.Millisecond))

	if len(executeErrors) > 0 {
		// Return first error for simplicity, but log all errors
		for _, err := range executeErrors {
			logger.Error(err)
		}
		return fmt.Errorf("parallel execution failed with %d errors", len(executeErrors))
	}

	return nil
}

func (it *ParallelStateCommand) Execute(
	targetPath string,
	arguments []string,
	_ []entities.Dependency,
) error {
	// Check if this should be executed in parallel
	if !it.shouldExecuteInParallel(arguments) {
		// Not a parallel command, should not reach here
		return errors.New("command is not a parallel command")
	}

	// Determine max jobs: use parallel=N value if present, otherwise default
	maxJobs := defaultMaxJobs
	if parallelValue, found := GetParallelValue(arguments); found {
		maxJobs = parallelValue
		logger.Infof("Using %d parallel threads", maxJobs)
	}

	// Execute in parallel
	return it.executeInParallel(targetPath, arguments, maxJobs)
}
