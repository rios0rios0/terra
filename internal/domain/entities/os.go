package entities

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	downloadTimeout = 10 * time.Minute
)

type OS interface {
	Download(url, tempFilePath string) error
	Extract(tempFilePath, destPath string) error
	Move(tempFilePath, destPath string) error
	Remove(tempFilePath string) error
	MakeExecutable(filePath string) error
	GetTempDir() string
	GetInstallationPath() string
}

// downloadFile provides a common implementation for downloading files via HTTP.
func downloadFile(url, tempFilePath string) error {
	// Create context with timeout for the download
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
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
