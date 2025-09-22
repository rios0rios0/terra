package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/command_doubles"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallDependenciesController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when command and dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock command and test dependencies
		mockCommand := &command_doubles.StubInstallDependenciesCommand{}
		dependencies := []entities.Dependency{
			{
				Name:              "Test Dependency",
				CLI:               "test",
				BinaryURL:         "https://example.com/test",
				VersionURL:        "https://example.com/version",
				RegexVersion:      `"version":"([^"]+)"`,
				FormattingCommand: []string{"format"},
			},
		}

		// WHEN: Creating a new install dependencies controller
		controller := controllers.NewInstallDependenciesController(mockCommand, dependencies)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestInstallDependenciesController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct bind when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An install dependencies controller with mock command and empty dependencies
		mockCommand := &command_doubles.StubInstallDependenciesCommand{}
		dependencies := []entities.Dependency{}
		controller := controllers.NewInstallDependenciesController(mockCommand, dependencies)

		// WHEN: Getting the controller bind
		bind := controller.GetBind()

		// THEN: Should return correct bind configuration
		assert.Equal(t, "install", bind.Use)
		assert.Equal(
			t,
			"Install or update Terraform and Terragrunt to the latest versions",
			bind.Short,
		)
		assert.Equal(
			t,
			"Install all the dependencies required to run Terra, or update them if newer versions are available. Dependencies are installed to ~/.local/bin on Linux.",
			bind.Long,
		)
	})
}

func TestInstallDependenciesController_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute command when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An install dependencies controller with mock command and test dependencies
		mockCommand := &command_doubles.StubInstallDependenciesCommand{}
		dependencies := []entities.Dependency{
			{
				Name: "Test Dependency",
				CLI:  "test",
			},
			{
				Name: "Another Dependency",
				CLI:  "another",
			},
		}
		controller := controllers.NewInstallDependenciesController(mockCommand, dependencies)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command with correct dependencies
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
		assert.Len(t, mockCommand.LastDependencies, len(dependencies))
		assert.Equal(t, dependencies[0].Name, mockCommand.LastDependencies[0].Name)
		assert.Equal(t, dependencies[1].Name, mockCommand.LastDependencies[1].Name)
	})

	t.Run("should execute command multiple times when called repeatedly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An install dependencies controller with mock command and test dependencies
		mockCommand := &command_doubles.StubInstallDependenciesCommand{}
		dependencies := []entities.Dependency{
			{Name: "Test", CLI: "test"},
		}
		controller := controllers.NewInstallDependenciesController(mockCommand, dependencies)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller multiple times
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)

		// THEN: Should execute the command the correct number of times
		assert.Equal(t, 3, mockCommand.ExecuteCallCount)
	})
}
