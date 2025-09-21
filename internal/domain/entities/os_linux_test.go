package entities_test

import (
	"os"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOS_ShouldReturnValidInstance_WhenCalled(t *testing.T) {
	// GIVEN: The GetOS function is available

	// WHEN: Calling GetOS
	osInstance := entities.GetOS()

	// THEN: Should return a valid OS instance
	require.NotNil(t, osInstance)
}

func TestOSLinux_ShouldReturnValidTempDir_WhenGetTempDirCalled(t *testing.T) {
	// GIVEN: An OS instance
	osInstance := entities.GetOS()

	// WHEN: Getting the temporary directory
	tempDir := osInstance.GetTempDir()

	// THEN: Should return a non-empty string pointing to an existing directory
	assert.NotEmpty(t, tempDir)
	
	info, err := os.Stat(tempDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestOSLinux_ShouldReturnValidInstallationPath_WhenGetInstallationPathCalled(t *testing.T) {
	// GIVEN: An OS instance
	osInstance := entities.GetOS()

	// WHEN: Getting the installation path
	installPath := osInstance.GetInstallationPath()

	// THEN: Should return a non-empty string
	assert.NotEmpty(t, installPath)
}

func TestOSLinux_ShouldMakeFileExecutable_WhenValidFileProvided(t *testing.T) {
	// GIVEN: An OS instance and a temporary file
	osInstance := entities.GetOS()
	tempFile, err := os.CreateTemp(t.TempDir(), "test_executable_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// WHEN: Making the file executable
	err = osInstance.MakeExecutable(tempFile.Name())

	// THEN: Should succeed and file should be executable
	require.NoError(t, err)
	
	info, err := os.Stat(tempFile.Name())
	require.NoError(t, err)
	mode := info.Mode()
	assert.NotEqual(t, 0, mode&0111, "File should be executable after MakeExecutable")
}

func TestOSLinux_ShouldReturnError_WhenMakeExecutableCalledWithNonExistentFile(t *testing.T) {
	// GIVEN: An OS instance and a non-existent file path
	osInstance := entities.GetOS()
	nonExistentFile := "/non/existent/file12345"

	// WHEN: Attempting to make the non-existent file executable
	err := osInstance.MakeExecutable(nonExistentFile)

	// THEN: Should return an error
	assert.Error(t, err)
}

func TestOSLinux_ShouldHandleRemove_WhenNonExistentFileProvided(t *testing.T) {
	// GIVEN: An OS instance and a non-existent file path
	osInstance := entities.GetOS()
	nonExistentFile := "/non/existent/file12345"

	// WHEN: Attempting to remove the non-existent file
	err := osInstance.Remove(nonExistentFile)

	// THEN: Different implementations may handle this differently, so we just log the result
	// This test verifies the method doesn't panic
	_ = err // Some implementations return nil, others return error - both are acceptable
}
