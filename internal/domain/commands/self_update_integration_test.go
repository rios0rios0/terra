//go:build integration

package commands_test

import (
	"strings"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
)

// skipOnRateLimit skips the test if the error indicates a GitHub API rate limit (HTTP 403).
// This is a known flaky condition in CI environments without GitHub tokens.
func skipOnRateLimit(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "403") {
		t.Skip("Skipping test due to GitHub API rate limiting (HTTP 403)")
	}
}

func TestSelfUpdateCommand_Execute_Integration(t *testing.T) {
	t.Run("should perform dry run successfully when valid release available", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		// GIVEN
		cmd := commands.NewSelfUpdateCommand()

		// WHEN
		err := cmd.Execute(true, false)

		// THEN: Dry run should succeed (or fail only due to rate limiting)
		skipOnRateLimit(t, err)
		assert.NoError(t, err)
	})
}

func TestSelfUpdateCommand_RealGitHubAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real API test in short mode")
	}

	t.Run("should successfully connect to GitHub API", func(t *testing.T) {
		// GIVEN
		cmd := commands.NewSelfUpdateCommand()

		// WHEN
		err := cmd.Execute(true, false)

		// THEN: Should successfully reach GitHub API and find the correct asset
		skipOnRateLimit(t, err)
		assert.NoError(t, err)
	})
}
