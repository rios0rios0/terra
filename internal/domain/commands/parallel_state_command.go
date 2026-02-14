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
// 2. It has --parallel=N flag (new functionality for any command).
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

// buildInclusionPaths validates and returns full paths for the given inclusion filter values.
func (it *ParallelStateCommand) buildInclusionPaths(
	targetPath string,
	inclusions []string,
) []string {
	var paths []string

	for _, inclusion := range inclusions {
		fullPath := filepath.Join(targetPath, inclusion)

		info, err := os.Stat(fullPath)
		if err == nil && info.IsDir() {
			paths = append(paths, fullPath)
		} else {
			logger.Warnf("Filter path does not exist or is not a directory: %s", fullPath)
		}
	}

	return paths
}

// applyExclusions removes paths matching any exclusion from the given path list.
func (it *ParallelStateCommand) applyExclusions(
	targetPath string,
	paths []string,
	exclusions []string,
) []string {
	var filteredPaths []string

	for _, currentPath := range paths {
		relPath, err := filepath.Rel(targetPath, currentPath)
		if err != nil {
			relPath = filepath.Base(currentPath)
		}

		if !it.isExcluded(relPath, exclusions) {
			filteredPaths = append(filteredPaths, currentPath)
		}
	}

	return filteredPaths
}

// isExcluded checks whether a relative path matches any exclusion entry.
func (it *ParallelStateCommand) isExcluded(relPath string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if relPath == exclusion || filepath.Base(relPath) == exclusion {
			return true
		}
	}

	return false
}

// buildFilteredPaths builds full paths by concatenating filter values with the target path.
// Handles both inclusions and exclusions (values starting with !).
func (it *ParallelStateCommand) buildFilteredPaths(
	targetPath string,
	filterValues []string,
) []string {
	parsed := ParseFilterValues(filterValues)
	var paths []string

	if len(parsed.Inclusions) > 0 {
		paths = it.buildInclusionPaths(targetPath, parsed.Inclusions)
	} else {
		allModules, err := it.findSubdirectories(targetPath)
		if err != nil {
			logger.Warnf("Failed to find subdirectories for exclusion filter: %s", err)
			return paths
		}

		paths = allModules
	}

	if len(parsed.Exclusions) > 0 {
		paths = it.applyExclusions(targetPath, paths, parsed.Exclusions)
	}

	return paths
}

// resolveModules discovers which module directories to process based on arguments and filters.
func (it *ParallelStateCommand) resolveModules(
	targetPath string,
	arguments []string,
) ([]string, error) {
	if filterValues, hasFilter := GetFilterValue(arguments); hasFilter {
		modules := it.buildFilteredPaths(targetPath, filterValues)
		if len(modules) == 0 {
			return nil, errors.New("no valid filter paths found")
		}

		logger.Infof("Using filter: %d modules to process", len(modules))

		return modules, nil
	}

	modules, err := it.findSubdirectories(targetPath)
	if err != nil {
		return nil, err
	}

	logger.Infof("Found %d modules to process", len(modules))

	return modules, nil
}

// runWorkers spawns goroutine workers that execute the command across modules and collects errors.
func (it *ParallelStateCommand) runWorkers(
	modules []string,
	filteredArguments []string,
	maxJobs int,
) []error {
	jobs := make(chan string, len(modules))
	results := make(chan error, len(modules))

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

	go func() {
		defer close(jobs)
		for _, module := range modules {
			jobs <- module
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var executeErrors []error
	for err := range results {
		if err != nil {
			executeErrors = append(executeErrors, err)
		}
	}

	return executeErrors
}

// executeInParallel executes the command in parallel across multiple directories.
func (it *ParallelStateCommand) executeInParallel(
	targetPath string,
	arguments []string,
	maxJobs int,
) error {
	startTime := time.Now()

	modules, err := it.resolveModules(targetPath, arguments)
	if err != nil {
		return err
	}

	if maxJobs > len(modules) {
		maxJobs = len(modules)
		logger.Infof("Reducing thread count to %d (number of modules)", maxJobs)
	}

	filteredArguments := it.removeParallelFlags(arguments)
	executeErrors := it.runWorkers(modules, filteredArguments, maxJobs)
	duration := time.Since(startTime)

	successful := len(modules) - len(executeErrors)
	logger.Infof(
		"Parallel execution completed: %d successful, %d failed (threads: %d, duration: %s)",
		successful, len(executeErrors), maxJobs, duration.Round(time.Millisecond),
	)

	if len(executeErrors) > 0 {
		for _, workerErr := range executeErrors {
			logger.Error(workerErr)
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
