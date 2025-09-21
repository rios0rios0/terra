package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
)

const contextTimeout = 10 * time.Second

type InstallDependenciesCommand struct{}

func NewInstallDependenciesCommand() *InstallDependenciesCommand {
	return &InstallDependenciesCommand{}
}

func (it *InstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	for _, dependency := range dependencies {
		latestVersion := fetchLatestVersion(dependency.VersionURL, dependency.RegexVersion)

		if !isDependencyCLIAvailable(dependency.CLI) {
			logger.Warnf("%s is not installed, installing now...", dependency.Name)
			install(dependency.GetBinaryURL(latestVersion), dependency.CLI)
		} else {
			// Dependency is installed, check if it's the latest version
			currentVersion := getCurrentVersion(dependency.CLI)
			if currentVersion == "" {
				logger.Warnf("Could not determine current version of %s, skipping update check", dependency.Name)
				continue
			}

			comparison := compareVersions(currentVersion, latestVersion)
			switch {
			case comparison < 0:
				// Current version is older than latest
				if promptForUpdate(dependency.Name, currentVersion, latestVersion) {
					logger.Infof("Updating %s from %s to %s...", dependency.Name, currentVersion, latestVersion)
					install(dependency.GetBinaryURL(latestVersion), dependency.CLI)
				} else {
					logger.Infof("Skipping update for %s", dependency.Name)
				}
			case comparison == 0:
				logger.Infof("%s is already up to date (version %s)", dependency.Name, currentVersion)
			default:
				logger.Infof("%s version %s is newer than latest available %s", dependency.Name, currentVersion, latestVersion)
			}
		}
	}
}

// fetch the latest version of software from a URL
func fetchLatestVersion(url, regexPattern string) string {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Fatalf("Error creating request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Fatalf("Error fetching version info: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("Error reading response body: %s", err)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1]
	}

	logger.Fatalf("No version match found, check the regex pattern: %s", regexPattern)
	return ""
}

// checking if a dependency is available
func isDependencyCLIAvailable(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, "-v")
	return cmd.Run() == nil
}

// get current version of installed dependency
func getCurrentVersion(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	var cmd *exec.Cmd
	switch name {
	case "terraform":
		cmd = exec.CommandContext(ctx, name, "--version")
	case "terragrunt":
		cmd = exec.CommandContext(ctx, name, "--version")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		logger.Debugf("Failed to get %s version: %s", name, err)
		return ""
	}

	version := strings.TrimSpace(string(output))

	// Extract version number from output
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// compare two semantic versions (returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2)
func compareVersions(v1, v2 string) int {
	// Only compare numeric parts, reject versions with non-numeric components
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Validate that all parts are numeric for proper semantic version comparison
	for _, part := range parts1 {
		if _, err := strconv.Atoi(part); err != nil {
			logger.Warnf(
				"Version %s contains non-numeric parts, cannot perform reliable comparison",
				v1,
			)
			return strings.Compare(v1, v2)
		}
	}
	for _, part := range parts2 {
		if _, err := strconv.Atoi(part); err != nil {
			logger.Warnf(
				"Version %s contains non-numeric parts, cannot perform reliable comparison",
				v2,
			)
			return strings.Compare(v1, v2)
		}
	}

	// Pad versions to same length
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for len(parts1) < maxLen {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < maxLen {
		parts2 = append(parts2, "0")
	}

	for i := range maxLen {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	return 0
}

// prompt user for update confirmation
func promptForUpdate(dependencyName, currentVersion, latestVersion string) bool {
	logger.Infof("%s is installed (version %s) but a newer version is available (%s).",
		dependencyName, currentVersion, latestVersion)
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

// findBinaryInArchive recursively searches for a binary in extracted archive
func findBinaryInArchive(extractDir, binaryName string) (string, error) {
	var foundPath string

	err := filepath.WalkDir(extractDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileName := d.Name()

		// Check for exact match first
		if fileName == binaryName {
			foundPath = path
			return filepath.SkipAll // Stop searching once found
		}

		// Check for pattern matches (e.g., terraform_1.5.0_linux_amd64 contains terraform)
		if strings.Contains(fileName, binaryName) && !strings.Contains(fileName, ".") {
			// Additional validation: check if it's likely an executable (no extension or common binary patterns)
			if !strings.Contains(fileName, ".txt") && !strings.Contains(fileName, ".md") &&
				!strings.Contains(fileName, ".json") && !strings.Contains(fileName, ".yml") &&
				!strings.Contains(fileName, ".yaml") {
				foundPath = path
				return filepath.SkipAll
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("could not find %s binary", binaryName)
	}

	return foundPath, nil
}

// installing dependencies doesn't matter the operating system
func install(url, name string) {
	currentOS := entities.GetOS()
	tempFilePath := path.Join(currentOS.GetTempDir(), name)
	destPath := path.Join(currentOS.GetInstallationPath(), name)

	// Ensure installation directory exists
	installDir := currentOS.GetInstallationPath()
	if err := os.MkdirAll(installDir, 0750); err != nil {
		logger.Fatalf("Failed to create installation directory %s: %s", installDir, err)
	}

	logger.Infof("Downloading %s from %s...", name, url)
	if err := currentOS.Download(url, tempFilePath); err != nil {
		logger.Fatalf("Failed to download %s: %s", name, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	fileTypeCmd := exec.CommandContext(ctx, "file", tempFilePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
	}

	//nolint:nestif // Complex file type detection logic requires nested conditions
	if strings.Contains(string(fileTypeOutput), "Zip archive data") {
		logger.Infof("%s is a zip file, extracting...", name)
		// Create a temporary directory for extraction
		extractDir := path.Join(currentOS.GetTempDir(), name+"_extract")
		if err = currentOS.Extract(tempFilePath, extractDir); err != nil {
			logger.Fatalf("Failed to extract %s: %s", name, err)
		}

		// Find the actual binary in the extracted directory using recursive search
		var binaryPath string
		binaryPath, err = findBinaryInArchive(extractDir, name)
		if err != nil {
			logger.Fatalf("Failed to find %s binary in extracted archive: %s", name, err)
		}

		// Move the binary to the destination
		if err = currentOS.Move(binaryPath, destPath); err != nil {
			logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
		}

		// Clean up
		if err = currentOS.Remove(tempFilePath); err != nil {
			logger.Fatalf("Failed to remove %s: %s", name, err)
		}
		if err = os.RemoveAll(extractDir); err != nil {
			logger.Fatalf("Failed to remove extraction directory %s: %s", extractDir, err)
		}
	} else {
		if err = currentOS.Move(tempFilePath, destPath); err != nil {
			logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
		}
	}

	if err = currentOS.MakeExecutable(destPath); err != nil {
		logger.Fatalf("Failed to make %s executable: %s", name, err)
	}
}
