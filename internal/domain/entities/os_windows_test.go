package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestOSWindows_ShouldDownloadSuccessfully_WhenValidURLProvided(t *testing.T) {
	// GIVEN: A Windows OS implementation
	osImpl := &entities.OSWindows{}

	// WHEN: Testing download functionality
	// THEN: Should download successfully (using helper for consistent testing)
	entities.HelperDownloadSuccess(t, osImpl, "test_download_windows")
}

func TestOSWindows_ShouldReturnError_WhenHTTPErrorOccurs(t *testing.T) {
	// GIVEN: A Windows OS implementation
	osImpl := &entities.OSWindows{}

	// WHEN: Testing download with HTTP error
	// THEN: Should handle HTTP error appropriately (using helper for consistent testing)
	entities.HelperDownloadHTTPError(t, osImpl, "test_download_windows")
}
