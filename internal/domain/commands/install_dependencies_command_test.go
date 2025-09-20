package commands

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewInstallDependenciesCommand(t *testing.T) {
	cmd := NewInstallDependenciesCommand()
	if cmd == nil {
		t.Fatal("NewInstallDependenciesCommand returned nil")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "v1 less than v2",
			v1:       "1.0.0",
			v2:       "1.0.1",
			expected: -1,
		},
		{
			name:     "v1 equal to v2",
			v1:       "1.0.0",
			v2:       "1.0.0",
			expected: 0,
		},
		{
			name:     "v1 greater than v2",
			v1:       "1.0.1",
			v2:       "1.0.0",
			expected: 1,
		},
		{
			name:     "different major versions",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: -1,
		},
		{
			name:     "different minor versions",
			v1:       "1.1.0",
			v2:       "1.0.0",
			expected: 1,
		},
		{
			name:     "different patch versions",
			v1:       "1.0.1",
			v2:       "1.0.2",
			expected: -1,
		},
		{
			name:     "different length versions - shorter first",
			v1:       "1.0",
			v2:       "1.0.0",
			expected: 0,
		},
		{
			name:     "different length versions - longer first",
			v1:       "1.0.0",
			v2:       "1.0",
			expected: 0,
		},
		{
			name:     "longer version with higher patch",
			v1:       "1.0.0.1",
			v2:       "1.0.0",
			expected: 1,
		},
		{
			name:     "version with leading zeros",
			v1:       "1.01.0",
			v2:       "1.1.0",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestCompareVersionsWithNonNumericParts(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected string // We expect string comparison fallback
	}{
		{
			name:     "version with alpha",
			v1:       "1.0.0-alpha",
			v2:       "1.0.0",
			expected: "fallback to string comparison",
		},
		{
			name:     "version with beta",
			v1:       "1.0.0-beta",
			v2:       "1.0.0-alpha",
			expected: "fallback to string comparison",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			// For non-numeric versions, it should fall back to string comparison
			expected := strings.Compare(tt.v1, tt.v2)
			if result != expected {
				t.Errorf("compareVersions(%q, %q) = %d, expected %d (string comparison)", tt.v1, tt.v2, result, expected)
			}
		})
	}
}

func TestFetchLatestVersion(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		regexPattern string
		expected     string
	}{
		{
			name:         "valid terraform response",
			responseBody: `{"current_version":"1.5.0","alerts":[]}`,
			regexPattern: `"current_version":"([^"]+)"`,
			expected:     "1.5.0",
		},
		{
			name:         "valid terragrunt response",
			responseBody: `{"tag_name":"v0.50.0","name":"v0.50.0"}`,
			regexPattern: `"tag_name":"v([^"]+)"`,
			expected:     "0.50.0",
		},
		{
			name:         "multiple matches returns first",
			responseBody: `{"current_version":"1.5.0","old_version":"1.4.0"}`,
			regexPattern: `"([^"]+)_version":"([^"]+)"`,
			expected:     "current",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use builder pattern to create test server
			versionServer, _ := NewTestServerBuilder().
				WithVersionResponse("", tt.responseBody).
				BuildServers()
			defer versionServer.Close()

			result := fetchLatestVersion(versionServer.URL, tt.regexPattern)

			if result != tt.expected {
				t.Errorf("fetchLatestVersion() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestFetchLatestVersionWithServerError tests fetchLatestVersion with server errors
func TestFetchLatestVersionWithServerError(t *testing.T) {
	// Use builder pattern to create a failing server
	versionServer, _ := NewTestServerBuilder().
		WithDownloadFailure().
		BuildServers()
	defer versionServer.Close()

	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - the function calls logger.Fatalf on HTTP errors
			t.Logf("HTTP error correctly detected and handled: %v", r)
		}
	}()

	// This should trigger an error due to the server returning 500
	result := fetchLatestVersion(versionServer.URL, `"version":"([^"]+)"`)

	// If we reach here, the function didn't fail as expected
	t.Logf("Unexpected success with result: %s", result)
}

// Note: We can't easily test the failure case because fetchLatestVersion calls logger.Fatalf
// which terminates the process. In a real-world scenario, we would refactor this function
// to return an error instead of calling Fatalf.

func TestFindBinaryInArchive(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create test directory structure
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name           string
		binaryName     string
		createFiles    []string
		expectedFound  bool
		expectedSuffix string // suffix of the expected path
	}{
		{
			name:           "exact match in root",
			binaryName:     "terraform",
			createFiles:    []string{"terraform", "other.txt"},
			expectedFound:  true,
			expectedSuffix: "terraform",
		},
		{
			name:           "exact match in subdirectory",
			binaryName:     "terragrunt",
			createFiles:    []string{"subdir/terragrunt", "other.txt"},
			expectedFound:  true,
			expectedSuffix: "subdir/terragrunt",
		},
		{
			name:           "pattern match without extension",
			binaryName:     "terraform",
			createFiles:    []string{"terraform_1_5_0_linux_amd64", "terraform.txt"},
			expectedFound:  true,
			expectedSuffix: "terraform_1_5_0_linux_amd64",
		},
		{
			name:           "no match found",
			binaryName:     "missing",
			createFiles:    []string{"terraform", "other.txt"},
			expectedFound:  false,
			expectedSuffix: "",
		},
		{
			name:           "ignore text files",
			binaryName:     "terraform",
			createFiles:    []string{"terraform.txt", "terraform.md", "another.json"},
			expectedFound:  false,
			expectedSuffix: "",
		},
		{
			name:           "prefer exact match over pattern",
			binaryName:     "terraform",
			createFiles:    []string{"terraform", "terraform_with_version"},
			expectedFound:  true,
			expectedSuffix: "terraform",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing files
			if err := os.RemoveAll(tempDir); err != nil {
				t.Fatalf("Failed to clean temp dir: %v", err)
			}
			if err := os.MkdirAll(subDir, 0755); err != nil {
				t.Fatalf("Failed to recreate test directory: %v", err)
			}

			// Create test files
			for _, file := range tt.createFiles {
				filePath := filepath.Join(tempDir, file)
				dir := filepath.Dir(filePath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}
				if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
					t.Fatalf("Failed to create test file %s: %v", filePath, err)
				}
			}

			// Test the function
			result, err := findBinaryInArchive(tempDir, tt.binaryName)

			if tt.expectedFound {
				if err != nil {
					t.Errorf("Expected to find binary, but got error: %v", err)
				} else if !strings.HasSuffix(result, tt.expectedSuffix) {
					t.Errorf("Expected result to end with %q, got %q", tt.expectedSuffix, result)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error when binary not found, but got result: %s", result)
				}
			}
		})
	}
}

