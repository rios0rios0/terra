//go:build unit && !windows

package entities_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// zipEntry describes a single member of an archive built for a test.
type zipEntry struct {
	name    string
	content string
	mode    os.FileMode
}

// buildZipArchive writes a ZIP archive holding the given entries into a
// temporary directory and returns its path.
func buildZipArchive(t *testing.T, entries []zipEntry) string {
	t.Helper()

	archivePath := filepath.Join(t.TempDir(), "archive.zip")
	file, err := os.Create(archivePath)
	require.NoError(t, err)
	defer file.Close()

	writer := zip.NewWriter(file)
	for _, entry := range entries {
		header := &zip.FileHeader{Name: entry.name, Method: zip.Deflate}
		header.SetMode(entry.mode)

		target, createErr := writer.CreateHeader(header)
		require.NoError(t, createErr)

		_, writeErr := target.Write([]byte(entry.content))
		require.NoError(t, writeErr)
	}
	require.NoError(t, writer.Close())

	return archivePath
}

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

	t.Run("should preserve file contents when moving across directories", func(t *testing.T) {
		t.Parallel()
		// given
		srcDir := t.TempDir()
		destDir := t.TempDir()
		srcPath := filepath.Join(srcDir, "binary")
		destPath := filepath.Join(destDir, "binary")
		require.NoError(t, os.WriteFile(srcPath, []byte("payload"), 0o755))
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Move(srcPath, destPath)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(destPath)
		require.NoError(t, readErr)
		assert.Equal(t, "payload", string(content))
	})

	t.Run("should return error when source file does not exist", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Move(filepath.Join(tempDir, "nonexistent"), filepath.Join(tempDir, "dest"))

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform moving")
	})
}

func TestOSUnix_Extract(t *testing.T) {
	t.Parallel()

	t.Run("should write archive members when archive is valid", func(t *testing.T) {
		t.Parallel()
		// given
		archivePath := buildZipArchive(t, []zipEntry{
			{name: "terraform", content: "binary payload", mode: 0o755},
			{name: "README.md", content: "docs", mode: 0o644},
		})
		destDir := t.TempDir()
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Extract(archivePath, destDir)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(filepath.Join(destDir, "terraform"))
		require.NoError(t, readErr)
		assert.Equal(t, "binary payload", string(content))
		_, statErr := os.Stat(filepath.Join(destDir, "README.md"))
		require.NoError(t, statErr)
	})

	t.Run("should create parent directories owner-only when entries are nested", func(t *testing.T) {
		t.Parallel()
		// given the archive advertises a world-accessible directory
		archivePath := buildZipArchive(t, []zipEntry{
			{name: "dist/", content: "", mode: os.ModeDir | 0o777},
			{name: "dist/bin/terragrunt", content: "nested payload", mode: 0o755},
		})
		destDir := t.TempDir()
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Extract(archivePath, destDir)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(filepath.Join(destDir, "dist", "bin", "terragrunt"))
		require.NoError(t, readErr)
		assert.Equal(t, "nested payload", string(content))
		for _, dir := range []string{"dist", filepath.Join("dist", "bin")} {
			info, statErr := os.Stat(filepath.Join(destDir, dir))
			require.NoError(t, statErr)
			assert.Equal(
				t, os.FileMode(0o700), info.Mode().Perm(),
				"Extracted directory %q must not be reachable by group or other", dir,
			)
		}
	})

	t.Run("should keep the executable bit when entry is executable", func(t *testing.T) {
		t.Parallel()
		// given
		archivePath := buildZipArchive(t, []zipEntry{
			{name: "terraform", content: "binary payload", mode: 0o755},
		})
		destDir := t.TempDir()
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Extract(archivePath, destDir)

		// then
		require.NoError(t, err)
		info, statErr := os.Stat(filepath.Join(destDir, "terraform"))
		require.NoError(t, statErr)
		assert.NotZero(t, info.Mode()&0o111, "Extracted binary should stay executable")
	})

	t.Run("should unpack the members when an entry names the destination", func(t *testing.T) {
		t.Parallel()
		// given an archive opening with the `./` entry some packers emit
		archivePath := buildZipArchive(t, []zipEntry{
			{name: "./", content: "", mode: os.ModeDir | 0o755},
			{name: "terragrunt", content: "root payload", mode: 0o755},
		})
		osImpl := &entities.OSUnix{}
		destDir := t.TempDir()

		// when
		extractErr := osImpl.Extract(archivePath, destDir)

		// then
		require.NoError(t, extractErr)
		unpacked := filepath.Join(destDir, "terragrunt")
		require.FileExists(t, unpacked)
		payload, readErr := os.ReadFile(unpacked)
		require.NoError(t, readErr)
		assert.Equal(t, "root payload", string(payload))
	})

	// Every shape of a "Zip Slip" name has to be refused before the entry
	// reaches the filesystem, so they share one expectation.
	escapingEntries := map[string]string{
		"an entry escapes the destination": "../escaped.txt",
		"an entry is an absolute path":     "/etc/cron.d/escaped.txt",
		"a nested entry climbs out":        "dist/../../escaped.txt",
	}
	for description, entryName := range escapingEntries {
		t.Run("should reject the archive when "+description, func(t *testing.T) {
			t.Parallel()
			// given
			archivePath := buildZipArchive(t, []zipEntry{
				{name: entryName, content: "malicious", mode: 0o644},
			})
			destDir := t.TempDir()
			osImpl := &entities.OSUnix{}

			// when
			err := osImpl.Extract(archivePath, destDir)

			// then
			require.Error(t, err)
			assert.Contains(t, err.Error(), "escapes the destination directory")
			_, statErr := os.Stat(filepath.Join(filepath.Dir(destDir), "escaped.txt"))
			assert.True(t, os.IsNotExist(statErr), "Escaping entry must not be written")
		})
	}

	t.Run("should return error when called with non-existent archive", func(t *testing.T) {
		t.Parallel()
		// given
		osImpl := &entities.OSUnix{}
		destDir := t.TempDir()

		// when
		err := osImpl.Extract("/non/existent/archive.zip", destDir)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform decompressing")
	})

	t.Run("should return error when called with invalid archive", func(t *testing.T) {
		t.Parallel()
		// given
		osImpl := &entities.OSUnix{}
		fakePath := filepath.Join(t.TempDir(), "fake.zip")
		require.NoError(t, os.WriteFile(fakePath, []byte("not a zip"), 0o644))
		destDir := t.TempDir()

		// when
		err := osImpl.Extract(fakePath, destDir)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform decompressing")
	})
}

func TestOSUnix_Remove(t *testing.T) {
	t.Parallel()

	t.Run("should delete the file when it exists", func(t *testing.T) {
		t.Parallel()
		// given
		filePath := filepath.Join(t.TempDir(), "disposable.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("content"), 0o644))
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Remove(filePath)

		// then
		require.NoError(t, err)
		_, statErr := os.Stat(filePath)
		assert.True(t, os.IsNotExist(statErr), "File should no longer exist")
	})

	t.Run("should return error when file does not exist", func(t *testing.T) {
		t.Parallel()
		// given
		osImpl := &entities.OSUnix{}
		missingPath := filepath.Join(t.TempDir(), "nonexistent")

		// when
		err := osImpl.Remove(missingPath)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform deleting")
	})
}
