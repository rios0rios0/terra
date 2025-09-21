//nolint:testpackage // Testing private functions and fields requires same package
package entities

import (
	"testing"
)

func TestOSWindows_Download(t *testing.T) {
	testDownloadSuccess(t, &OSWindows{}, "test_download_windows")
}

func TestOSWindows_DownloadHTTPError(t *testing.T) {
	testDownloadHTTPError(t, &OSWindows{}, "test_download_windows")
}
