//nolint:staticcheck // Test package naming follows established project structure
package repository_builders

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// TestServerBuilder helps create mock servers with a fluent interface.
type TestServerBuilder struct {
	versionResponses map[string]string
	binaryResponse   []byte
	binaryStatus     int
	contentType      string
	shouldFail       bool
}

// NewTestServerBuilder creates a new test server builder.
func NewTestServerBuilder() *TestServerBuilder {
	return &TestServerBuilder{
		versionResponses: make(map[string]string),
		binaryResponse:   []byte("#!/bin/bash\necho 'mock binary'\n"),
		binaryStatus:     http.StatusOK,
		contentType:      "application/octet-stream",
		shouldFail:       false,
	}
}

// WithVersionResponse adds a version response for a specific path pattern.
func (b *TestServerBuilder) WithVersionResponse(pathPattern, response string) *TestServerBuilder {
	b.versionResponses[pathPattern] = response
	return b
}

// WithTerraformVersion adds a terraform version response.
func (b *TestServerBuilder) WithTerraformVersion(version string) *TestServerBuilder {
	response := `{"current_version":"` + version + `"}`
	return b.WithVersionResponse("terraform", response)
}

// WithTerragruntVersion adds a terragrunt version response.
func (b *TestServerBuilder) WithTerragruntVersion(version string) *TestServerBuilder {
	response := `{"tag_name":"v` + version + `"}`
	return b.WithVersionResponse("terragrunt", response)
}

// WithBinaryResponse sets the binary response content.
func (b *TestServerBuilder) WithBinaryResponse(content []byte) *TestServerBuilder {
	b.binaryResponse = content
	return b
}

// WithBinaryStatus sets the binary server response status.
func (b *TestServerBuilder) WithBinaryStatus(status int) *TestServerBuilder {
	b.binaryStatus = status
	return b
}

// WithZipContent sets up the server to return zip content.
func (b *TestServerBuilder) WithZipContent() *TestServerBuilder {
	b.binaryResponse = []byte("PK\x03\x04test")
	b.contentType = "application/zip"
	return b
}

// WithDownloadFailure sets up the server to simulate download failures.
func (b *TestServerBuilder) WithDownloadFailure() *TestServerBuilder {
	b.shouldFail = true
	b.binaryStatus = http.StatusInternalServerError
	return b
}

// BuildServers creates and returns the version and binary servers.
func (b *TestServerBuilder) BuildServers() (*httptest.Server, *httptest.Server) {
	versionServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			// Check for specific patterns first (longer matches first).
			for pattern, response := range b.versionResponses {
				if pattern != "" && strings.Contains(r.URL.Path, pattern) {
					_, _ = w.Write([]byte(response))
					return
				}
			}

			// Use default response if available.
			if defaultResponse, exists := b.versionResponses[""]; exists {
				_, _ = w.Write([]byte(defaultResponse))
				return
			}

			// Final fallback.
			_, _ = w.Write([]byte(`{"version":"1.0.0"}`))
		}),
	)

	binaryServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			if b.shouldFail {
				w.WriteHeader(b.binaryStatus)
				_, _ = w.Write([]byte("download failed"))
				return
			}
			w.Header().Set("Content-Type", b.contentType)
			w.WriteHeader(b.binaryStatus)
			_, _ = w.Write(b.binaryResponse)
		}),
	)

	return versionServer, binaryServer
}
