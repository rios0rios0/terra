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

type OSWindows struct{}

func (it *OSWindows) Download(url, tempFilePath string) error {
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

func (it *OSWindows) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("powershell", "Expand-Archive", "-Path", tempFilePath, "-DestinationPath", destPath, "-Force")
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	err := unzipCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform decompressing using 'powershell': %w", err)
	}
	return err
}

func (it *OSWindows) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("move", tempFilePath, destPath)
	err := mvCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform moving folder using 'move': %w", err)
	}
	return err
}

func (it *OSWindows) Remove(tempFilePath string) error {
	rmCmd := exec.Command("del", tempFilePath)
	err := rmCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform deleting folder using 'del': %w", err)
	}
	return err
}

func (it *OSWindows) MakeExecutable(_ string) error {
	return nil // Windows doesn't need to explicitly make files executable
}

func (it *OSWindows) GetTempDir() string {
	return os.Getenv("TEMP")
}

func (it *OSWindows) GetInstallationPath() string {
	return os.Getenv("ProgramFiles")
}

func GetOS() *OSWindows {
	return &OSWindows{}
}
