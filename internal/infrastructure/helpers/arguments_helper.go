package helpers

import (
	"os"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
)

type ArgumentsHelper struct{}

func (it ArgumentsHelper) RemovePathFromArguments(arguments []string) []string {
	_, position := findRelativePath(arguments)
	if position == -1 {
		// No valid directory found, return arguments as-is
		return arguments
	}
	return append(arguments[:position], arguments[position+1:]...)
}

func (it ArgumentsHelper) FindAbsolutePath(arguments []string) string {
	relativePath, position := findRelativePath(arguments)

	// If no valid directory was found in arguments, check if the resolved path is valid
	if position == -1 {
		// User specified a path but it wasn't found to be valid
		// Check if any of the arguments look like a path (contains "/" or looks like directory)
		for _, arg := range arguments {
			if arg != "" && (filepath.IsAbs(arg) || filepath.Clean(arg) != arg || arg == "." || arg == ".." ||
				(len(arg) > 0 && (arg[0] == '/' || arg[0] == '.' || arg[0] == '~'))) {
				// This looks like a path argument, validate it
				if info, err := os.Stat(arg); err != nil {
					if os.IsNotExist(err) {
						logger.Fatalf("Directory does not exist: %s", arg)
					} else {
						logger.Fatalf("Cannot access directory: %s - %s", arg, err)
					}
				} else if !info.IsDir() {
					logger.Fatalf("Path is not a directory: %s", arg)
				}
			}
		}
		// Check if current directory is accessible (fallback to ".")
		if _, err := os.Stat(relativePath); err != nil {
			logger.Fatalf("No valid directory found in arguments and current directory is not accessible: %s", err)
		}
	} else {
		// Check if the specified path exists and is a directory
		if info, err := os.Stat(relativePath); err != nil {
			if os.IsNotExist(err) {
				logger.Fatalf("Directory does not exist: %s", relativePath)
			} else {
				logger.Fatalf("Cannot access directory: %s - %s", relativePath, err)
			}
		} else if !info.IsDir() {
			logger.Fatalf("Path is not a directory: %s", relativePath)
		}
	}

	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		logger.Fatalf("Error resolving directory path: %s", err)
	}
	return absolutePath
}

func findRelativePath(arguments []string) (string, int) {
	position := -1
	relativePath := "."

	// check if the first or last argument is a directory
	if _, err := os.Stat(arguments[0]); err == nil {
		position = 0
		relativePath = arguments[position]
	} else if _, err := os.Stat(arguments[len(arguments)-1]); err == nil {
		position = len(arguments) - 1
		relativePath = arguments[position]
	}
	return relativePath, position
}
