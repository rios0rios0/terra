package entities

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	osOrwxGrxUx      = 0o755
	operationTimeout = 30 * time.Second
)

type OSLinux struct{}

func (it *OSLinux) Download(url, tempFilePath string) error {
	return downloadFile(url, tempFilePath)
}

func (it *OSLinux) Extract(tempFilePath, destPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()
	unzipCmd := exec.CommandContext(ctx, "unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	err := unzipCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform decompressing using 'zip': %w", err)
	}
	return err
}

func (it *OSLinux) Move(tempFilePath, destPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()
	mvCmd := exec.CommandContext(ctx, "mv", tempFilePath, destPath)
	err := mvCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform moving folder using 'mv': %w", err)
	}
	return err
}

func (it *OSLinux) Remove(tempFilePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()
	rmCmd := exec.CommandContext(ctx, "rm", tempFilePath)
	err := rmCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform deleting folder using 'rm': %w", err)
	}
	return err
}

func (it *OSLinux) MakeExecutable(filePath string) error {
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	err := os.Chmod(filePath, osOrwxGrxUx)
	if err != nil {
		err = fmt.Errorf("failed to perform change binary permissions using 'chmod': %w", err)
	}
	return err
}

func (it *OSLinux) GetTempDir() string {
	return os.TempDir()
}

func (it *OSLinux) GetInstallationPath() string {
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

func GetOS() *OSLinux {
	return &OSLinux{}
}
