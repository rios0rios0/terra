package helpers

import (
	"os"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
)

type ArgumentsHelper struct{}

func (it ArgumentsHelper) RemovePathFromArguments(arguments []string) []string {
	_, position := findRelativePath(arguments)
	return append(arguments[:position], arguments[position+1:]...)
}

func (it ArgumentsHelper) FindAbsolutePath(arguments []string) string {
	relativePath, _ := findRelativePath(arguments)
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
