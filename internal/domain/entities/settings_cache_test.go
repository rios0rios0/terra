//go:build unit

package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings_GetModuleCacheDir(t *testing.T) {
	t.Parallel()

	t.Run("should return custom path when TerraModuleCacheDir is set", func(t *testing.T) {
		t.Parallel()
		// given
		settings := &entities.Settings{TerraModuleCacheDir: "/custom/modules"}

		// when
		dir, err := settings.GetModuleCacheDir()

		// then
		require.NoError(t, err)
		assert.Equal(t, "/custom/modules", dir)
	})

	t.Run("should return default path when TerraModuleCacheDir is empty", func(t *testing.T) {
		t.Parallel()
		// given
		settings := &entities.Settings{}

		// when
		dir, err := settings.GetModuleCacheDir()

		// then
		require.NoError(t, err)
		assert.Contains(t, dir, ".cache/terra/modules")
	})
}

func TestSettings_GetProviderCacheDir(t *testing.T) {
	t.Parallel()

	t.Run("should return custom path when TerraProviderCacheDir is set", func(t *testing.T) {
		t.Parallel()
		// given
		settings := &entities.Settings{TerraProviderCacheDir: "/custom/providers"}

		// when
		dir, err := settings.GetProviderCacheDir()

		// then
		require.NoError(t, err)
		assert.Equal(t, "/custom/providers", dir)
	})

	t.Run("should return default path when TerraProviderCacheDir is empty", func(t *testing.T) {
		t.Parallel()
		// given
		settings := &entities.Settings{}

		// when
		dir, err := settings.GetProviderCacheDir()

		// then
		require.NoError(t, err)
		assert.Contains(t, dir, ".cache/terra/providers")
	})
}
