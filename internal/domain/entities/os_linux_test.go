package entities

import (
	"testing"
)

func TestOSLinux_Download(t *testing.T) {
	testDownloadSuccess(t, &OSLinux{}, "test_download_linux")
}

func TestOSLinux_DownloadHTTPError(t *testing.T) {
	testDownloadHTTPError(t, &OSLinux{}, "test_download_linux")
}
