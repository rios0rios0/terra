//go:build unit && !windows

package entities_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const downloadTimeoutEnvVar = "TERRA_DOWNLOAD_TIMEOUT"

func TestOSUnix_Download(t *testing.T) {
	t.Parallel()

	t.Run("should download file successfully when valid URL provided", func(t *testing.T) {
		t.Parallel()
		// given
		expectedContent := "hello-binary-content"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(expectedContent))
		}))
		defer server.Close()

		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "downloaded_file")
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Download(server.URL, destPath)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(destPath)
		require.NoError(t, readErr)
		assert.Equal(t, expectedContent, string(content))
	})

	t.Run("should return error when server returns non-200 status", func(t *testing.T) {
		t.Parallel()
		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "downloaded_file")
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Download(server.URL, destPath)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("should return error when URL is unreachable", func(t *testing.T) {
		t.Parallel()
		// given
		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "downloaded_file")
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Download("http://127.0.0.1:1/unreachable", destPath)

		// then
		assert.Error(t, err)
	})

}

// `TestOSUnix_DownloadTimeout` sets `TERRA_DOWNLOAD_TIMEOUT` via
// `t.Setenv`, which requires the test (and its parent) to be
// non-parallel. Kept in a separate top-level function instead of as
// subtests of `TestOSUnix_Download` (which calls `t.Parallel()`)
// because Go rejects `Setenv` inside any test marked Parallel.
func TestOSUnix_DownloadTimeout(t *testing.T) {
	t.Run("should honor TERRA_DOWNLOAD_TIMEOUT when the response is slow", func(t *testing.T) {
		// given
		// Server sleeps longer than the configured override, so a
		// correctly applied deadline will abort the body read with
		// `context deadline exceeded` before the server replies.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("would-be-content"))
		}))
		defer server.Close()

		t.Setenv(downloadTimeoutEnvVar, "100ms")

		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "downloaded_file")
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Download(server.URL, destPath)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "deadline exceeded")
	})

	t.Run("should fall back to default when TERRA_DOWNLOAD_TIMEOUT is malformed", func(t *testing.T) {
		// given
		// A malformed value must not break installs -- the function
		// logs a warning and uses the default 10-minute timeout, so
		// the fast test server still completes well inside the
		// fallback ceiling.
		expectedContent := "fallback-still-works"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(expectedContent))
		}))
		defer server.Close()

		t.Setenv(downloadTimeoutEnvVar, "this-is-not-a-duration")

		tempDir := t.TempDir()
		destPath := filepath.Join(tempDir, "downloaded_file")
		osImpl := &entities.OSUnix{}

		// when
		err := osImpl.Download(server.URL, destPath)

		// then
		require.NoError(t, err)
		content, readErr := os.ReadFile(destPath)
		require.NoError(t, readErr)
		assert.Equal(t, expectedContent, string(content))
	})
}
