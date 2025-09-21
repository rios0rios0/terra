package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRunFromRootCommand is a mock implementation of the RunFromRoot interface
type MockRunFromRootCommand struct {
	ExecuteCallCount int
	LastTargetPath   string
	LastArguments    []string
	LastDependencies []entities.Dependency
}

func (m *MockRunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	m.ExecuteCallCount++
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
	m.LastDependencies = dependencies
}

func TestNewRunFromRootController_ShouldCreateInstance_WhenCommandAndDependenciesProvided(t *testing.T) {
	// GIVEN: A mock command and test dependencies
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Test Dependency",
			CLI:  "test",
		},
	}

	// WHEN: Creating a new run from root controller
	controller := controllers.NewRunFromRootController(mockCommand, dependencies)

	// THEN: Should create a valid controller instance
	require.NotNil(t, controller)
}

func TestRunFromRootController_ShouldReturnCorrectBind_WhenGetBindCalled(t *testing.T) {
	// GIVEN: A run from root controller with mock command and empty dependencies
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}
	controller := controllers.NewRunFromRootController(mockCommand, dependencies)

	// WHEN: Getting the controller bind
	bind := controller.GetBind()

	// THEN: Should return correct bind configuration
	assert.Equal(t, "terra [flags] [terragrunt command] [directory]", bind.Use)
	assert.Equal(t, "Terra is a CLI wrapper for Terragrunt", bind.Short)
	assert.Equal(t, "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. It also allows changing the account/subscription and workspace for AWS and Azure.", bind.Long)
}

func TestRunFromRootController_ShouldExecuteCommand_WhenExecuteCalled(t *testing.T) {
	// GIVEN: A run from root controller with mock command and test dependencies
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Terraform",
			CLI:  "terraform",
		},
	}
	controller := controllers.NewRunFromRootController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{"apply", "."}

	// WHEN: Executing the controller
	controller.Execute(cmd, args)

	// THEN: Should execute the command with correct dependencies and arguments
	assert.Equal(t, 1, mockCommand.ExecuteCallCount)
	assert.Equal(t, len(dependencies), len(mockCommand.LastDependencies))
	assert.Equal(t, dependencies[0].Name, mockCommand.LastDependencies[0].Name)
	assert.NotEmpty(t, mockCommand.LastArguments)
}

func TestRunFromRootController_ShouldExecuteCommand_WhenDifferentArgumentsProvided(t *testing.T) {
	// GIVEN: A run from root controller with mock command and empty dependencies
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}
	controller := controllers.NewRunFromRootController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{"plan", "--dry-run"}

	// WHEN: Executing the controller with different arguments
	controller.Execute(cmd, args)

	// THEN: Should execute the command with the provided arguments
	assert.Equal(t, 1, mockCommand.ExecuteCallCount)
	assert.NotEmpty(t, mockCommand.LastArguments)
}

func TestRunFromRootController_ShouldExecuteCommandMultipleTimes_WhenCalledRepeatedly(t *testing.T) {
	// GIVEN: A run from root controller with mock command and empty dependencies
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}
	controller := controllers.NewRunFromRootController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{"plan"}

	// WHEN: Executing the controller multiple times
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)

	// THEN: Should execute the command the correct number of times
	assert.Equal(t, 3, mockCommand.ExecuteCallCount)
}
