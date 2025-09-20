package entities

import (
	"fmt"
	"os"
	"os/exec"
)

const osOrwxGrxUx = 0o755

type OSLinux struct{}

func (it *OSLinux) Download(url, tempFilePath string) error {
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout
	err := curlCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform download using 'cURL': %w", err)
	}
	return err
}

func (it *OSLinux) Extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	err := unzipCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform decompressing using 'zip': %w", err)
	}
	return err
}

func (it *OSLinux) Move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("mv", tempFilePath, destPath)
	err := mvCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform moving folder using 'mv': %w", err)
	}
	return err
}

func (it *OSLinux) Remove(tempFilePath string) error {
	rmCmd := exec.Command("rm", tempFilePath)
	err := rmCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform deleting folder using 'rm': %w", err)
	}
	return err
}

func (it *OSLinux) MakeExecutable(filePath string) error {
	err := os.Chmod(filePath, osOrwxGrxUx)
	if err != nil {
		err = fmt.Errorf("failed to perform change binary permissions using 'chmod': %w", err)
	}
	return err
}

func (it *OSLinux) GetTempDir() string {
	return "/tmp"
}

func (it *OSLinux) GetInstallationPath() string {
	return "~/.local/bin"
}

func GetOS() *OSLinux {
	return &OSLinux{}
}
