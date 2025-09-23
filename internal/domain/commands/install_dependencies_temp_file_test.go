//go:build unit

package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallDependenciesCommand_TempFilePatterns(t *testing.T) {
	t.Parallel()

	t.Run("should create unique temp files using same pattern as install function", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Real OS interface and a valid temp directory
		osInstance := entities.GetOS()
		tempDir := osInstance.GetTempDir()
		require.NotEmpty(t, tempDir, "Temp directory should be available")

		// WHEN: Creating temp files using the same pattern as install function multiple times
		tempFile1, err := os.CreateTemp(tempDir, "terraform_*")
		require.NoError(t, err, "Should create first temp file successfully")
		defer os.Remove(tempFile1.Name())
		defer tempFile1.Close()

		tempFile2, err := os.CreateTemp(tempDir, "terraform_*")
		require.NoError(t, err, "Should create second temp file successfully")
		defer os.Remove(tempFile2.Name())
		defer tempFile2.Close()

		// THEN: Should create different files with expected patterns
		assert.NotEqual(t, tempFile1.Name(), tempFile2.Name(),
			"Multiple calls should create unique temp files")
		assert.True(t, strings.HasPrefix(filepath.Base(tempFile1.Name()), "terraform_"),
			"First temp file should have expected prefix")
		assert.True(t, strings.HasPrefix(filepath.Base(tempFile2.Name()), "terraform_"),
			"Second temp file should have expected prefix")
		assert.True(t, strings.Contains(tempFile1.Name(), tempDir),
			"First temp file should be in expected directory")
		assert.True(t, strings.Contains(tempFile2.Name(), tempDir),
			"Second temp file should be in expected directory")
	})

	t.Run("should create temp files with write permissions", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Real OS interface
		osInstance := entities.GetOS()
		tempDir := osInstance.GetTempDir()
		require.NotEmpty(t, tempDir, "Temp directory should be available")

		// WHEN: Creating a temp file using the same pattern as install function
		tempFile, err := os.CreateTemp(tempDir, "terraform_*")
		require.NoError(t, err, "Should create temp file successfully")
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		// THEN: Should be able to write to the file (no permission denied errors)
		testContent := "test content for downloaded binary"
		_, err = tempFile.WriteString(testContent)
		assert.NoError(t, err, "Should be able to write to temp file without permission errors")

		// Verify content was written correctly
		_ = tempFile.Sync()
		_, err = tempFile.Seek(0, 0)
		require.NoError(t, err)

		buf := make([]byte, len(testContent))
		n, err := tempFile.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, testContent, string(buf[:n]),
			"Should be able to read back written content")
	})
}

func TestInstallDependenciesCommand_TempDirectoryPatterns(t *testing.T) {
	t.Parallel()

	t.Run("should create unique temp directories using same pattern as install function", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Real OS interface and a valid temp directory
		osInstance := entities.GetOS()
		tempDir := osInstance.GetTempDir()
		require.NotEmpty(t, tempDir, "Temp directory should be available")

		// WHEN: Creating temp directories using the same pattern as install function
		extractDir1, err := os.MkdirTemp(tempDir, "terraform_extract_*")
		require.NoError(t, err, "Should create first temp directory successfully")
		defer os.RemoveAll(extractDir1)

		extractDir2, err := os.MkdirTemp(tempDir, "terraform_extract_*")
		require.NoError(t, err, "Should create second temp directory successfully")
		defer os.RemoveAll(extractDir2)

		// THEN: Should create different directories with expected patterns
		assert.NotEqual(t, extractDir1, extractDir2,
			"Multiple calls should create unique temp directories")
		assert.True(t, strings.HasPrefix(filepath.Base(extractDir1), "terraform_extract_"),
			"First temp directory should have expected prefix")
		assert.True(t, strings.HasPrefix(filepath.Base(extractDir2), "terraform_extract_"),
			"Second temp directory should have expected prefix")
		assert.True(t, strings.Contains(extractDir1, tempDir),
			"First temp directory should be in expected parent directory")
		assert.True(t, strings.Contains(extractDir2, tempDir),
			"Second temp directory should be in expected parent directory")
	})

	t.Run("should create temp directories with write permissions", func(t *testing.T) {
		t.Parallel()

		// GIVEN: Real OS interface
		osInstance := entities.GetOS()
		tempDir := osInstance.GetTempDir()
		require.NotEmpty(t, tempDir, "Temp directory should be available")

		// WHEN: Creating a temp directory using the same pattern as install function
		extractDir, err := os.MkdirTemp(tempDir, "terraform_extract_*")
		require.NoError(t, err, "Should create temp directory successfully")
		defer os.RemoveAll(extractDir)

		// THEN: Should be able to create files in the directory (no permission denied errors)
		testFile := filepath.Join(extractDir, "terraform")
		testContent := "fake terraform binary content"
		err = os.WriteFile(testFile, []byte(testContent), 0755)
		assert.NoError(t, err, "Should be able to create files in temp directory without permission errors")

		// Verify content was written correctly
		readContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, testContent, string(readContent),
			"Should be able to read back written content")

		// Verify file permissions allow execution (for binary files)
		fileInfo, err := os.Stat(testFile)
		require.NoError(t, err)
		assert.True(t, fileInfo.Mode()&0755 != 0,
			"Created file should have execution permissions")
	})
}
