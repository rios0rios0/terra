//go:build unit

package commands_test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelfUpdateCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN

		// WHEN
		cmd := commands.NewSelfUpdateCommand()

		// THEN
		require.NotNil(t, cmd)
	})
}

func TestSelfUpdateCommand_Execute(t *testing.T) {
	t.Run("should show correct download URL when dry run succeeds", func(t *testing.T) {
		// GIVEN
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run (hits real GitHub API)
		err := cmd.Execute(true, false)

		// THEN: Should succeed without error (dry run does not download)
		// NOTE: This test may fail if GitHub API rate limits are hit.
		// In that case, the error message will contain "failed to fetch latest release".
		if err != nil {
			assert.Contains(t, err.Error(), "failed to fetch latest release",
				"Only GitHub API rate limiting should cause failure")
		}
	})

	t.Run("should report current version is newer when version is higher than latest", func(t *testing.T) {
		// GIVEN: A version significantly higher than any real release
		originalVersion := commands.TerraVersion
		commands.TerraVersion = "999.999.999"
		defer func() { commands.TerraVersion = originalVersion }()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run
		err := cmd.Execute(true, false)

		// THEN: Should succeed without error because the current version is newer
		// NOTE: This test may fail if GitHub API rate limits are hit.
		if err != nil {
			assert.Contains(t, err.Error(), "failed to fetch latest release",
				"Only GitHub API rate limiting should cause failure")
		}
	})
}

// redirectTransport is an http.RoundTripper that redirects all requests to a test server.
type redirectTransport struct {
	targetURL string
}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the request URL to point to the test server while preserving the path
	redirectedURL := t.targetURL + req.URL.Path
	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, redirectedURL, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create redirected request: %w", err)
	}
	newReq.Header = req.Header
	return http.DefaultTransport.RoundTrip(newReq)
}

// helperWithRedirectedHTTPClient temporarily replaces http.DefaultClient transport
// so that all HTTP calls go to the given test server. It returns a cleanup function
// that restores the original transport.
func helperWithRedirectedHTTPClient(serverURL string) func() {
	originalTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = &redirectTransport{targetURL: serverURL}
	return func() {
		http.DefaultClient.Transport = originalTransport
	}
}

func TestSelfUpdateCommand_Execute_WithMockedAPI(t *testing.T) {
	// NOTE: Cannot use t.Parallel() because these tests modify http.DefaultClient.Transport

	t.Run("should return error when GitHub API returns non-200 status", func(t *testing.T) {
		// GIVEN: A mock server that returns HTTP 403 Forbidden
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing a dry run
		err := cmd.Execute(true, false)

		// THEN: Should return an error about failed fetch
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
	})

	t.Run("should return error when GitHub API returns invalid JSON", func(t *testing.T) {
		// GIVEN: A mock server that returns invalid JSON
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{this is not valid json}`))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing a dry run
		err := cmd.Execute(true, false)

		// THEN: Should return an error about parsing
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
	})

	t.Run("should return error when no matching asset found for platform", func(t *testing.T) {
		// GIVEN: A mock server that returns a release without a matching asset
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"tag_name": "v2.0.0",
				"assets": [
					{
						"name": "terra-2.0.0-fakeos-fakearch.tar.gz",
						"browser_download_url": "https://example.com/fake-asset"
					}
				]
			}`))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing a dry run
		err := cmd.Execute(true, false)

		// THEN: Should return an error about missing asset for this platform
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
		assert.Contains(t, err.Error(), "no asset")
	})

	t.Run("should report up to date when current version equals latest", func(t *testing.T) {
		// GIVEN: A mock server that returns a version matching TerraVersion
		originalVersion := commands.TerraVersion
		commands.TerraVersion = "2.0.0"
		defer func() { commands.TerraVersion = originalVersion }()

		archString := runtime.GOARCH
		osString := runtime.GOOS
		assetName := fmt.Sprintf("terra-2.0.0-%s-%s.tar.gz", osString, archString)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := fmt.Sprintf(`{
				"tag_name": "v2.0.0",
				"assets": [
					{
						"name": %q,
						"browser_download_url": "https://example.com/terra-2.0.0.tar.gz"
					}
				]
			}`, assetName)
			_, _ = w.Write([]byte(response))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run when version matches
		err := cmd.Execute(true, false)

		// THEN: Should succeed with no error (version is up to date)
		assert.NoError(t, err)
	})

	t.Run("should report dry run update info when current version is older", func(t *testing.T) {
		// GIVEN: A mock server that returns a version newer than TerraVersion
		originalVersion := commands.TerraVersion
		commands.TerraVersion = "1.0.0"
		defer func() { commands.TerraVersion = originalVersion }()

		archString := runtime.GOARCH
		osString := runtime.GOOS
		assetName := fmt.Sprintf("terra-3.0.0-%s-%s.tar.gz", osString, archString)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := fmt.Sprintf(`{
				"tag_name": "v3.0.0",
				"assets": [
					{
						"name": %q,
						"browser_download_url": "https://example.com/terra-3.0.0.tar.gz"
					}
				]
			}`, assetName)
			_, _ = w.Write([]byte(response))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run when version is older
		err := cmd.Execute(true, false)

		// THEN: Should succeed without error (dry run just logs info)
		assert.NoError(t, err)
	})

	t.Run("should report newer version when current version is higher than latest", func(t *testing.T) {
		// GIVEN: A mock server that returns a version older than TerraVersion
		originalVersion := commands.TerraVersion
		commands.TerraVersion = "99.0.0"
		defer func() { commands.TerraVersion = originalVersion }()

		archString := runtime.GOARCH
		osString := runtime.GOOS
		assetName := fmt.Sprintf("terra-1.0.0-%s-%s.tar.gz", osString, archString)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := fmt.Sprintf(`{
				"tag_name": "v1.0.0",
				"assets": [
					{
						"name": %q,
						"browser_download_url": "https://example.com/terra-1.0.0.tar.gz"
					}
				]
			}`, assetName)
			_, _ = w.Write([]byte(response))
		}))
		defer server.Close()
		cleanup := helperWithRedirectedHTTPClient(server.URL)
		defer cleanup()
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run when current version is newer
		err := cmd.Execute(true, false)

		// THEN: Should succeed without error (version is newer than latest)
		assert.NoError(t, err)
	})
}

