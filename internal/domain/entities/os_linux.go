package entities

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const osOrwxGrxUx = 0o755

type OSLinux struct{}

func (it *OSLinux) Download(url, tempFilePath string) error {
	// Create context with timeout for the download
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform download: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	// Create the destination file
	out, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", tempFilePath, err)
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write downloaded content to file: %w", err)
	}

	return nil
}

func (it *OSLinux) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	err := unzipCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform decompressing using 'zip': %w", err)
	}
	return err
}

func (it *OSLinux) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("mv", tempFilePath, destPath)
	err := mvCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform moving folder using 'mv': %w", err)
	}
	return err
}

func (it *OSLinux) Remove(tempFilePath string) error {
	rmCmd := exec.Command("rm", tempFilePath)
	err := rmCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform deleting folder using 'rm': %w", err)
	}
	return err
}

func (it *OSLinux) MakeExecutable(filePath string) error {
	err := os.Chmod(filePath, osOrwxGrxUx)
	if err != nil {
		err = fmt.Errorf("failed to perform change binary permissions using 'chmod': %w", err)
	}
	return err
}

func (it *OSLinux) GetTempDir() string {
	return "/tmp"
}

func (it *OSLinux) GetInstallationPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.local/bin" // Fallback to original path
	}
	return fmt.Sprintf("%s/.local/bin", homeDir)
}

func GetOS() *OSLinux {
	return &OSLinux{}
}
