package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
)

const (
	selfUpdateTimeout = 30 * time.Second
	terraRepoOwner    = "rios0rios0"
	terraRepoName     = "terra"
	githubAPIBaseURL  = "https://api.github.com"
	windowsOS         = "windows"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type SelfUpdateCommand struct{}

func NewSelfUpdateCommand() *SelfUpdateCommand {
	return &SelfUpdateCommand{}
}

func (it *SelfUpdateCommand) Execute(dryRun, force bool) error {
	logger.Info("Checking for terra updates...")

	// Get current version
	currentVersion := TerraVersion
	logger.Infof("Current terra version: %s", currentVersion)

	// Fetch latest release from GitHub
	latestVersion, downloadURL, err := it.fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}

	logger.Infof("Latest terra version: %s", latestVersion)

	// Compare versions
	comparison := compareVersions(currentVersion, latestVersion)
	switch {
	case comparison < 0:
		// Current version is older than latest
		if dryRun {
			logger.Infof("Dry run: Would update terra from %s to %s", currentVersion, latestVersion)
			logger.Infof("Download URL: %s", downloadURL)
			return nil
		}

		if !force && !it.promptForUpdate(currentVersion, latestVersion) {
			logger.Info("Update cancelled by user")
			return nil
		}

		logger.Infof("Updating terra from %s to %s...", currentVersion, latestVersion)
		return it.performUpdate(downloadURL)

	case comparison == 0:
		logger.Info("terra is already up to date")
		return nil

	default:
		logger.Infof("Current terra version %s is newer than latest available %s", currentVersion, latestVersion)
		return nil
	}
}

func (it *SelfUpdateCommand) fetchLatestRelease() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), selfUpdateTimeout)
	defer cancel()

	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", githubAPIBaseURL, terraRepoOwner, terraRepoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("error fetching release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %w", err)
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", "", fmt.Errorf("error parsing release JSON: %w", err)
	}

	// Extract version from tag name (remove 'v' prefix if present)
	version := strings.TrimPrefix(release.TagName, "v")

	// Find the appropriate archive for current platform.
	// GoReleaser naming: terra-<version>-<os>-<arch>.tar.gz (.zip on Windows)
	platform := entities.GetPlatformInfo()
	ext := "tar.gz"
	if platform.GetOSString() == windowsOS {
		ext = "zip"
	}
	expectedAssetName := fmt.Sprintf(
		"terra-%s-%s-%s.%s", version, platform.GetOSString(), platform.GetTerraformArchString(), ext,
	)

	for _, asset := range release.Assets {
		if asset.Name == expectedAssetName {
			return version, asset.BrowserDownloadURL, nil
		}
	}

	return "", "", fmt.Errorf("no asset %q found for platform %s", expectedAssetName, platform.GetPlatformString())
}

func (it *SelfUpdateCommand) promptForUpdate(currentVersion, latestVersion string) bool {
	logger.Infof("terra version %s is available (current: %s)", latestVersion, currentVersion)
	logger.Info("Do you want to update? [y/N]: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logger.Errorf("Error reading input: %v", err)
		}
		return false
	}
	response := strings.TrimSpace(strings.ToLower(scanner.Text()))

	return response == "y" || response == "yes"
}

func (it *SelfUpdateCommand) performUpdate(downloadURL string) error {
	currentOS := entities.GetOS()

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	currentExe, err = filepath.EvalSymlinks(currentExe)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Create temporary workspace for download and extraction
	tempDir, err := os.MkdirTemp("", "terra-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			logger.Warnf("Failed to cleanup temp directory %s: %v", tempDir, removeErr)
		}
	}()

	tempArchive := filepath.Join(tempDir, "terra-archive")

	logger.Info("Downloading new version...")
	err = currentOS.Download(downloadURL, tempArchive)
	if err != nil {
		return fmt.Errorf("failed to download new version: %w", err)
	}

	// Extract the archive to get the binary
	logger.Info("Extracting archive...")
	err = it.extractArchive(tempArchive, tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Locate the extracted binary
	binaryName := "terra"
	if entities.GetPlatformInfo().GetOSString() == windowsOS {
		binaryName = "terra.exe"
	}
	extractedBinary := filepath.Join(tempDir, binaryName)
	if _, statErr := os.Stat(extractedBinary); os.IsNotExist(statErr) {
		return fmt.Errorf("binary %q not found in extracted archive", binaryName)
	}

	// Make the extracted binary executable
	err = currentOS.MakeExecutable(extractedBinary)
	if err != nil {
		return fmt.Errorf("failed to make downloaded file executable: %w", err)
	}

	// Create backup of current binary
	backupFile := currentExe + ".backup"
	err = currentOS.Move(currentExe, backupFile)
	if err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary to current location
	err = currentOS.Move(extractedBinary, currentExe)
	if err != nil {
		// Try to restore backup on error
		if restoreErr := currentOS.Move(backupFile, currentExe); restoreErr != nil {
			logger.Errorf("Failed to restore backup: %v", restoreErr)
		}
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Remove backup file
	err = currentOS.Remove(backupFile)
	if err != nil {
		logger.Warnf("Failed to remove backup file %s: %v", backupFile, err)
	}

	logger.Info("terra has been successfully updated!")
	logger.Info("Please restart your terminal or run 'terra version' to verify the update")

	return nil
}

func (it *SelfUpdateCommand) extractArchive(archivePath, destDir string) error {
	platform := entities.GetPlatformInfo()

	// On Windows, use the OS abstraction (PowerShell Expand-Archive) for .zip files
	if platform.GetOSString() == windowsOS {
		currentOS := entities.GetOS()
		return currentOS.Extract(archivePath, destDir)
	}

	// On Unix, extract .tar.gz using tar
	ctx, cancel := context.WithTimeout(context.Background(), selfUpdateTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "tar", "-xzf", archivePath, "-C", destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tar.gz archive: %w", err)
	}
	return nil
}