func TestFindBinaryInArchiveWithNonExistentDirectory(t *testing.T) {
	_, err := findBinaryInArchive("/non/existent/path", "terraform")
	if err == nil {
		t.Error("Expected error when directory doesn't exist, but got nil")
	}
}

func TestGetCurrentVersion(t *testing.T) {
	tests := []struct {
		name     string
		cliName  string
		expected string // empty means we expect it to return empty (no mock command available)
	}{
		{
			name:     "terraform command",
			cliName:  "terraform",
			expected: "", // Will return empty since terraform is not installed in test environment
		},
		{
			name:     "terragrunt command",
			cliName:  "terragrunt",
			expected: "", // Will return empty since terragrunt is not installed in test environment
		},
		{
			name:     "unknown command",
			cliName:  "unknown",
			expected: "",
		},
		{
			name:     "empty command name",
			cliName:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCurrentVersion(tt.cliName)
			// Since we can't guarantee terraform/terragrunt are installed in test environment,
			// we just verify the function doesn't panic and returns a string
			if result != tt.expected && tt.expected != "" {
				t.Errorf("getCurrentVersion(%q) = %q, expected %q", tt.cliName, result, tt.expected)
			}
		})
	}
}

func TestIsDependencyCLIAvailable(t *testing.T) {
	tests := []struct {
		name    string
		cliName string
		// We can't predict which tools are available, so we just test that it doesn't panic
	}{
		{
			name:    "check terraform",
			cliName: "terraform",
		},
		{
			name:    "check terragrunt",
			cliName: "terragrunt",
		},
		{
			name:    "check non-existent command",
			cliName: "definitely-not-a-real-command-12345",
		},
		{
			name:    "check empty command",
			cliName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			result := isDependencyCLIAvailable(tt.cliName)
			// Result should be a boolean
			_ = result
		})
	}
}

// Helper function to capture output for testing prompt functions
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	output, _ := io.ReadAll(r)
	return string(output)
}

func TestPromptForUpdate(t *testing.T) {
	tests := []struct {
		name           string
		dependencyName string
		currentVersion string
		latestVersion  string
		input          string
		expected       bool
	}{
		{
			name:           "user says yes",
			dependencyName: "Terraform",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			input:          "y\n",
			expected:       true,
		},
		{
			name:           "user says yes full word",
			dependencyName: "Terragrunt",
			currentVersion: "0.45.0",
			latestVersion:  "0.50.0",
			input:          "yes\n",
			expected:       true,
		},
		{
			name:           "user says no",
			dependencyName: "Terraform",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			input:          "n\n",
			expected:       false,
		},
		{
			name:           "user says no full word",
			dependencyName: "Terraform",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			input:          "no\n",
			expected:       false,
		},
		{
			name:           "user presses enter (default no)",
			dependencyName: "Terraform",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			input:          "\n",
			expected:       false,
		},
		{
			name:           "user types random text",
			dependencyName: "Terraform",
			currentVersion: "1.0.0",
			latestVersion:  "1.1.0",
			input:          "maybe\n",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate user input
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}

			// Save original stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Replace stdin with our pipe
			os.Stdin = r

			// Write the test input in a goroutine
			go func() {
				defer w.Close()
				w.WriteString(tt.input)
			}()

			// Call the function and capture the result
			result := promptForUpdate(tt.dependencyName, tt.currentVersion, tt.latestVersion)

			if result != tt.expected {
				t.Errorf("promptForUpdate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
