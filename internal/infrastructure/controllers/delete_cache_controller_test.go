package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteCacheController(t *testing.T) {
	t.Parallel()
	
	t.Run("should create instance when command provided", func(t *testing.T) {
		// GIVEN: A mock delete cache command
		mockCommand := &test.StubDeleteCacheCommand{}

		// WHEN: Creating a new delete cache controller
		controller := controllers.NewDeleteCacheController(mockCommand)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestDeleteCacheController_GetBind(t *testing.T) {
	t.Parallel()
	
	t.Run("should return correct bind when called", func(t *testing.T) {
		// GIVEN: A delete cache controller
		mockCommand := &test.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)

		// WHEN: Getting the controller bind
		bind := controller.GetBind()

		// THEN: Should return correct bind configuration
		assert.Equal(t, "clear", bind.Use)
		assert.Equal(t, "Clear all cache and modules directories", bind.Short)
		assert.Equal(t, "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.", bind.Long)
	})
}

func TestDeleteCacheController_Execute(t *testing.T) {
	t.Parallel()
	
	t.Run("should execute command when called", func(t *testing.T) {
		// GIVEN: A delete cache controller with mock command
		mockCommand := &test.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command with correct directories
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
		expectedDirs := []string{".terraform", ".terragrunt-cache"}
		assert.Equal(t, expectedDirs, mockCommand.LastToBeDeleted)
	})
	
	t.Run("should execute command multiple times when called repeatedly", func(t *testing.T) {
		// GIVEN: A delete cache controller with mock command
		mockCommand := &test.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)
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
