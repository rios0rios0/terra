package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockInstallDependenciesCommand is a mock implementation of the InstallDependencies interface
type MockInstallDependenciesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *MockInstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}

func TestNewInstallDependenciesController_ShouldCreateInstance_WhenCommandAndDependenciesProvided(t *testing.T) {
	// GIVEN: A mock command and test dependencies
	mockCommand := &MockInstallDependenciesCommand{}
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
}

func TestInstallDependenciesController_ShouldReturnCorrectBind_WhenGetBindCalled(t *testing.T) {
	// GIVEN: An install dependencies controller with mock command and empty dependencies
	mockCommand := &MockInstallDependenciesCommand{}
	dependencies := []entities.Dependency{}
	controller := controllers.NewInstallDependenciesController(mockCommand, dependencies)

	// WHEN: Getting the controller bind
	bind := controller.GetBind()

	// THEN: Should return correct bind configuration
	assert.Equal(t, "install", bind.Use)
	assert.Equal(t, "Install or update Terraform and Terragrunt to the latest versions", bind.Short)
	assert.Equal(t, "Install all the dependencies required to run Terra, or update them if newer versions are available. Dependencies are installed to ~/.local/bin on Linux.", bind.Long)
}

func TestInstallDependenciesController_ShouldExecuteCommand_WhenExecuteCalled(t *testing.T) {
	// GIVEN: An install dependencies controller with mock command and test dependencies
	mockCommand := &MockInstallDependenciesCommand{}
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
	assert.Equal(t, len(dependencies), len(mockCommand.LastDependencies))
	assert.Equal(t, dependencies[0].Name, mockCommand.LastDependencies[0].Name)
	assert.Equal(t, dependencies[1].Name, mockCommand.LastDependencies[1].Name)
}

func TestInstallDependenciesController_ShouldExecuteCommandMultipleTimes_WhenCalledRepeatedly(t *testing.T) {
	// GIVEN: An install dependencies controller with mock command and test dependencies
	mockCommand := &MockInstallDependenciesCommand{}
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
}
