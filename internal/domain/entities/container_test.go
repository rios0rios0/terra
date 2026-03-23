//go:build unit

package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
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
		err := entities.RegisterProviders(container)

		// THEN
		assert.NoError(t, err)
	})
}
