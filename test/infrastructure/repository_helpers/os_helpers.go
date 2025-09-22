//go:build integration || unit || test

package repository_helpers //nolint:staticcheck // Test package naming follows established project structure

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HelperDownloadSuccess is a helper function to test successful download for any OS implementation.
func HelperDownloadSuccess(t *testing.T, osImpl entities.OS, testPrefix string) {
	// Create a test server.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("test file content"))
	}))
	defer server.Close()

	// Create a secure temporary file.
	tempFile, err := os.CreateTemp(t.TempDir(), testPrefix+"_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	if closeErr := tempFile.Close(); closeErr != nil { // Close the file so Download can create it.
		t.Fatalf("Failed to close temporary file: %v", closeErr)
	}
	defer os.Remove(tempFile.Name())

	// Test the download.
	err = osImpl.Download(server.URL, tempFile.Name())
	require.NoError(t, err, "Download should succeed")

	// Verify the file was created and has the correct content.
	content, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err, "Should be able to read downloaded file")

	expectedContent := "test file content"
	assert.Equal(t, expectedContent, string(content),
		"Downloaded content should match expected content")
}

// HelperDownloadHTTPError is a helper function to test HTTP error handling for any OS implementation.
func HelperDownloadHTTPError(t *testing.T, osImpl entities.OS, testPrefix string) {
	t.Helper()
	// Create a test server that returns an error.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create a secure temporary file.
	tempFile, err := os.CreateTemp(t.TempDir(), testPrefix+"_error_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	if closeErr := tempFile.Close(); closeErr != nil { // Close the file so Download can create it.
		t.Fatalf("Failed to close temporary file: %v", closeErr)
	}
	defer os.Remove(tempFile.Name())

	// Test the download - should fail.
	err = osImpl.Download(server.URL, tempFile.Name())
	require.Error(t, err, "Download should fail with HTTP 500")

	// Verify the error message contains HTTP status information.
	assert.Contains(t, err.Error(), "HTTP 500",
		"Error message should contain HTTP 500 status")
}
