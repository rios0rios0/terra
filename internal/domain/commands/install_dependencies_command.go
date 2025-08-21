package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
)

const contextTimeout = 10 * time.Second

type InstallDependenciesCommand struct{}

// getHTTPClient returns an HTTP client configured with proxy settings if available
func getHTTPClient() *http.Client {
	transport := &http.Transport{}

	if httpsProxy := os.Getenv("TERRA_HTTPS_PROXY"); httpsProxy != "" {
		if proxyURL, err := url.Parse(httpsProxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	} else if httpProxy := os.Getenv("TERRA_HTTP_PROXY"); httpProxy != "" {
		if proxyURL, err := url.Parse(httpProxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   contextTimeout,
	}
}

func NewInstallDependenciesCommand() *InstallDependenciesCommand {
	return &InstallDependenciesCommand{}
}

func (it *InstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	// Pre-check connectivity to all required endpoints
	checkConnectivity(dependencies)

	for _, dependency := range dependencies {
		latestVersion := fetchLatestVersion(dependency.VersionURL, dependency.RegexVersion)

		if !isDependencyCLIAvailable(dependency.CLI) {
			ensureRootPrivileges()
			logger.Warnf("%s is not installed, installing now...", dependency.Name)
			install(fmt.Sprintf(dependency.BinaryURL, latestVersion), dependency.CLI)
		}
	}
}

// check connectivity to required endpoints before attempting downloads
func checkConnectivity(dependencies []entities.Dependency) {
	logger.Info("Checking connectivity to dependency download endpoints...")

	var unreachableEndpoints []string

	for _, dependency := range dependencies {
		// Check version URL
		if !isEndpointReachable(dependency.VersionURL) {
			unreachableEndpoints = append(unreachableEndpoints, dependency.VersionURL)
		}

		// Extract base URL from binary URL for connectivity check
		if baseURL := extractBaseURL(dependency.BinaryURL); baseURL != "" && !isEndpointReachable(baseURL) {
			unreachableEndpoints = append(unreachableEndpoints, baseURL)
		}
	}

	if len(unreachableEndpoints) > 0 {
		logger.Warnf("Unable to reach the following endpoints:")
		for _, endpoint := range unreachableEndpoints {
			logger.Warnf("  - %s", endpoint)
		}
		logger.Warnf("This may indicate firewall restrictions or network connectivity issues.")
		logger.Warnf("Required firewall rules:")
		logger.Warnf("  - Allow outbound HTTPS (port 443) to:")
		logger.Warnf("    * releases.hashicorp.com (Terraform downloads)")
		logger.Warnf("    * checkpoint-api.hashicorp.com (Terraform version checks)")
		logger.Warnf("    * github.com (Terragrunt downloads)")
		logger.Warnf("    * api.github.com (Terragrunt version checks)")
		logger.Warnf("Consider using proxy settings or URL overrides via environment variables.")
		logger.Warnf("See documentation for details on network requirements and configuration.")
	} else {
		logger.Info("All endpoints are reachable.")
	}
}

// check if an endpoint is reachable within the timeout
func isEndpointReachable(endpoint string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, endpoint, nil)
	if err != nil {
		return false
	}

	client := getHTTPClient()
	client.Timeout = 5 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx status codes as reachable
	return resp.StatusCode < 400
}

// extract base URL for connectivity testing from a template URL
func extractBaseURL(templateURL string) string {
	// Handle URLs with format specifiers
	if strings.Contains(templateURL, "%") {
		templateURL = strings.Split(templateURL, "%")[0]
		if strings.HasSuffix(templateURL, "/") {
			templateURL = strings.TrimSuffix(templateURL, "/")
		}
	}

	parsedURL, err := url.Parse(templateURL)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
}

// fetch the latest version of software from a URL
func fetchLatestVersion(url, regexPattern string) string {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Fatalf("Error creating request: %s", err)
	}

	resp, err := getHTTPClient().Do(req)
	if err != nil {
		logger.Fatalf("Error fetching version info from %s: %s\nThis may indicate network connectivity issues or firewall restrictions.\nSee documentation for network requirements and configuration options.", url, err)
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
	cmd := exec.Command(name, "-v")
	return cmd.Run() == nil
}

// check if the "terra" has root privileges to install dependencies
func ensureRootPrivileges() {
	if os.Geteuid() != 0 {
		logger.Fatalf("Run this command with root privileges to install the dependencies")
		return
	}
}

// installing dependencies doesn't matter the operating system
func install(url, name string) {
	currentOS := entities.GetOS()
	tempFilePath := path.Join(currentOS.GetTempDir(), name)
	destPath := path.Join(currentOS.GetInstallationPath(), name)

	logger.Infof("Downloading %s from %s...", name, url)
	if err := currentOS.Download(url, tempFilePath); err != nil {
		logger.Fatalf("Failed to download %s from %s: %s\nThis may indicate network connectivity issues or firewall restrictions.\nEnsure outbound HTTPS (port 443) access to required domains.\nSee documentation for network requirements and proxy configuration.", name, url, err)
	}

	fileTypeCmd := exec.Command("file", tempFilePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
	}

	if strings.Contains(string(fileTypeOutput), "Zip archive data") {
		logger.Infof("%s is a zip file, extracting...", name)
		if err := currentOS.Extract(tempFilePath, destPath); err != nil {
			logger.Fatalf("Failed to extract %s: %s", name, err)
		}
		if err := currentOS.Remove(tempFilePath); err != nil {
			logger.Fatalf("Failed to remove %s: %s", name, err)
		}
	} else {
		if err := currentOS.Move(tempFilePath, destPath); err != nil {
			logger.Fatalf("Failed to move %s to %s: %s", name, destPath, err)
		}
	}

	if err := currentOS.MakeExecutable(destPath); err != nil {
		logger.Fatalf("Failed to make %s executable: %s", name, err)
	}
}
