package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/command_doubles"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelfUpdateController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when command provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A self-update command
		command := commands.NewSelfUpdateCommand()

		// WHEN: Creating a new self-update controller
		controller := controllers.NewSelfUpdateController(command)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestSelfUpdateController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct binding information when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A self-update controller
		command := commands.NewSelfUpdateCommand()
		controller := controllers.NewSelfUpdateController(command)

		// WHEN: Getting the binding information
		bind := controller.GetBind()

		// THEN: Should return correct binding details
		assert.Equal(t, "self-update", bind.Use)
		assert.Equal(t, "Update terra to the latest version", bind.Short)
		assert.Equal(t, "Download and install the latest version of terra from GitHub releases.", bind.Long)
	})
}

func TestSelfUpdateController_Execute(t *testing.T) {
	t.Run("should call command execute with correct flags when dry run flag provided", func(t *testing.T) {
		// GIVEN: A mock command and controller
		mockCommand := command_doubles.NewStubSelfUpdateCommand()
		controller := controllers.NewSelfUpdateController(mockCommand)

		// Create a cobra command with flags
		cmd := &cobra.Command{}
		cmd.Flags().Bool("dry-run", true, "test flag")
		cmd.Flags().Bool("force", false, "test flag")

		// WHEN: Executing the controller
		// Note: This will call the real command which will fail due to GitHub API
		// but we can verify the controller structure is correct
		controller.Execute(cmd, []string{})

		// THEN: The mock command should have been called
		assert.True(t, mockCommand.WasCalled())
		assert.True(t, mockCommand.DryRunFlag)
		assert.False(t, mockCommand.ForceFlag)
	})

	t.Run("should call command execute with correct flags when force flag provided", func(t *testing.T) {
		// GIVEN: A mock command and controller
		mockCommand := command_doubles.NewStubSelfUpdateCommand()
		controller := controllers.NewSelfUpdateController(mockCommand)

		// Create a cobra command with flags
		cmd := &cobra.Command{}
		cmd.Flags().Bool("dry-run", false, "test flag")
		cmd.Flags().Bool("force", true, "test flag")

		// WHEN: Executing the controller
		controller.Execute(cmd, []string{})

		// THEN: The mock command should have been called with correct flags
		assert.True(t, mockCommand.WasCalled())
		assert.False(t, mockCommand.DryRunFlag)
		assert.True(t, mockCommand.ForceFlag)
	})
}
