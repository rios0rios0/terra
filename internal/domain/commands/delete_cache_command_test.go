//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	testhelpers "github.com/rios0rios0/terra/test/infrastructure/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteCacheCommand_Execute(t *testing.T) {
	t.Run("should remove directories when matching the given names", func(t *testing.T) {
		// given
		helper := testhelpers.DirectoryHelper{}
		toBeDeleted := helper.CreateTestDirectories()

		command := commands.NewDeleteCacheCommand()

		listeners := interfaces.DeleteCacheListeners{OnSuccess: func() {
			// then
			for _, dir := range toBeDeleted {
				assert.False(t, helper.DirectoryExists(dir), "directory should be removed: "+dir)
			}
		}}

		// when
		command.Execute(toBeDeleted, listeners)
	})

	t.Run("should handle non-existent directories gracefully", func(t *testing.T) {
		// given
		helper := testhelpers.DirectoryHelper{}
		toBeDeleted := []string{"non-existent"}
		command := commands.NewDeleteCacheCommand()

		listeners := interfaces.DeleteCacheListeners{OnSuccess: func() {
			// then
			assert.False(t, helper.DirectoryExists("non-existent"), "non-existent directory should not cause errors")
		}}

		// when
		command.Execute(toBeDeleted, listeners)
	})

	t.Run("should throw an error when some unexpected condition happens", func(t *testing.T) {
		// given
		toBeDeleted := []string{"."}
		command := commands.NewDeleteCacheCommand()

		listeners := interfaces.DeleteCacheListeners{OnError: func(err error) {
			// then
			assert.Error(t, err, "should return an error when trying to remove and it couldn't")
		}}

		// when
		command.Execute(toBeDeleted, listeners)
	})
}
