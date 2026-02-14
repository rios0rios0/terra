//go:build unit && !windows

package entities_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSUnix_Move(t *testing.T) {
	t.Parallel()

	t.Run("should move file successfully when valid paths provided", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		srcPath := filepath.Join(tempDir, "source.txt")
		destPath := filepath.Join(tempDir, "dest.txt")
		require.NoError(t, os.WriteFile(srcPath, []byte("content"), 0o644))
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Move(srcPath, destPath)

		// then
		require.NoError(t, err)
		_, statErr := os.Stat(destPath)
		assert.False(t, os.IsNotExist(statErr), "Destination file should exist")
		_, statErr = os.Stat(srcPath)
		assert.True(t, os.IsNotExist(statErr), "Source file should no longer exist")
	})

	t.Run("should return error when source file does not exist", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Move(filepath.Join(tempDir, "nonexistent"), filepath.Join(tempDir, "dest"))

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mv")
	})
}
