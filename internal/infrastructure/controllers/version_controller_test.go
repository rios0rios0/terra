//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/command_doubles"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVersionController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when command provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock version command
		mockCommand := &command_doubles.StubVersionCommand{}

		// WHEN: Creating a new version controller
		controller := controllers.NewVersionController(mockCommand)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestVersionController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct bind when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A version controller with mock command
		mockCommand := &command_doubles.StubVersionCommand{}
		controller := controllers.NewVersionController(mockCommand)

		// WHEN: Getting the controller bind
		bind := controller.GetBind()

		// THEN: Should return correct bind configuration
		assert.Equal(t, "version", bind.Use)
		assert.Equal(t, "Show Terra, Terraform, and Terragrunt versions", bind.Short)
		assert.Equal(
			t,
			"Display the version information for Terra and its dependencies (Terraform and Terragrunt).",
			bind.Long,
		)
	})
}

func TestVersionController_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute command when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A version controller with mock command
		mockCommand := &command_doubles.StubVersionCommand{}
		controller := controllers.NewVersionController(mockCommand)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command once
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
	})

	t.Run("should execute command multiple times when called repeatedly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A version controller with mock command
		mockCommand := &command_doubles.StubVersionCommand{}
		controller := controllers.NewVersionController(mockCommand)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller multiple times
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)

		// THEN: Should execute the command the correct number of times
		assert.Equal(t, 2, mockCommand.ExecuteCallCount)
	})
}
