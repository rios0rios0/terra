//go:build unit

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"should return 0 when versions are equal", "1.5.0", "1.5.0", 0},
		{"should return -1 when v1 is older major", "1.0.0", "2.0.0", -1},
		{"should return 1 when v1 is newer major", "2.0.0", "1.0.0", 1},
		{"should return -1 when v1 is older minor", "1.4.0", "1.5.0", -1},
		{"should return 1 when v1 is newer minor", "1.6.0", "1.5.0", 1},
		{"should return -1 when v1 is older patch", "1.5.0", "1.5.1", -1},
		{"should return 1 when v1 is newer patch", "1.5.2", "1.5.1", 1},
		{"should handle different length versions", "1.5", "1.5.0", 0},
		{"should handle different length versions with diff", "1.5", "1.5.1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.CompareVersionsPublic(tt.v1, tt.v2))
		})
	}
}

func TestFindBinaryInArchive(t *testing.T) {
	t.Parallel()

	t.Run("should find exact binary match when present", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		binaryPath := filepath.Join(tempDir, "terraform")
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.WriteFile(binaryPath, []byte("binary"), 0755))

		// when
		found, err := commands.FindBinaryInArchivePublic(tempDir, "terraform")

		// then
		require.NoError(t, err)
		assert.Equal(t, binaryPath, found)
	})

	t.Run("should find binary in nested directory when present", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		nestedDir := filepath.Join(tempDir, "subdir")
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(nestedDir, 0o750))
		binaryPath := filepath.Join(nestedDir, "terraform")
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.WriteFile(binaryPath, []byte("binary"), 0755))

		// when
		found, err := commands.FindBinaryInArchivePublic(tempDir, "terraform")

		// then
		require.NoError(t, err)
		assert.Equal(t, binaryPath, found)
	})

	t.Run("should return error when binary not found", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()

		// when
		_, err := commands.FindBinaryInArchivePublic(tempDir, "nonexistent")

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find nonexistent binary")
	})

	t.Run("should find pattern match binary when present", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		binaryPath := filepath.Join(tempDir, "terraform_1_5_0_linux_amd64")
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.WriteFile(binaryPath, []byte("binary"), 0755))

		// when
		found, err := commands.FindBinaryInArchivePublic(tempDir, "terraform")

		// then
		require.NoError(t, err)
		assert.Equal(t, binaryPath, found)
	})
}
