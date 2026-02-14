//go:build unit

package internal_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubController is a minimal Controller implementation for testing.
type stubController struct{}

func (s *stubController) GetBind() entities.ControllerBind {
	return entities.ControllerBind{Use: "test", Short: "test command", Long: "test long"}
}

func (s *stubController) Execute(_ *cobra.Command, _ []string) {}

func TestNewAppInternal(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when controllers provided", func(t *testing.T) {
		t.Parallel()
		// given
		controllers := &[]entities.Controller{}

		// when
		app := internal.NewAppInternal(controllers)

		// then
		require.NotNil(t, app)
	})
}

func TestAppInternal_GetControllers(t *testing.T) {
	t.Parallel()

	t.Run("should return empty slice when no controllers provided", func(t *testing.T) {
		t.Parallel()
		// given
		controllers := &[]entities.Controller{}
		app := internal.NewAppInternal(controllers)

		// when
		result := app.GetControllers()

		// then
		assert.Empty(t, result)
	})

	t.Run("should return controllers when provided", func(t *testing.T) {
		t.Parallel()
		// given
		ctrl := &stubController{}
		controllers := &[]entities.Controller{ctrl}
		app := internal.NewAppInternal(controllers)

		// when
		result := app.GetControllers()

		// then
		assert.Len(t, result, 1)
		assert.Equal(t, "test", result[0].GetBind().Use)
	})
}
