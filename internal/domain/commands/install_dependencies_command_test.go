//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/domain/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Run("should install dependencies if not available", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{Name: "Go", CLI: "go", VersionURL: "https://golang.org/dl/", RegexVersion: `go(\d+\.\d+\.\d+)`, BinaryURL: "https://golang.org/dl/go%s.linux-amd64.tar.gz"},
			{Name: "Python", CLI: "python", VersionURL: "https://www.python.org/downloads/", RegexVersion: `Python (\d+\.\d+\.\d+)`, BinaryURL: "https://www.python.org/ftp/python/%s/Python-%s.tgz"},
		}
		repository := repositories.NewShellRepositoryStub().WithSuccess()
		command := commands.NewInstallDependenciesCommand()

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.True(t, repository.CommandExecuted(dependency.CLI, "install"), "command should be executed for: "+dependency.CLI)
		}
	})

	t.Run("should handle already installed dependencies gracefully", func(t *testing.T) {
		// given
		dependencies := []entities.Dependency{
			{Name: "Go", CLI: "go", VersionURL: "https://golang.org/dl/", RegexVersion: `go(\d+\.\d+\.\d+)`, BinaryURL: "https://golang.org/dl/go%s.linux-amd64.tar.gz"},
		}
		repository := repositories.NewShellRepositoryStub().WithError()
		command := commands.NewInstallDependenciesCommand()

		// when
		command.Execute(dependencies)

		// then
		for _, dependency := range dependencies {
			assert.False(t, repository.CommandExecuted(dependency.CLI, "install"), "command should not be executed for: "+dependency.CLI)
		}
	})
}

// Mock functions for testing
func fetchLatestVersion(url, regexPattern string) string {
	// Mock implementation for testing
	return "1.0.0"
}

func isDependencyCLIAvailable(name string) bool {
	// Mock implementation for testing
	return false
}

func ensureRootPrivileges() {
	// Mock implementation for testing
}

func install(url, name string) {
	// Mock implementation for testing
}
