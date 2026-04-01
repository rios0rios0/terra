//go:build unit

package commands_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

