//go:build integration

package repositories_test

import (
	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdShellRepository_ExecuteCommand(t *testing.T) {
	t.Run("should return success when executing a valid command", func(t *testing.T) {
		// given
		repository := repositories.NewStdShellRepository()

		// when
		err := repository.ExecuteCommand("echo", []string{"Hello, World!"}, ".")

		// then
		assert.NoError(t, err, "should not return an error when executing a valid command")
	})

	t.Run("should return an error when executing an invalid command", func(t *testing.T) {
		// given
		repository := repositories.NewStdShellRepository()

		// when
		err := repository.ExecuteCommand("non-existent", []string{}, ".")

		// then
		assert.Error(t, err, "should return an error when executing an invalid command")
	})
}
