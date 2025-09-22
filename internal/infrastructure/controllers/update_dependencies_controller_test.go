//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdateDependenciesController(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Valid dependencies
		dependencies := []entities.Dependency{
			{Name: "terraform", CLI: "terraform"},
			{Name: "terragrunt", CLI: "terragrunt"},
		}
		installCommand := &commanddoubles.StubInstallDependencies{}

		// WHEN: Creating a new update controller
		controller := controllers.NewUpdateDependenciesController(installCommand, dependencies)

		// THEN: Should create instance
		assert.NotNil(t, controller)
	})
}

func TestUpdateDependenciesController_GetBind(t *testing.T) {
	t.Parallel()

	t.Run("should return correct bind when called", func(t *testing.T) {
		t.Parallel()

		// GIVEN: An update controller
		dependencies := []entities.Dependency{}
		installCommand := &commanddoubles.StubInstallDependencies{}
		controller := controllers.NewUpdateDependenciesController(installCommand, dependencies)

		// WHEN: Getting the bind
		bind := controller.GetBind()

		// THEN: Should return correct bind information
		assert.Equal(t, "update", bind.Use)
		assert.Equal(t, "Install or update Terraform and Terragrunt to the latest versions", bind.Short)
		assert.Contains(t, bind.Long, "This is an alias for the 'install' command")
	})
}

func TestUpdateDependenciesController_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute install command when called", func(t *testing.T) {
		t.Parallel()

		// GIVEN: An update controller with dependencies
		dependencies := []entities.Dependency{
			{Name: "terraform", CLI: "terraform"},
			{Name: "terragrunt", CLI: "terragrunt"},
		}
		installCommand := &commanddoubles.StubInstallDependencies{}
		controller := controllers.NewUpdateDependenciesController(installCommand, dependencies)

		// WHEN: Executing the controller
		controller.Execute(nil, []string{})

		// THEN: Should call install command with dependencies
		assert.True(t, installCommand.ExecuteCalled)
		assert.Equal(t, dependencies, installCommand.LastDependencies)
	})
}