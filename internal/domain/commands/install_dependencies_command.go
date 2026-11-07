package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rios0rios0/cliforge/pkg/selfupdate"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const contextTimeout = 10 * time.Second

// assumeYesEnvVar, when set to a truthy value, makes dependency-update prompts
// auto-confirm. This keeps `terra install` non-interactive in CI, where
// blocking on stdin would hang the job forever.
const assumeYesEnvVar = "TERRA_ASSUME_YES"

// githubTokenEnvVar names the token `terra install` sends to the GitHub API when
// it is set. Without one, GitHub allows 60 lookups an hour per IP address, and
// machines sharing an address can use that up between them.
const githubTokenEnvVar = "GITHUB_TOKEN" //nolint:gosec // the name of an environment variable, not a credential

// githubAPIURL scopes that token: no other host may ever receive it.
const githubAPIURL = "https://api.github.com/"

const (
	// versionLookupAttempts is how many times a version lookup is tried.
	versionLookupAttempts = 3
	// versionLookupBackoff is the wait before the second attempt. It grows with
	// each attempt after that.
	versionLookupBackoff = time.Second
	// versionLookupErrorBodyLimit bounds how much of a failed response is logged.
	versionLookupErrorBodyLimit = 200
)

var (
	errNoVersionMatch      = errors.New("no version match found, check the regex pattern")
	errVersionLookupStatus = errors.New("version lookup failed")
)

type InstallDependenciesCommand struct{}

func NewInstallDependenciesCommand() *InstallDependenciesCommand {
	return &InstallDependenciesCommand{}
}

func (it *InstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	for _, dependency := range dependencies {
		latestVersion, err := fetchLatestVersion(dependency.VersionURL, dependency.RegexVersion)
		if err != nil {
			handleFailedVersionLookup(dependency, err)
			continue
		}

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

			comparison := selfupdate.CompareVersions(currentVersion, latestVersion)
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
				logger.Infof(
					"%s version %s is newer than latest available %s",
					dependency.Name,
					currentVersion,
					latestVersion,
				)
			}
		}
	}
}

// handleFailedVersionLookup keeps a tool that is already installed, because the
// lookup only decides whether to update it. A missing tool has no version to
// install, so that case stays fatal.
func handleFailedVersionLookup(dependency entities.Dependency, err error) {
	if !isDependencyCLIAvailable(dependency.CLI) {
		logger.Fatalf("Could not find the latest %s version to install: %s", dependency.Name, err)
	}
	logger.Warnf("Could not check for a newer %s, keeping the installed one: %s", dependency.Name, err)
}

// fetch the latest version of software from a URL, retrying the failures that
// can pass on their own.
func fetchLatestVersion(url, regexPattern string) (string, error) {
	var lastErr error
	for attempt := range versionLookupAttempts {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * versionLookupBackoff)
		}

		body, retryable, err := requestVersionInfo(url)
		if err == nil {
			return matchVersion(body, regexPattern)
		}

		lastErr = err
		if !retryable {
			break
		}
		logger.Debugf("Version lookup attempt %d of %d failed: %s", attempt+1, versionLookupAttempts, err)
	}

	return "", lastErr
}

// FetchLatestVersionPublic exposes fetchLatestVersion for external (black-box) tests.
func FetchLatestVersionPublic(url, regexPattern string) (string, error) {
	return fetchLatestVersion(url, regexPattern)
}

// requestVersionInfo makes one lookup. It returns the response body and, on a
// failure, whether trying again could help. A non-200 is reported with its status
// and the start of its body, so a rate limit does not read as a regex mismatch.
func requestVersionInfo(url string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", false, fmt.Errorf("creating the version request: %w", err)
	}
	authorizeVersionRequest(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", true, fmt.Errorf("fetching version info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", true, fmt.Errorf("reading the version response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", isRetryableVersionStatus(resp.StatusCode), fmt.Errorf(
			"%w: HTTP %d from %s: %s", errVersionLookupStatus, resp.StatusCode, url, summarizeBody(body),
		)
	}

	return string(body), false, nil
}

