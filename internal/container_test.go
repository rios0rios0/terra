//go:build unit

package internal_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func TestRegisterProviders(t *testing.T) {
	t.Parallel()

	t.Run("should register all providers without error", func(t *testing.T) {
		t.Parallel()
		// given
		container := dig.New()

		// when
		err := internal.RegisterProviders(container)

		// then
		require.NoError(t, err)
	})
}
