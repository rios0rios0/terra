package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/infrastructure/repository_helpers"
)

func TestOSWindows_Download(t *testing.T) {
	t.Parallel()
	
	t.Run("should download successfully when valid URL provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Windows OS implementation
		osImpl := &entities.OSWindows{}

		// WHEN: Testing download functionality
		// THEN: Should download successfully (using helper for consistent testing)
		repositories_helpers.HelperDownloadSuccess(t, osImpl, "test_download_windows")
	})
	
	t.Run("should return error when HTTP error occurs", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A Windows OS implementation
		osImpl := &entities.OSWindows{}

		// WHEN: Testing download with HTTP error
		// THEN: Should handle HTTP error appropriately (using helper for consistent testing)
		repositories_helpers.HelperDownloadHTTPError(t, osImpl, "test_download_windows")
	})
}
