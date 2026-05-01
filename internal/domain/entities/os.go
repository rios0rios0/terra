package entities

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	logger "github.com/sirupsen/logrus"
)

const (
	// defaultDownloadTimeout caps a single download HTTP request +
	// body read for `terra install`. Configurable via the
	// `TERRA_DOWNLOAD_TIMEOUT` environment variable when slower
	// transports (corporate proxies, QEMU-emulated multi-arch
	// container builds, low-bandwidth links) need a longer ceiling.
	// Accepts any `time.ParseDuration` value, e.g. `30m`, `1h`,
	// `20m30s`.
	defaultDownloadTimeout = 10 * time.Minute

	// downloadTimeoutEnvVar is the environment-variable name that
	// overrides `defaultDownloadTimeout`.
	downloadTimeoutEnvVar = "TERRA_DOWNLOAD_TIMEOUT"
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

// resolveDownloadTimeout returns the user-configured download
// timeout from `TERRA_DOWNLOAD_TIMEOUT`, falling back to
// `defaultDownloadTimeout` when the env var is unset, empty, or
// fails to parse. A parse failure logs a warning and uses the
// default so a malformed override never silently breaks installs.
func resolveDownloadTimeout() time.Duration {
	raw := os.Getenv(downloadTimeoutEnvVar)
	if raw == "" {
		return defaultDownloadTimeout
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		logger.Warnf(
			"invalid %s=%q (%s); falling back to default %s",
			downloadTimeoutEnvVar, raw, err, defaultDownloadTimeout,
		)
		return defaultDownloadTimeout
	}

	return parsed
}

// downloadFile provides a common implementation for downloading files via HTTP.
func downloadFile(url, tempFilePath string) error {
	// Create context with timeout for the download. Resolved at
	// call time (not init) so tests / operators can override
	// `TERRA_DOWNLOAD_TIMEOUT` without a process restart.
	ctx, cancel := context.WithTimeout(context.Background(), resolveDownloadTimeout())
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
