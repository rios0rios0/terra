//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteCacheController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when command provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock delete cache command
		mockCommand := &commanddoubles.StubDeleteCacheCommand{}

		// WHEN: Creating a new delete cache controller
		controller := controllers.NewDeleteCacheController(mockCommand)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestDeleteCacheController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct bind when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache controller
		mockCommand := &commanddoubles.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)

		// WHEN: Getting the controller bind
		bind := controller.GetBind()

		// THEN: Should return correct bind configuration
		assert.Equal(t, "clear", bind.Use)
		assert.Equal(t, "Clear all cache and modules directories", bind.Short)
		assert.Contains(t, bind.Long, "Clear all temporary directories")
		assert.Contains(t, bind.Long, "--global")
	})
}

func TestDeleteCacheController_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute command when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache controller with mock command
		mockCommand := &commanddoubles.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)
		cmd := &cobra.Command{} //nolint:exhaustruct // minimal test setup
		cmd.Flags().Bool("global", false, "")
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command with correct directories and global=false
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
		expectedDirs := []string{".terraform", ".terragrunt-cache"}
		assert.Equal(t, expectedDirs, mockCommand.LastToBeDeleted)
		assert.False(t, mockCommand.LastGlobal)
	})

	t.Run("should pass global flag as true when set", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache controller with mock command and global flag set to true
		mockCommand := &commanddoubles.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)
		cmd := &cobra.Command{} //nolint:exhaustruct // minimal test setup
		cmd.Flags().Bool("global", false, "")
		require.NoError(t, cmd.Flags().Set("global", "true"))
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command with correct directories and global=true
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
		expectedDirs := []string{".terraform", ".terragrunt-cache"}
		assert.Equal(t, expectedDirs, mockCommand.LastToBeDeleted)
		assert.True(t, mockCommand.LastGlobal)
	})

	t.Run("should execute command multiple times when called repeatedly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A delete cache controller with mock command
		mockCommand := &commanddoubles.StubDeleteCacheCommand{}
		controller := controllers.NewDeleteCacheController(mockCommand)
		cmd := &cobra.Command{} //nolint:exhaustruct // minimal test setup
		cmd.Flags().Bool("global", false, "")
		args := []string{}

		// WHEN: Executing the controller multiple times
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)

		// THEN: Should execute the command the correct number of times
		assert.Equal(t, 3, mockCommand.ExecuteCallCount)
	})
}
