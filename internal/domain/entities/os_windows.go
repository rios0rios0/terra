package entities

import (
	"os"
)

type OSWindows struct{}

func (it *OSWindows) Download(url, tempFilePath string) error {
	return downloadFile(url, tempFilePath)
}

func (it *OSWindows) Extract(tempFilePath, destPath string) error {
	return extractZipArchive(tempFilePath, destPath)
}

func (it *OSWindows) Move(tempFilePath, destPath string) error {
	return moveFile(tempFilePath, destPath)
}

func (it *OSWindows) Remove(tempFilePath string) error {
	return removeFile(tempFilePath)
}

func (it *OSWindows) MakeExecutable(_ string) error {
	return nil // Windows doesn't need to explicitly make files executable
}

func (it *OSWindows) GetTempDir() string {
	return os.Getenv("TEMP")
}

func (it *OSWindows) GetInstallationPath() string {
	// Allow override via environment variable (used by tests to avoid
	// overwriting real binaries).
	if envPath := os.Getenv("TERRA_INSTALL_PATH"); envPath != "" {
		return envPath
	}
	return os.Getenv("ProgramFiles")
}

func GetOS() *OSWindows {
	return &OSWindows{}
}
