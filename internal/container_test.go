//go:build unit

package internal_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func TestRegisterProviders(t *testing.T) {
	t.Parallel()

	t.Run("should register all providers without error", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		container := dig.New()

		// WHEN
		err := internal.RegisterProviders(container)

		// THEN
		require.NoError(t, err)
	})

	t.Run("should return error when called twice on same container", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container that already has all providers registered
		container := dig.New()
		firstErr := internal.RegisterProviders(container)
		require.NoError(t, firstErr)

		// WHEN: Attempting to register all providers again
		err := internal.RegisterProviders(container)

		// THEN: Should return an error due to duplicate registration
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already provided")
	})
}
