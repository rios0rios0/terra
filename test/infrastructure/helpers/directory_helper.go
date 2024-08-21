package helpers

import (
	"os"
	"path/filepath"
)

type DirectoryHelper struct{}

func (it *DirectoryHelper) CreateTestDirectories() []string {
	dirs := []string{"cache", "temp"}
	for _, dir := range dirs {
		_ = os.MkdirAll(filepath.Join(".", dir), os.ModePerm)
	}
	return dirs
}

func (it *DirectoryHelper) DirectoryExists(directory string) bool {
	info, err := os.Stat(filepath.Join(".", directory))
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
