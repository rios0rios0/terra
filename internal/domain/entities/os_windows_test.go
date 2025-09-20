package entities

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestOSWindows_Download(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test file content"))
	}))
	defer server.Close()

	// Create a secure temporary file
	tempFile, err := os.CreateTemp("", "test_download_windows_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close() // Close the file so Download can create it
	defer os.Remove(tempFile.Name())

	// Test the download
	osWindows := &OSWindows{}
	err = osWindows.Download(server.URL, tempFile.Name())
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
		t.Errorf("Downloaded content doesn't match. Expected: '%s', Got: '%s'", expectedContent, string(content))
	}
}

func TestOSWindows_DownloadHTTPError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Create a secure temporary file
	tempFile, err := os.CreateTemp("", "test_download_error_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close() // Close the file so Download can create it
	defer os.Remove(tempFile.Name())

	// Test the download - should fail
	osWindows := &OSWindows{}
	err = osWindows.Download(server.URL, tempFile.Name())
	if err == nil {
		t.Error("Expected download to fail with HTTP 500, but it succeeded")
	}

	// Verify the error message contains HTTP status information
	if !containsString(err.Error(), "HTTP 500") {
		t.Errorf("Expected error to contain 'HTTP 500', but got: %s", err.Error())
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}())
}
