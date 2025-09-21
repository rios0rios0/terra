package entities

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// testDownloadSuccess is a helper function to test successful download for any OS implementation
func testDownloadSuccess(t *testing.T, osImpl OS, testPrefix string) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("test file content"))
	}))
	defer server.Close()

	// Create a secure temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), testPrefix+"_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close() // Close the file so Download can create it
	defer os.Remove(tempFile.Name())

	// Test the download
	err = osImpl.Download(server.URL, tempFile.Name())
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Verify the file was created and has the correct content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	expectedContent := "test file content"
	if string(content) != expectedContent {
		t.Errorf(
			"Downloaded content doesn't match. Expected: '%s', Got: '%s'",
			expectedContent,
			string(content),
		)
	}
}

// testDownloadHTTPError is a helper function to test HTTP error handling for any OS implementation
func testDownloadHTTPError(t *testing.T, osImpl OS, testPrefix string) {
	t.Helper()
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create a secure temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), testPrefix+"_error_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close() // Close the file so Download can create it
	defer os.Remove(tempFile.Name())

	// Test the download - should fail
	err = osImpl.Download(server.URL, tempFile.Name())
	if err == nil {
		t.Error("Expected download to fail with HTTP 500, but it succeeded")
	}

	// Verify the error message contains HTTP status information
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("Expected error to contain 'HTTP 500', but got: %s", err.Error())
	}
}