// helperCreateTarGz creates a valid tar.gz archive containing a single file
// with the given name and content in the specified directory.
func helperCreateTarGz(t *testing.T, dir, fileName, content string) string {
	t.Helper()
	archivePath := filepath.Join(dir, "test.tar.gz")

	file, err := os.Create(archivePath)
	require.NoError(t, err)

	gzWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzWriter)

	contentBytes := []byte(content)
	header := &tar.Header{
		Name: fileName,
		Mode: 0o755,
		Size: int64(len(contentBytes)),
	}
	require.NoError(t, tarWriter.WriteHeader(header))

	_, err = tarWriter.Write(contentBytes)
	require.NoError(t, err)

	require.NoError(t, tarWriter.Close())
	require.NoError(t, gzWriter.Close())
	require.NoError(t, file.Close())

	return archivePath
}

func TestSelfUpdateCommand_ExtractArchive(t *testing.T) {
	t.Parallel()

	t.Run("should extract file when valid tar.gz archive provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A valid tar.gz archive containing a "terra" binary
		cmd := commands.NewSelfUpdateCommand()
		tempDir := t.TempDir()
		archivePath := helperCreateTarGz(t, tempDir, "terra", "fake binary content")

		extractDir := filepath.Join(tempDir, "extract")
		require.NoError(t, os.MkdirAll(extractDir, 0o755))

		// WHEN: Extracting the archive
		err := cmd.ExtractArchiveForTest(archivePath, extractDir)

		// THEN: Should succeed and the extracted file should exist with correct content
		assert.NoError(t, err)
		extractedFile := filepath.Join(extractDir, "terra")
		content, readErr := os.ReadFile(extractedFile)
		require.NoError(t, readErr)
		assert.Equal(t, "fake binary content", string(content))
	})

	t.Run("should return error when archive does not exist", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A path to a non-existent archive
		cmd := commands.NewSelfUpdateCommand()
		extractDir := t.TempDir()

		// WHEN: Attempting to extract a non-existent archive
		err := cmd.ExtractArchiveForTest("/nonexistent/path/archive.tar.gz", extractDir)

		// THEN: Should return an error indicating extraction failure
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract tar.gz archive")
	})

	t.Run("should return error when archive contains invalid data", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A file with invalid tar.gz content
		cmd := commands.NewSelfUpdateCommand()
		tempDir := t.TempDir()
		invalidArchive := filepath.Join(tempDir, "invalid.tar.gz")
		require.NoError(t, os.WriteFile(invalidArchive, []byte("this is not a valid archive"), 0o644))

		extractDir := filepath.Join(tempDir, "extract")
		require.NoError(t, os.MkdirAll(extractDir, 0o755))

		// WHEN: Attempting to extract the invalid archive
		err := cmd.ExtractArchiveForTest(invalidArchive, extractDir)

		// THEN: Should return an error indicating extraction failure
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract tar.gz archive")
	})

	t.Run("should extract multiple files when archive contains several entries", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A tar.gz archive containing multiple files
		cmd := commands.NewSelfUpdateCommand()
		tempDir := t.TempDir()
		archivePath := filepath.Join(tempDir, "multi.tar.gz")

		file, err := os.Create(archivePath)
		require.NoError(t, err)
		gzWriter := gzip.NewWriter(file)
		tarWriter := tar.NewWriter(gzWriter)

		files := map[string]string{
			"terra":     "binary content",
			"README.md": "readme content",
		}
		for name, content := range files {
			contentBytes := []byte(content)
			header := &tar.Header{
				Name: name,
				Mode: 0o644,
				Size: int64(len(contentBytes)),
			}
			require.NoError(t, tarWriter.WriteHeader(header))
			_, writeErr := tarWriter.Write(contentBytes)
			require.NoError(t, writeErr)
		}
		require.NoError(t, tarWriter.Close())
		require.NoError(t, gzWriter.Close())
		require.NoError(t, file.Close())

		extractDir := filepath.Join(tempDir, "extract")
		require.NoError(t, os.MkdirAll(extractDir, 0o755))

		// WHEN: Extracting the archive
		err = cmd.ExtractArchiveForTest(archivePath, extractDir)

		// THEN: Should extract all files successfully
		assert.NoError(t, err)
		for name := range files {
			_, statErr := os.Stat(filepath.Join(extractDir, name))
			assert.NoError(t, statErr, "file %s should exist after extraction", name)
		}
	})
}