// authorizeVersionRequest adds the GitHub token to a request for the GitHub API,
// and to nothing else.
func authorizeVersionRequest(req *http.Request) {
	token := os.Getenv(githubTokenEnvVar)
	if token != "" && strings.HasPrefix(req.URL.String(), githubAPIURL) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// AuthorizeVersionRequestPublic exposes authorizeVersionRequest for external (black-box) tests.
func AuthorizeVersionRequestPublic(req *http.Request) {
	authorizeVersionRequest(req)
}

// isRetryableVersionStatus reports whether a status can pass on its own. GitHub
// answers 403 as well as 429 when a client is rate limited.
func isRetryableVersionStatus(status int) bool {
	return status == http.StatusForbidden ||
		status == http.StatusTooManyRequests ||
		status >= http.StatusInternalServerError
}

// summarizeBody keeps a failed response short enough to log.
func summarizeBody(body []byte) string {
	summary := strings.TrimSpace(string(body))
	if len(summary) > versionLookupErrorBodyLimit {
		return summary[:versionLookupErrorBodyLimit] + "..."
	}

	return summary
}

// matchVersion reads the version out of a successful response.
func matchVersion(body, regexPattern string) (string, error) {
	matches := regexp.MustCompile(regexPattern).FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("%w: %s", errNoVersionMatch, regexPattern)
}

// checking if a dependency is available.
func isDependencyCLIAvailable(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, "-v")
	return cmd.Run() == nil
}

// get current version of installed dependency.
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

// isTruthy reports whether an environment-variable value means affirmative.
func isTruthy(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "1" || normalized == "true" || normalized == "on" {
		return true
	}

	return strings.HasPrefix(normalized, "y")
}

// IsTruthyPublic exposes isTruthy for external (black-box) tests.
func IsTruthyPublic(value string) bool {
	return isTruthy(value)
}

// stdinIsInteractive reports whether stdin is a real terminal. A pipe, a file,
// or a closed stdin (as in CI) is not a terminal, so the caller must never
// block on it.
func stdinIsInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// promptForUpdate asks whether to update a dependency. In a non-interactive
// session (CI) it never blocks on stdin: TERRA_ASSUME_YES=true auto-confirms,
// otherwise the update is skipped with a warning.
func promptForUpdate(dependencyName, currentVersion, latestVersion string) bool {
	logger.Infof("%s is installed (version %s) but a newer version is available (%s).",
		dependencyName, currentVersion, latestVersion)

	if isTruthy(os.Getenv(assumeYesEnvVar)) {
		logger.Infof("%s is set - updating %s automatically.", assumeYesEnvVar, dependencyName)

		return true
	}
	if !stdinIsInteractive() {
		logger.Warnf(
			"Non-interactive session - skipping update of %s. Set %s=true to update automatically (e.g. in CI).",
			dependencyName, assumeYesEnvVar,
		)

		return false
	}

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

// findBinaryInArchive recursively searches for a binary in extracted archive.
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

// setupInstallationEnvironment prepares the temporary file and installation directory.
func setupInstallationEnvironment(name string, currentOS entities.OS) (string, string) {
	// Create a unique temporary file to avoid permission conflicts
	tempFile, err := os.CreateTemp(currentOS.GetTempDir(), name+"_*")
	if err != nil {
		logger.Fatalf("Failed to create temporary file for %s: %s", name, err)
	}
	tempFilePath := tempFile.Name()
	if closeErr := tempFile.Close(); closeErr != nil {
		logger.Warnf("Failed to close temporary file %s: %s", tempFilePath, closeErr)
	} // Close immediately since we'll overwrite it during download

	destPath := path.Join(currentOS.GetInstallationPath(), name)

	// Ensure installation directory exists
	installDir := currentOS.GetInstallationPath()
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	if mkdirErr := os.MkdirAll(installDir, 0750); mkdirErr != nil {
		logger.Fatalf("Failed to create installation directory %s: %s", installDir, mkdirErr)
	}

	return tempFilePath, destPath
}

// downloadDependency downloads the dependency to the specified temporary file.
func downloadDependency(url, name, tempFilePath string, currentOS entities.OS) {
	logger.Infof("Downloading %s from %s...", name, url)
	if downloadErr := currentOS.Download(url, tempFilePath); downloadErr != nil {
		logger.Fatalf("Failed to download %s: %s", name, downloadErr)
	}
}

// zipFileSignature is the local file header signature that opens every
// non-empty ZIP archive ("PK\x03\x04").
const zipFileSignature = "PK\x03\x04"

// detectFileType determines if the downloaded file is a zip archive by
// reading its magic number. Inspecting the bytes directly replaces a call
// to `file`, which has to be resolved through PATH and does not exist on
// Windows at all.
func detectFileType(tempFilePath, name string) bool {
	handle, err := os.Open(tempFilePath)
	if err != nil {
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
		return false
	}
	defer handle.Close()

	signature := make([]byte, len(zipFileSignature))
	if _, err = io.ReadFull(handle, signature); err != nil {
		// A file too short to hold the signature simply is not an archive.
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return false
		}
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
		return false
	}

	return string(signature) == zipFileSignature
}

