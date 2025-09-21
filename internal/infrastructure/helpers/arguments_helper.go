package helpers

import (
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
)

type ArgumentsHelper struct{}

func (it ArgumentsHelper) RemovePathFromArguments(arguments []string) []string {
	_, position := findRelativePath(arguments)
	if position == -1 {
		return arguments
	}
	return append(arguments[:position], arguments[position+1:]...)
}

func (it ArgumentsHelper) FindAbsolutePath(arguments []string) string {
	relativePath, _ := findRelativePath(arguments)
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		logger.Fatalf("Error resolving directory path: %s", err)
	}

	// Validate that the path exists and is a directory
	info, err := os.Stat(absolutePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Fatalf("Directory does not exist: %s", absolutePath)
		}
		logger.Fatalf("Error accessing directory %s: %s", absolutePath, err)
	}

	if !info.IsDir() {
		logger.Fatalf("Path is not a directory: %s", absolutePath)
	}

	return absolutePath
}

func findRelativePath(arguments []string) (string, int) {
	position := -1
	relativePath := "."

	if len(arguments) == 0 {
		return relativePath, position
	}

	// Helper function to check if an argument looks like a path
	isPathLike := func(arg string) bool {
		return strings.HasPrefix(arg, "/") || strings.HasPrefix(arg, "./") ||
			strings.HasPrefix(arg, "../") ||
			arg == "."
	}

	// Check if the first argument is a directory (existing) or looks like a path
	if _, err := os.Stat(arguments[0]); err == nil {
		position = 0
		relativePath = arguments[position]
	} else if isPathLike(arguments[0]) {
		position = 0
		relativePath = arguments[position]
	} else if len(arguments) > 1 {
		// Check the last argument
		lastArg := arguments[len(arguments)-1]
		if _, err := os.Stat(lastArg); err == nil {
			position = len(arguments) - 1
			relativePath = lastArg
		} else if isPathLike(lastArg) {
			position = len(arguments) - 1
			relativePath = lastArg
		}
	}
	return relativePath, position
}
