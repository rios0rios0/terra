//go:build unit

package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgumentsHelper_RemovePathFromArguments(t *testing.T) {
	t.Parallel()
	helper := helpers.ArgumentsHelper{}

	t.Run("should remove first argument when it is an existing directory", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		args := []string{tempDir, "plan", "--detailed-exitcode"}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, result)
	})

	t.Run("should remove last argument when it is an existing directory", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		args := []string{"plan", tempDir}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{"plan"}, result)
	})

	t.Run("should return unchanged when no path found", func(t *testing.T) {
		t.Parallel()
		// given
		args := []string{"plan", "--detailed-exitcode"}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{"plan", "--detailed-exitcode"}, result)
	})

	t.Run("should return unchanged when empty arguments", func(t *testing.T) {
		t.Parallel()
		// given
		args := []string{}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{}, result)
	})

	t.Run("should detect path-like first argument starting with ./", func(t *testing.T) {
		t.Parallel()
		// given
		args := []string{"./some/path", "plan"}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{"plan"}, result)
	})

	t.Run("should detect path-like first argument starting with /", func(t *testing.T) {
		t.Parallel()
		// given
		args := []string{"/some/path", "plan"}

		// when
		result := helper.RemovePathFromArguments(args)

		// then
		assert.Equal(t, []string{"plan"}, result)
	})
}

func TestArgumentsHelper_FindAbsolutePath(t *testing.T) {
	t.Run("should return absolute path when existing directory provided as first argument", func(t *testing.T) {
		// given
		tempDir := t.TempDir()
		helper := helpers.ArgumentsHelper{}
		args := []string{tempDir, "plan"}

		// when
		result := helper.FindAbsolutePath(args)

		// then
		expected, err := filepath.Abs(tempDir)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return current directory when no path in arguments", func(t *testing.T) {
		// given
		helper := helpers.ArgumentsHelper{}
		args := []string{"plan", "--detailed-exitcode"}

		// when
		result := helper.FindAbsolutePath(args)

		// then
		cwd, err := os.Getwd()
		require.NoError(t, err)
		assert.Equal(t, cwd, result)
	})

	t.Run("should return absolute path when existing directory provided as last argument", func(t *testing.T) {
		// given
		tempDir := t.TempDir()
		helper := helpers.ArgumentsHelper{}
		args := []string{"plan", tempDir}

		// when
		result := helper.FindAbsolutePath(args)

		// then
		expected, err := filepath.Abs(tempDir)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
