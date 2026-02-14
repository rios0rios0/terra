package commands

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	logger "github.com/sirupsen/logrus"
)

// acquireRepoLock acquires an exclusive file lock for the current repository.
// This prevents concurrent terra processes from corrupting shared terragrunt caches
// (e.g. .terragrunt-cache directories in shared dependency paths).
// The lock is advisory and cross-platform (Unix and Windows).
func acquireRepoLock() (*flock.Flock, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(cwd)))
	lockPath := filepath.Join(os.TempDir(), fmt.Sprintf("terra-%s.lock", hash[:12]))

	fl := flock.New(lockPath)
	logger.Debugf("Acquiring repo lock on %s...", lockPath)

	if err = fl.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire lock on %s: %w", lockPath, err)
	}

	logger.Debugf("Repo lock acquired on %s", lockPath)

	return fl, nil
}

// releaseRepoLock releases the exclusive file lock.
func releaseRepoLock(fl *flock.Flock) {
	if fl == nil {
		return
	}

	_ = fl.Unlock()
	logger.Debugf("Repo lock released")
}
