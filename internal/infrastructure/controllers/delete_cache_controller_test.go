//go:build unit

package controllers_test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	testcommands "github.com/rios0rios0/terra/test/domain/commands/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DeleteCacheControllerSuite struct {
	t *testing.T
}

func (suite *DeleteCacheControllerSuite) TestExecuteSuccess(t *testing.T) {
	// given
	command := testcommands.NewDeleteCacheCommandStub().WithSuccess()
	controller := controllers.NewDeleteCacheController(command)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.NoError(t, err, "should not return an error when executing a valid command")
}

func (suite *DeleteCacheControllerSuite) TestExecuteError(t *testing.T) {
	// given
	command := testcommands.NewDeleteCacheCommandStub().WithError()
	controller := controllers.NewDeleteCacheController(command)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.Error(t, err, "should return an error when executing an invalid command")
}

func TestDeleteCacheController(t *testing.T) {
	suite := &DeleteCacheControllerSuite{t: t}

	t.Run("should return success when clearing cache",
		suite.TestExecuteSuccess,
	)

	t.Run("should return an error when failing to clear cache",
		suite.TestExecuteError,
	)
}

func (suite *DeleteCacheControllerSuite) TestGetBind(t *testing.T) {
	// given
	controller := controllers.NewDeleteCacheController(nil)
	expectedBind := entities.ControllerBind{
		Use:   "clear",
		Short: "Clear all cache and modules directories",
		Long:  "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution.",
	}

	// when
	bind := controller.GetBind()

	// then
	assert.Equal(t, expectedBind, bind, "expected bind is the same as the controller's bind")
}

func TestDeleteCacheController_GetBind(t *testing.T) {
	suite := &DeleteCacheControllerSuite{t: t}

	t.Run("should return the correct bind information",
		suite.TestGetBind,
	)
}
