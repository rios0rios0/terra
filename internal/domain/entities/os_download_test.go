//go:build unit && !windows

package entities_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
