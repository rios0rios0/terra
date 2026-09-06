//go:build !windows

package entities

import (
	"fmt"
	"os"
)

const osOrwxGrxUx = 0o755

type OSUnix struct{}

func (it *OSUnix) Download(url, tempFilePath string) error {
	return downloadFile(url, tempFilePath)
}

func (it *OSUnix) Extract(tempFilePath, destPath string) error {
	return extractZipArchive(tempFilePath, destPath)
}

func (it *OSUnix) Move(tempFilePath, destPath string) error {
	return moveFile(tempFilePath, destPath)
}

func (it *OSUnix) Remove(tempFilePath string) error {
	return removeFile(tempFilePath)
}

func (it *OSUnix) MakeExecutable(filePath string) error {
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	err := os.Chmod(filePath, osOrwxGrxUx)
	if err != nil {
		err = fmt.Errorf("failed to perform change binary permissions using 'chmod': %w", err)
	}
	return err
}

func (it *OSUnix) GetTempDir() string {
	return os.TempDir()
}

func (it *OSUnix) GetInstallationPath() string {
	// Allow override via environment variable (used by tests to avoid
	// overwriting real binaries like terraform in ~/.local/bin).
	if envPath := os.Getenv("TERRA_INSTALL_PATH"); envPath != "" {
		return envPath
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.local/bin" // Fallback to original path
	}
	return fmt.Sprintf("%s/.local/bin", homeDir)
}

func GetOS() *OSUnix {
	return &OSUnix{}
}
