//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestRegisterProviders(t *testing.T) {
	t.Parallel()

	t.Run("should register all providers when container is valid", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		container := dig.New()

		// WHEN
		err := commands.RegisterProviders(container)

		// THEN
		assert.NoError(t, err)
	})
}
