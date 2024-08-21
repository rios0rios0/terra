//go:build unit

package controllers_test

import (
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RunFromRootControllerSuite struct {
	t *testing.T
}

func (suite *RunFromRootControllerSuite) TestExecuteSuccess(t *testing.T) {
	// given
	command := testcommands.NewRunFromRootCommandStub().OnSuccess()
	dependencies := []entities.Dependency{}
	controller := controllers.NewRunFromRootController(command, dependencies)
	arguments := []string{"arg1", "arg2"}

	// when
	err := controller.Execute(nil, arguments)

	// then
	assert.NoError(t, err, "should not return an error when executing a valid command")
}

func (suite *RunFromRootControllerSuite) TestExecuteError(t *testing.T) {
	// given
	command := testcommands.NewRunFromRootCommandStub().OnError()
	dependencies := []entities.Dependency{}
	controller := controllers.NewRunFromRootController(command, dependencies)
	arguments := []string{"arg1", "arg2"}

	// when
	err := controller.Execute(nil, arguments)

	// then
	assert.Error(t, err, "should return an error when executing an invalid command")
}

func TestRunFromRootController(t *testing.T) {
	suite := &RunFromRootControllerSuite{t: t}

	t.Run("should return success when running from root",
		suite.TestExecuteSuccess,
	)

	t.Run("should return an error when failing to run from root",
		suite.TestExecuteError,
	)
}

func (suite *RunFromRootControllerSuite) TestGetBind(t *testing.T) {
	// given
	controller := controllers.NewRunFromRootController(nil, nil)
	expectedBind := entities.ControllerBind{
		Use:   "terra [flags] [terragrunt command] [directory]",
		Short: "Terra is a CLI wrapper for Terragrunt",
		Long: "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. " +
			"It also allows changing the account/subscription and workspace for AWS and Azure.",
	}

	// when
	bind := controller.GetBind()

	// then
	assert.Equal(t, expectedBind, bind, "expected bind is the same as the controller's bind")
}

func TestRunFromRootController_GetBind(t *testing.T) {
	suite := &RunFromRootControllerSuite{t: t}

	t.Run("should return the correct bind information",
		suite.TestGetBind,
	)
}
