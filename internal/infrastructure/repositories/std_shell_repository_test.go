//go:build integration

package repositories_test

import (
	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StdShellRepositorySuite struct {
	t *testing.T
}

func (suite *StdShellRepositorySuite) TestExecuteCommandSuccess(t *testing.T) {
	// given
	repository := repositories.NewStdShellRepository()
	command := "echo"
	arguments := []string{"Hello, World!"}
	directory := "."

	// when
	err := repository.ExecuteCommand(command, arguments, directory)

	// then
	assert.NoError(t, err, "should not return an error when executing a valid command")
}

func (suite *StdShellRepositorySuite) TestExecuteCommandError(t *testing.T) {
	// given
	repository := repositories.NewStdShellRepository()
	command := "nonexistentcommand"
	arguments := []string{}
	directory := "."

	// when
	err := repository.ExecuteCommand(command, arguments, directory)

	// then
	assert.Error(t, err, "should return an error when executing an invalid command")
}

func TestStdShellRepository(t *testing.T) {
	suite := &StdShellRepositorySuite{t: t}

	t.Run("should return success when executing a valid command",
		suite.TestExecuteCommandSuccess,
	)

	t.Run("should return an error when executing an invalid command",
		suite.TestExecuteCommandError,
	)
}
