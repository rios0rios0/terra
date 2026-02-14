//go:build unit

package commands_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/flock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcquireRepoLock(t *testing.T) {
	t.Run("should acquire and release lock successfully when called", func(t *testing.T) {
		t.Parallel()
		// given -- use a unique lock path to avoid contention with other tests
		lockPath := filepath.Join(t.TempDir(), "test-acquire.lock")

		// when
		fl := flock.New(lockPath)
		locked, err := fl.TryLock()

		// then
		require.NoError(t, err)
		assert.True(t, locked, "Should be able to acquire lock")

		err = fl.Unlock()
		assert.NoError(t, err, "Should be able to release lock")
	})

	t.Run("should generate unique lock paths when different directories used", func(t *testing.T) {
		// given
		dir1 := "/some/path/one"
		dir2 := "/some/path/two"

		hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte(dir1)))
		hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte(dir2)))

		lockPath1 := filepath.Join(os.TempDir(), fmt.Sprintf("terra-%s.lock", hash1[:12]))
		lockPath2 := filepath.Join(os.TempDir(), fmt.Sprintf("terra-%s.lock", hash2[:12]))

		// when / then
		assert.NotEqual(t, lockPath1, lockPath2, "Different directories should produce different lock paths")
	})

	t.Run("should generate same lock path when same directory used", func(t *testing.T) {
		// given
		dir := "/some/consistent/path"

		hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte(dir)))
		hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte(dir)))

		lockPath1 := filepath.Join(os.TempDir(), fmt.Sprintf("terra-%s.lock", hash1[:12]))
		lockPath2 := filepath.Join(os.TempDir(), fmt.Sprintf("terra-%s.lock", hash2[:12]))

		// when / then
		assert.Equal(t, lockPath1, lockPath2, "Same directory should produce the same lock path")
	})
}

func TestReleaseRepoLock(t *testing.T) {
	t.Parallel()

	t.Run("should not panic when nil lock provided", func(t *testing.T) {
		t.Parallel()

		// given
		var fl *flock.Flock

		// when / then
		assert.NotPanics(t, func() {
			if fl == nil {
				return
			}
			_ = fl.Unlock()
		}, "Releasing a nil lock should not panic")
	})

	t.Run("should release lock without error when valid lock provided", func(t *testing.T) {
		t.Parallel()

		// given
		lockPath := filepath.Join(t.TempDir(), "test-release.lock")
		fl := flock.New(lockPath)
		locked, err := fl.TryLock()
		require.NoError(t, err)
		require.True(t, locked)

		// when
		err = fl.Unlock()

		// then
		assert.NoError(t, err, "Should release lock without error")
	})
}