// processArchive extracts the zip archive and moves the binary to destination.
func processArchive(tempFilePath, destPath, name string, currentOS entities.OS) {
	logger.Infof("%s is a zip file, extracting...", name)

	// Create a unique temporary directory for extraction
	extractDir, extractErr := os.MkdirTemp(currentOS.GetTempDir(), name+"_extract_*")
	if extractErr != nil {
		logger.Fatalf("Failed to create temporary extraction directory for %s: %s", name, extractErr)
	}

	if extractErr = currentOS.Extract(tempFilePath, extractDir); extractErr != nil {
		logger.Fatalf("Failed to extract %s: %s", name, extractErr)
	}

	// Find the actual binary in the extracted directory using recursive search
	binaryPath, err := findBinaryInArchive(extractDir, name)
	if err != nil {
		logger.Fatalf("Failed to find %s binary in extracted archive: %s", name, err)
	}

	// Move the binary to the destination
	if err = currentOS.Move(binaryPath, destPath); err != nil {
		logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
	}

	// Clean up temporary files
	cleanupArchiveFiles(tempFilePath, extractDir, name, currentOS)
}

// cleanupArchiveFiles removes temporary files after archive processing.
func cleanupArchiveFiles(tempFilePath, extractDir, name string, currentOS entities.OS) {
	if err := currentOS.Remove(tempFilePath); err != nil {
		logger.Fatalf("Failed to remove %s: %s", name, err)
	}
	if err := os.RemoveAll(extractDir); err != nil {
		logger.Fatalf("Failed to remove extraction directory %s: %s", extractDir, err)
	}
}

// installBinary moves the binary to destination and makes it executable.
func installBinary(tempFilePath, destPath, name string, currentOS entities.OS) {
	if err := currentOS.Move(tempFilePath, destPath); err != nil {
		logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
	}

	if err := currentOS.MakeExecutable(destPath); err != nil {
		logger.Fatalf("Failed to make %s executable: %s", name, err)
	}
}

// installing dependencies doesn't matter the operating system.
func install(url, name string) {
	currentOS := entities.GetOS()

	// Setup environment
	tempFilePath, destPath := setupInstallationEnvironment(name, currentOS)
	defer os.Remove(tempFilePath) // Ensure cleanup of the temporary file

	// Download the dependency
	downloadDependency(url, name, tempFilePath, currentOS)

	// Process based on file type
	if detectFileType(tempFilePath, name) {
		processArchive(tempFilePath, destPath, name, currentOS)
	} else {
		installBinary(tempFilePath, destPath, name, currentOS)
	}
}

// CompareVersionsPublic is a public wrapper for testing.
func CompareVersionsPublic(v1, v2 string) int {
	return selfupdate.CompareVersions(v1, v2)
}

// FindBinaryInArchivePublic is a public wrapper for testing the private findBinaryInArchive function.
func FindBinaryInArchivePublic(extractDir, binaryName string) (string, error) {
	return findBinaryInArchive(extractDir, binaryName)
}
