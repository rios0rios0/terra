package commands_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/infrastructure/repositorybuilders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// scriptedVersionServer answers each request with the next status, repeating the
// last one, and counts the requests it receives.
func scriptedVersionServer(t *testing.T, body string, statuses ...int) (*httptest.Server, *atomic.Int32) {
	t.Helper()

	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		attempt := int(requests.Add(1)) - 1
		writer.WriteHeader(statuses[min(attempt, len(statuses)-1)])
		_, _ = writer.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return server, &requests
}

func TestNewInstallDependenciesCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: The NewInstallDependenciesCommand constructor is available

		// WHEN: Creating a new install dependencies command
		cmd := commands.NewInstallDependenciesCommand()

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should complete without error when empty dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An install dependencies command and empty dependencies list
		cmd := commands.NewInstallDependenciesCommand()
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command with empty dependencies
		// Note: This should complete quickly without attempting downloads
		cmd.Execute(dependencies)

		// THEN: Should complete without panicking or errors (verified by not crashing)
	})

	t.Run("should install dependency when dependency not available", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock server and dependency for non-existent CLI
		versionServer, binaryServer := repositorybuilders.NewTestServerBuilder().
			WithTerraformVersion("1.0.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entitybuilders.NewDependencyBuilder().
			WithName("TestTool").
			WithCLI("non-existent-cli-tool-12345").
			WithBinaryURL(binaryServer.URL + "/testtool_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			BuildDependency()

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should attempt to install (tested through integration test success)
		// This tests the isDependencyCLIAvailable (100% covered) and install (partial coverage) methods
	})

	t.Run("should skip update check when current version cannot be determined", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A mock server and dependency for terraform (which should be available)
		versionServer, binaryServer := repositorybuilders.NewTestServerBuilder().
			WithTerraformVersion("1.0.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependency := entitybuilders.NewDependencyBuilder().
			WithName("Terraform").
			WithCLI("terraform").
			WithBinaryURL(binaryServer.URL + "/terraform_%s").
			WithVersionURL(versionServer.URL + "/terraform").
			WithTerraformPattern().
			BuildDependency()

		// WHEN: Executing the command with terraform (if installed)
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute([]entities.Dependency{dependency})

		// THEN: Should handle version determination logic
		// This tests getCurrentVersion method indirectly when terraform is available
	})

	t.Run("should handle version comparison scenarios", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Multiple dependencies to test different version scenarios
		versionServer, binaryServer := repositorybuilders.NewTestServerBuilder().
			WithTerraformVersion("1.0.0").
			WithTerragruntVersion("0.50.0").
			BuildServers()
		defer versionServer.Close()
		defer binaryServer.Close()

		dependencies := []entities.Dependency{
			entitybuilders.NewDependencyBuilder().
				WithName("Terraform").
				WithCLI("terraform").
				WithBinaryURL(binaryServer.URL + "/terraform_%s").
				WithVersionURL(versionServer.URL + "/terraform").
				WithTerraformPattern().
				BuildDependency(),
			entitybuilders.NewDependencyBuilder().
				WithName("Terragrunt").
				WithCLI("terragrunt").
				WithBinaryURL(binaryServer.URL + "/terragrunt_%s").
				WithVersionURL(versionServer.URL + "/terragrunt").
				WithTerragruntPattern().
				BuildDependency(),
		}

		// WHEN: Executing the command
		cmd := commands.NewInstallDependenciesCommand()
		cmd.Execute(dependencies)

		// THEN: Should handle version comparison logic
		// This tests compareVersions method indirectly through version checking logic
	})
}

