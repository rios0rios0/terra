package repository_builders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// GitHubAPIServerBuilder helps create mock GitHub API servers for testing.
type GitHubAPIServerBuilder struct {
	t              *testing.T
	releaseVersion string
	assetName      string
	assetURL       string
	statusCode     int
	errorResponse  string
}

// NewGitHubAPIServerBuilder creates a new builder for GitHub API test servers.
func NewGitHubAPIServerBuilder(t *testing.T) *GitHubAPIServerBuilder {
	t.Helper()
	return &GitHubAPIServerBuilder{
		t:          t,
		statusCode: http.StatusOK,
	}
}

// WithRelease sets up a mock release response.
func (b *GitHubAPIServerBuilder) WithRelease(version, assetName, assetURL string) *GitHubAPIServerBuilder {
	b.releaseVersion = version
	b.assetName = assetName
	b.assetURL = assetURL
	return b
}

// WithStatusCode sets the HTTP status code to return.
func (b *GitHubAPIServerBuilder) WithStatusCode(code int) *GitHubAPIServerBuilder {
	b.statusCode = code
	return b
}

// WithErrorResponse sets an error response body.
func (b *GitHubAPIServerBuilder) WithErrorResponse(errorMsg string) *GitHubAPIServerBuilder {
	b.errorResponse = errorMsg
	return b
}

// BuildServer creates and returns a test server.
func (b *GitHubAPIServerBuilder) BuildServer() *httptest.Server {
	b.t.Helper()
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(b.statusCode)
		
		if b.errorResponse != "" {
			fmt.Fprint(w, b.errorResponse)
			return
		}
		
		if b.statusCode != http.StatusOK {
			fmt.Fprint(w, `{"message": "API Error"}`)
			return
		}
		
		response := map[string]interface{}{
			"tag_name": b.releaseVersion,
			"assets": []map[string]interface{}{
				{
					"name":                 b.assetName,
					"browser_download_url": b.assetURL,
				},
			},
		}
		
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			b.t.Fatalf("Failed to encode JSON response: %v", err)
		}
	}))
	
	return server
}

// BuildBinaryServer creates a test server that serves binary files.
func (b *GitHubAPIServerBuilder) BuildBinaryServer() *httptest.Server {
	b.t.Helper()
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		// Simulate a binary file with some dummy content
		fmt.Fprint(w, "fake binary content for testing")
	}))
	
	return server
}