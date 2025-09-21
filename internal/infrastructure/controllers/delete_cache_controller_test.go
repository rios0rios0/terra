package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDeleteCacheCommand is a mock implementation of the DeleteCache interface
type MockDeleteCacheCommand struct {
	ExecuteCallCount int
	LastToBeDeleted  []string
}

func (m *MockDeleteCacheCommand) Execute(toBeDeleted []string) {
	m.ExecuteCallCount++
	m.LastToBeDeleted = toBeDeleted
}

func TestNewDeleteCacheController_ShouldCreateInstance_WhenCommandProvided(t *testing.T) {
	// GIVEN: A mock delete cache command
	mockCommand := &MockDeleteCacheCommand{}

	// WHEN: Creating a new delete cache controller
	controller := controllers.NewDeleteCacheController(mockCommand)

	// THEN: Should create a valid controller instance
	require.NotNil(t, controller)
}

func TestDeleteCacheController_ShouldReturnCorrectBind_WhenGetBindCalled(t *testing.T) {
	// GIVEN: A delete cache controller
	mockCommand := &MockDeleteCacheCommand{}
	controller := controllers.NewDeleteCacheController(mockCommand)

	// WHEN: Getting the controller bind
	bind := controller.GetBind()

	// THEN: Should return correct bind configuration
	assert.Equal(t, "clear", bind.Use)
	assert.Equal(t, "Clear all cache and modules directories", bind.Short)
	assert.Equal(t, "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.", bind.Long)
}

func TestDeleteCacheController_ShouldExecuteCommand_WhenExecuteCalled(t *testing.T) {
	// GIVEN: A delete cache controller with mock command
	mockCommand := &MockDeleteCacheCommand{}
	controller := controllers.NewDeleteCacheController(mockCommand)
	cmd := &cobra.Command{}
	args := []string{}

	// WHEN: Executing the controller
	controller.Execute(cmd, args)

	// THEN: Should execute the command with correct directories
	assert.Equal(t, 1, mockCommand.ExecuteCallCount)
	expectedDirs := []string{".terraform", ".terragrunt-cache"}
	assert.Equal(t, expectedDirs, mockCommand.LastToBeDeleted)
}

func TestDeleteCacheController_ShouldExecuteCommandMultipleTimes_WhenCalledRepeatedly(t *testing.T) {
	// GIVEN: A delete cache controller with mock command
	mockCommand := &MockDeleteCacheCommand{}
	controller := controllers.NewDeleteCacheController(mockCommand)
	cmd := &cobra.Command{}
	args := []string{}

	// WHEN: Executing the controller multiple times
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)

	// THEN: Should execute the command the correct number of times
	assert.Equal(t, 3, mockCommand.ExecuteCallCount)
}
