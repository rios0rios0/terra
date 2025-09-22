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
		// GIVEN: No dependencies required for constructor

		// WHEN: Creating a new self-update command
		cmd := commands.NewSelfUpdateCommand()

		// THEN: Should create a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestSelfUpdateCommand_Execute(t *testing.T) {
	t.Run("should return error when dry run requested and GitHub API fails", func(t *testing.T) {
		// GIVEN: A self-update command and invalid GitHub API scenario
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with dry run (this will hit real GitHub API and likely fail due to rate limiting)
		err := cmd.Execute(true, false)

		// THEN: Should return an error due to API limitations in test environment
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
	})

	t.Run("should return error when force flag used and GitHub API fails", func(t *testing.T) {
		// GIVEN: A self-update command and invalid GitHub API scenario
		cmd := commands.NewSelfUpdateCommand()

		// WHEN: Executing with force flag (this will hit real GitHub API and likely fail due to rate limiting)
		err := cmd.Execute(false, true)

		// THEN: Should return an error due to API limitations in test environment
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch latest release")
	})
}
