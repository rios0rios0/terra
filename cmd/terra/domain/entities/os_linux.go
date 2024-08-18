package entities

import (
	"os"
	"os/exec"
)

const osOrwxGrxUx = 0o755

type OSLinux struct{}

func (l OSLinux) Download(url, tempFilePath string) error {
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout
	return curlCmd.Run()
}

func (l OSLinux) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	return unzipCmd.Run()
}

func (l OSLinux) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("mv", tempFilePath, destPath)
	return mvCmd.Run()
}

func (l OSLinux) Remove(tempFilePath string) error {
	rmCmd := exec.Command("rm", tempFilePath)
	return rmCmd.Run()
}

func (l OSLinux) MakeExecutable(filePath string) error {
	return os.Chmod(filePath, osOrwxGrxUx)
}

func (l OSLinux) GetTempDir() string {
	return "/tmp"
}

func (l OSLinux) GetInstallationPath() string {
	return "/usr/local/bin"
}

func GetOS() OS {
	return OSLinux{}
}
