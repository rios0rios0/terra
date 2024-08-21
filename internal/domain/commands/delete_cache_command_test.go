//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteCacheCommand_Execute(t *testing.T) {
	t.Run("should remove directories matching the given names", func(t *testing.T) {
		// given
		toBeDeleted := createTestDirectories()
		command := commands.NewDeleteCacheCommand()

		// when
		command.Execute(toBeDeleted)

		// then
		for _, dir := range toBeDeleted {
			assert.False(t, directoryExists(dir), "directory should be removed: "+dir)
		}
	})

	t.Run("should handle non-existent directories gracefully", func(t *testing.T) {
		// given
		toBeDeleted := []string{"nonexistent"}
		command := commands.NewDeleteCacheCommand()

		// when
		command.Execute(toBeDeleted)

		// then
		assert.False(t, directoryExists("nonexistent"), "non-existent directory should not cause errors")
	})
}

func createTestDirectories() []string {
	dirs := []string{"cache", "temp"}
	for _, dir := range dirs {
		_ = os.MkdirAll(filepath.Join(".", dir), os.ModePerm)
	}
	return dirs
}

func directoryExists(dir string) bool {
	info, err := os.Stat(filepath.Join(".", dir))
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
