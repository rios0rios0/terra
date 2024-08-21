//go:build unit

package controllers_test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	testcommands "github.com/rios0rios0/terra/test/domain/commands/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InstallDependenciesControllerSuite struct {
	t *testing.T
}

func (suite *InstallDependenciesControllerSuite) TestExecuteSuccess(t *testing.T) {
	// given
	command := testcommands.NewInstallDependenciesCommandStub().WithSuccess()
	dependencies := []entities.Dependency{}
	controller := controllers.NewInstallDependenciesController(command, dependencies)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.NoError(t, err, "should not return an error when executing a valid command")
}

func (suite *InstallDependenciesControllerSuite) TestExecuteError(t *testing.T) {
	// given
	command := testcommands.NewInstallDependenciesCommandStub().WithError()
	dependencies := []entities.Dependency{}
	controller := controllers.NewInstallDependenciesController(command, dependencies)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.Error(t, err, "should return an error when executing an invalid command")
}

func TestInstallDependenciesController(t *testing.T) {
	suite := &InstallDependenciesControllerSuite{t: t}

	t.Run("should return success when installing dependencies",
		suite.TestExecuteSuccess,
	)

	t.Run("should return an error when failing to install dependencies",
		suite.TestExecuteError,
	)
}

func (suite *InstallDependenciesControllerSuite) TestGetBind(t *testing.T) {
	// given
	controller := controllers.NewInstallDependenciesController(nil, nil)
	expectedBind := entities.ControllerBind{
		Use:   "install",
		Short: "Install Terraform and Terragrunt (they are pre-requisites)",
		Long:  "Install all the dependencies required to run Terra. This command should be run with root privileges.",
	}

	// when
	bind := controller.GetBind()

	// then
	assert.Equal(t, expectedBind, bind, "expected bind is the same as the controller's bind")
}

func TestInstallDependenciesController_GetBind(t *testing.T) {
	suite := &InstallDependenciesControllerSuite{t: t}

	t.Run("should return the correct bind information",
		suite.TestGetBind,
	)
}
