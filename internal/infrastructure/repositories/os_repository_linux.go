package repositories

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const osOrwxGrxUx = 0o755

type OsRepository struct{}

func NewDefaultOSRepository() *OsRepository {
	return &OsRepository{}
}

func (it *OsRepository) Download(url, tempFilePath string) error {
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout
	err := curlCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform download using 'cURL': %w", err)
	}
	return err
}

func (it *OsRepository) GetTempDir() string {
	return "/tmp"
}

func (it *OsRepository) GetInstallationPath() string {
	return "/usr/local/bin"
}

func (it *OsRepository) extract(tempFilePath, destPath string) error {
	unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", destPath)
	unzipCmd.Stderr = os.Stderr
	unzipCmd.Stdout = os.Stdout
	err := unzipCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform decompressing using 'zip': %w", err)
	}
	return err
}

func (it *OsRepository) move(tempFilePath, destPath string) error {
	mvCmd := exec.Command("mv", tempFilePath, destPath)
	err := mvCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform moving folder using 'mv': %w", err)
	}
	return err
}

func (it *OsRepository) remove(tempFilePath string) error {
	rmCmd := exec.Command("rm", tempFilePath)
	err := rmCmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to perform deleting folder using 'rm': %w", err)
	}
	return err
}

func (it *OsRepository) makeExecutable(filePath string) error {
	err := os.Chmod(filePath, osOrwxGrxUx)
	if err != nil {
		err = fmt.Errorf("failed to perform change binary permissions using 'chmod': %w", err)
	}
	return err
}

func isZipFile(filePath string) bool {
	fileTypeCmd := exec.Command("file", filePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Errorf("Failed to determine file type of %s: %s", filePath, err)
		return false
	}
	return strings.Contains(string(fileTypeOutput), "Zip archive data")
}
