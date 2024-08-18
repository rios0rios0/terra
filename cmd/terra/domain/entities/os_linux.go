package entities

import (
	"os"
	"os/exec"
)

const osOrwxGrxUx = 0o755

type OSLinux struct{}

func (it OSLinux) Download(url, tempFilePath string) error {
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout
	return curlCmd.Run()
}

func (it OSLinux) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	return unzipCmd.Run()
}

func (it OSLinux) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("mv", tempFilePath, destPath)
	return mvCmd.Run()
}

func (it OSLinux) Remove(tempFilePath string) error {
	rmCmd := exec.Command("rm", tempFilePath)
	return rmCmd.Run()
}

func (it OSLinux) MakeExecutable(filePath string) error {
	return os.Chmod(filePath, osOrwxGrxUx)
}

func (it OSLinux) GetTempDir() string {
	return "/tmp"
}

func (it OSLinux) GetInstallationPath() string {
	return "/usr/local/bin"
}

func GetOS() OS {
	return OSLinux{}
}
