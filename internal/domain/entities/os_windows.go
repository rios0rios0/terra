package entities

import (
	"os"
	"os/exec"
)

type OSWindows struct{}

func (it OSWindows) Download(url, tempFilePath string) error {
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout
	return curlCmd.Run()
}

func (it OSWindows) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("powershell", "Expand-Archive", "-Path", tempFilePath, "-DestinationPath", destPath, "-Force")
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	return unzipCmd.Run()
}

func (it OSWindows) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("move", tempFilePath, destPath)
	return mvCmd.Run()
}

func (it OSWindows) Remove(tempFilePath string) error {
	rmCmd := exec.Command("del", tempFilePath)
	return rmCmd.Run()
}

func (it OSWindows) MakeExecutable(_ string) error {
	return nil // Windows doesn't need to explicitly make files executable
}

func (it OSWindows) GetTempDir() string {
	return os.Getenv("TEMP")
}

func (it OSWindows) GetInstallationPath() string {
	return os.Getenv("ProgramFiles")
}

func GetOS() OS {
	return OSWindows{}
}
