//go:build integration || unit || test

package repositoryhelpers //nolint:staticcheck // Test package naming follows established project structure

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// HelperCreateNonZipBinaryServer creates a server that serves a regular binary (not zip).
func HelperCreateNonZipBinaryServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Serve a simple binary content.
		binaryContent := []byte("#!/bin/bash\necho 'test binary'\n")
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(binaryContent)
	}))
}

// HelperCreateSimpleVersionServer creates a version server with given version.
func HelperCreateSimpleVersionServer(t *testing.T, version string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := `{"current_version":"` + version + `"}`
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
}

// HelperCreateSimpleZipServer creates a server that serves a zip file.
func HelperCreateSimpleZipServer(t *testing.T, binaryName string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Create a simple zip file in memory (reuse from existing zip helper).
		zipData := HelperCreateZipWithBinary(t, binaryName)

		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(zipData)
	}))
}

// HelperCreateZipServer creates a test server that serves a zip file containing a binary.
func HelperCreateZipServer(
	t *testing.T,
	binaryNameInZip, _ string,
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Create a zip file in memory.
		zipData := HelperCreateZipWithBinary(t, binaryNameInZip)

		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(zipData)
	}))
}

// HelperCreateVersionServer creates a test server that serves version information.
func HelperCreateVersionServer(t *testing.T, version string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := `{"current_version":"` + version + `"}`
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
}

// HelperCreateNestedZipServer creates a zip with binary in nested directory structure.
func HelperCreateNestedZipServer(t *testing.T, nestedPath, binaryName string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Create a zip file with nested directory structure.
		zipData := HelperCreateNestedZipWithBinary(t, nestedPath, binaryName)

		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(zipData)
	}))
}

// HelperCreateMixedContentZipServer creates a zip with various file types.
func HelperCreateMixedContentZipServer(t *testing.T, binaryName string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Create a zip file with mixed content.
		zipData := HelperCreateMixedContentZip(t, binaryName)

		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(zipData)
	}))
}
