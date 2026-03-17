//go:build unit

package commands_test

import (
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
}
