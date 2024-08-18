package entities

type OS interface {
	Download(url, tempFilePath string) error
	Extract(tempFilePath, destPath string) error
	Move(tempFilePath, destPath string) error
	Remove(tempFilePath string) error
	MakeExecutable(filePath string) error
	GetTempDir() string
	GetInstallationPath() string
}
