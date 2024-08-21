//go:build unit

package controllers_test

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	testcommands "github.com/rios0rios0/terra/test/domain/commands/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FormatFilesControllerSuite struct {
	t *testing.T
}

func (suite *FormatFilesControllerSuite) TestExecuteSuccess(t *testing.T) {
	// given
	command := testcommands.NewFormatFilesCommandStub().WithSuccess()
	dependencies := []entities.Dependency{}
	controller := controllers.NewFormatFilesController(command, dependencies)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.NoError(t, err, "should not return an error when executing a valid command")
}

func (suite *FormatFilesControllerSuite) TestExecuteError(t *testing.T) {
	// given
	command := testcommands.NewFormatFilesCommandStub().WithError()
	dependencies := []entities.Dependency{}
	controller := controllers.NewFormatFilesController(command, dependencies)

	// when
	err := controller.Execute(nil, nil)

	// then
	assert.Error(t, err, "should return an error when executing an invalid command")
}

func TestFormatFilesController(t *testing.T) {
	suite := &FormatFilesControllerSuite{t: t}

	t.Run("should return success when formatting files",
		suite.TestExecuteSuccess,
	)

	t.Run("should return an error when failing to format files",
		suite.TestExecuteError,
	)
}

func (suite *FormatFilesControllerSuite) TestGetBind(t *testing.T) {
	// given
	controller := controllers.NewFormatFilesController(nil, nil)
	expectedBind := entities.ControllerBind{
		Use:   "format",
		Short: "Format all files in the current directory",
		Long:  "Format all the Terraform and Terragrunt files in the current directory.",
	}

	// when
	bind := controller.GetBind()

	// then
	assert.Equal(t, expectedBind, bind, "expected bind is the same as the controller's bind")
}

func TestFormatFilesController_GetBind(t *testing.T) {
	suite := &FormatFilesControllerSuite{t: t}

	t.Run("should return the correct bind information",
		suite.TestGetBind,
	)
}
