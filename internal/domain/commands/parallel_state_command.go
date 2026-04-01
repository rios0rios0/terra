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
// Returns true if --parallel=N flag is present.
func (it *ParallelStateCommand) shouldExecuteInParallel(arguments []string) bool {
	return HasParallelFlag(arguments)
}

// removeParallelFlags removes --parallel=N, --only=, --skip=, and --reply/-r flags from arguments.
func (it *ParallelStateCommand) removeParallelFlags(arguments []string) []string {
	filtered := RemoveParallelFlag(arguments)
	filtered = RemoveSelectionFlags(filtered)
	return RemoveReplyFlag(filtered)
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

// buildOnlyPaths validates and returns full paths for the given --only values.
func (it *ParallelStateCommand) buildOnlyPaths(
	targetPath string,
	onlyModules []string,
) []string {
	var paths []string

	for _, module := range onlyModules {
		fullPath := filepath.Join(targetPath, module)

		info, err := os.Stat(fullPath)
		if err == nil && info.IsDir() {
			paths = append(paths, fullPath)
		} else {
			logger.Warnf("Module path does not exist or is not a directory: %s", fullPath)
		}
	}

	return paths
}

// applySkips removes paths matching any --skip value from the given path list.
func (it *ParallelStateCommand) applySkips(
	targetPath string,
	paths []string,
	skipModules []string,
) []string {
	var filteredPaths []string

	for _, currentPath := range paths {
		relPath, err := filepath.Rel(targetPath, currentPath)
		if err != nil {
			relPath = filepath.Base(currentPath)
		}

		if !it.isSkipped(relPath, skipModules) {
			filteredPaths = append(filteredPaths, currentPath)
		}
	}

	return filteredPaths
}

// isSkipped checks whether a relative path matches any --skip entry.
func (it *ParallelStateCommand) isSkipped(relPath string, skipModules []string) bool {
	for _, skip := range skipModules {
		if relPath == skip || filepath.Base(relPath) == skip {
			return true
		}
	}

	return false
}

// buildSelectedPaths builds full paths from the given --only/--skip values.
func (it *ParallelStateCommand) buildSelectedPaths(
	targetPath string,
	selection SelectionValues,
) []string {
	var paths []string

	if len(selection.Only) > 0 {
		paths = it.buildOnlyPaths(targetPath, selection.Only)
	} else {
		allModules, err := it.findSubdirectories(targetPath)
		if err != nil {
			logger.Warnf("Failed to find subdirectories for --skip filter: %s", err)
			return paths
		}

		paths = allModules
	}

	if len(selection.Skip) > 0 {
		paths = it.applySkips(targetPath, paths, selection.Skip)
	}

	return paths
}

// resolveModules discovers which module directories to process based on --only/--skip flags.
func (it *ParallelStateCommand) resolveModules(
	targetPath string,
	arguments []string,
) ([]string, error) {
	selection := GetSelectionValues(arguments)
	hasSelection := len(selection.Only) > 0 || len(selection.Skip) > 0

	if hasSelection {
		modules := it.buildSelectedPaths(targetPath, selection)
		if len(modules) == 0 {
			return nil, errors.New("no valid module paths found for --only/--skip selection")
		}

		logger.Infof("Using module selection: %d modules to process", len(modules))

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
		wg.Go(func() {
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
		})
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

	hadReplyFlag := HasReplyFlag(arguments)
	filteredArguments := it.removeParallelFlags(arguments)

	// When --reply was provided, inject --non-interactive so terragrunt skips prompts
	// in each worker (parallel workers cannot handle stdin prompts)
	if hadReplyFlag {
		filteredArguments = append(filteredArguments, "--non-interactive")
	}

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