func TestFetchLatestVersion(t *testing.T) {
	t.Parallel()

	const (
		pattern     = `"tag_name":"v([^"]+)"`
		versionBody = `{"tag_name":"v1.2.1"}`
		limitBody   = `{"message":"API rate limit exceeded"}`
	)

	t.Run("should return the version when the server answers 200", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server that answers with a release
		server, requests := scriptedVersionServer(t, versionBody, http.StatusOK)

		// WHEN: Looking up the latest version
		version, err := commands.FetchLatestVersionPublic(server.URL, pattern)

		// THEN: Should return it from the first request
		require.NoError(t, err)
		assert.Equal(t, "1.2.1", version)
		assert.Equal(t, int32(1), requests.Load())
	})

	t.Run("should report the status and body when the server keeps answering 403", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server that rate limits every request
		server, requests := scriptedVersionServer(t, limitBody, http.StatusForbidden)

		// WHEN: Looking up the latest version
		_, err := commands.FetchLatestVersionPublic(server.URL, pattern)

		// THEN: Should say what the server said, after trying three times
		require.ErrorContains(t, err, "HTTP 403")
		require.ErrorContains(t, err, "API rate limit exceeded")
		assert.NotContains(t, err.Error(), "regex")
		assert.Equal(t, int32(3), requests.Load())
	})

	t.Run("should succeed when a 503 passes before the attempts run out", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server that fails once and then recovers
		server, requests := scriptedVersionServer(t, versionBody, http.StatusServiceUnavailable, http.StatusOK)

		// WHEN: Looking up the latest version
		version, err := commands.FetchLatestVersionPublic(server.URL, pattern)

		// THEN: Should return the version from the second request
		require.NoError(t, err)
		assert.Equal(t, "1.2.1", version)
		assert.Equal(t, int32(2), requests.Load())
	})

	t.Run("should not retry when the server answers 404", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server that will never have the resource
		server, requests := scriptedVersionServer(t, "not found", http.StatusNotFound)

		// WHEN: Looking up the latest version
		_, err := commands.FetchLatestVersionPublic(server.URL, pattern)

		// THEN: Should stop at the first request
		require.ErrorContains(t, err, "HTTP 404")
		assert.Equal(t, int32(1), requests.Load())
	})

	t.Run("should name the regex when a 200 holds no version", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A server whose successful answer does not match the pattern
		server, _ := scriptedVersionServer(t, `{"tag_name":"1.2.1"}`, http.StatusOK)

		// WHEN: Looking up the latest version
		_, err := commands.FetchLatestVersionPublic(server.URL, pattern)

		// THEN: Should report the mismatch, which is what really happened here
		require.ErrorContains(t, err, "no version match found")
		require.ErrorContains(t, err, pattern)
	})

	t.Run("should keep an installed tool when the lookup keeps failing", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A tool that is installed, and a server that rate limits every request
		server, requests := scriptedVersionServer(t, limitBody, http.StatusForbidden)
		dependency := entitybuilders.NewDependencyBuilder().
			WithName("Shell").
			WithCLI("sh").
			WithBinaryURL(server.URL + "/sh_%s").
			WithVersionURL(server.URL).
			WithTerraformPattern().
			BuildDependency()

		// WHEN: Executing the command
		commands.NewInstallDependenciesCommand().Execute([]entities.Dependency{dependency})

		// THEN: Should return instead of exiting, without trying to install anything
		assert.Equal(t, int32(3), requests.Load())
	})
}

func TestAuthorizeVersionRequest(t *testing.T) {
	const token = "fixture-token-placeholder"

	newRequest := func(t *testing.T, url string) *http.Request {
		t.Helper()
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
		require.NoError(t, err)

		return request
	}

	t.Run("should send the token to the GitHub API when it is set", func(t *testing.T) {
		// GIVEN: A token and a request for the GitHub API
		t.Setenv("GITHUB_TOKEN", token)
		request := newRequest(t, "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest")

		// WHEN: Authorizing the request
		commands.AuthorizeVersionRequestPublic(request)

		// THEN: Should carry the token
		assert.Equal(t, "Bearer "+token, request.Header.Get("Authorization"))
	})

	t.Run("should not send the token to any other host", func(t *testing.T) {
		// GIVEN: A token and a request for another host
		t.Setenv("GITHUB_TOKEN", token)
		request := newRequest(t, "https://checkpoint-api.hashicorp.com/v1/check/terraform")

		// WHEN: Authorizing the request
		commands.AuthorizeVersionRequestPublic(request)

		// THEN: Should stay anonymous
		assert.Empty(t, request.Header.Get("Authorization"))
	})

	t.Run("should stay anonymous when no token is set", func(t *testing.T) {
		// GIVEN: No token and a request for the GitHub API
		t.Setenv("GITHUB_TOKEN", "")
		request := newRequest(t, "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest")

		// WHEN: Authorizing the request
		commands.AuthorizeVersionRequestPublic(request)

		// THEN: Should send no header
		assert.Empty(t, request.Header.Get("Authorization"))
	})
}
