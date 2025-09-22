//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatFilesController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when command and dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock command and test dependencies
		mockCommand := &commanddoubles.StubFormatFilesCommand{}
		dependencies := []entities.Dependency{
			{
				Name:              "Test Tool",
				CLI:               "test",
				FormattingCommand: []string{"format", "-recursive"},
			},
		}

		// WHEN: Creating a new format files controller
		controller := controllers.NewFormatFilesController(mockCommand, dependencies)

		// THEN: Should create a valid controller instance
		require.NotNil(t, controller)
	})
}

func TestFormatFilesController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct bind when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A format files controller with mock command and empty dependencies
		mockCommand := &commanddoubles.StubFormatFilesCommand{}
		dependencies := []entities.Dependency{}
		controller := controllers.NewFormatFilesController(mockCommand, dependencies)

		// WHEN: Getting the controller bind
		bind := controller.GetBind()

		// THEN: Should return correct bind configuration
		assert.Equal(t, "format", bind.Use)
		assert.Equal(t, "Format all files in the current directory", bind.Short)
		assert.Equal(
			t,
			"Format all the Terraform and Terragrunt files in the current directory.",
			bind.Long,
		)
	})
}

func TestFormatFilesController_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute command when called with dependencies", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A format files controller with mock command and test dependencies
		mockCommand := &commanddoubles.StubFormatFilesCommand{}
		terraformDep := entities.Dependency{
			Name: "Terraform",
			CLI:  "terraform",
		}
		terragruntDep := entities.Dependency{
			Name: "Terragrunt",
			CLI:  "terragrunt",
		}
		dependencies := []entities.Dependency{terraformDep, terragruntDep}
		controller := controllers.NewFormatFilesController(mockCommand, dependencies)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller
		controller.Execute(cmd, args)

		// THEN: Should execute the command with correct dependencies
		assert.Equal(t, 1, mockCommand.ExecuteCallCount)
		require.Len(t, mockCommand.LastDependencies, 2)
		assert.Equal(t, terraformDep.Name, mockCommand.LastDependencies[0].Name)
		assert.Equal(t, terragruntDep.Name, mockCommand.LastDependencies[1].Name)
	})

	t.Run("should execute command multiple times when called repeatedly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A format files controller with mock command and empty dependencies
		mockCommand := &commanddoubles.StubFormatFilesCommand{}
		dependencies := []entities.Dependency{}
		controller := controllers.NewFormatFilesController(mockCommand, dependencies)
		cmd := &cobra.Command{}
		args := []string{}

		// WHEN: Executing the controller multiple times
		controller.Execute(cmd, args)
		controller.Execute(cmd, args)

		// THEN: Should execute the command the correct number of times
		assert.Equal(t, 2, mockCommand.ExecuteCallCount)
	})
}
