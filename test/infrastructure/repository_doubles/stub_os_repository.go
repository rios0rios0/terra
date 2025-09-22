package repository_doubles

import (
	"fmt"
	"strings"
)

// StubOSRepository provides a stub implementation of OS interface for testing
type StubOSRepository struct {
	DownloadError      error
	DownloadCallCount  int
	DownloadURLs       []string
	DownloadPaths      []string
	
	MoveError         error
	MoveCallCount     int
	MoveSources       []string
	MoveDestinations  []string
	
	RemoveError      error
	RemoveCallCount  int
	RemovePaths      []string
	
	MakeExecutableError     error
	MakeExecutableCallCount int
	MakeExecutablePaths     []string
	
	TempDir         string
	InstallationPath string
}

func NewStubOSRepository() *StubOSRepository {
	return &StubOSRepository{
		TempDir:          "/tmp",
		InstallationPath: "/usr/local/bin",
		DownloadURLs:     make([]string, 0),
		DownloadPaths:    make([]string, 0),
		MoveSources:      make([]string, 0),
		MoveDestinations: make([]string, 0),
		RemovePaths:      make([]string, 0),
		MakeExecutablePaths: make([]string, 0),
	}
}

func (s *StubOSRepository) Download(url, tempFilePath string) error {
	s.DownloadCallCount++
	s.DownloadURLs = append(s.DownloadURLs, url)
	s.DownloadPaths = append(s.DownloadPaths, tempFilePath)
	return s.DownloadError
}

func (s *StubOSRepository) Extract(tempFilePath, destPath string) error {
	return fmt.Errorf("extract not implemented in stub")
}

func (s *StubOSRepository) Move(tempFilePath, destPath string) error {
	s.MoveCallCount++
	s.MoveSources = append(s.MoveSources, tempFilePath)
	s.MoveDestinations = append(s.MoveDestinations, destPath)
	return s.MoveError
}

func (s *StubOSRepository) Remove(tempFilePath string) error {
	s.RemoveCallCount++
	s.RemovePaths = append(s.RemovePaths, tempFilePath)
	return s.RemoveError
}

func (s *StubOSRepository) MakeExecutable(filePath string) error {
	s.MakeExecutableCallCount++
	s.MakeExecutablePaths = append(s.MakeExecutablePaths, filePath)
	return s.MakeExecutableError
}

func (s *StubOSRepository) GetTempDir() string {
	return s.TempDir
}

func (s *StubOSRepository) GetInstallationPath() string {
	return s.InstallationPath
}

// GetLastDownloadURL returns the last URL used for download
func (s *StubOSRepository) GetLastDownloadURL() string {
	if len(s.DownloadURLs) == 0 {
		return ""
	}
	return s.DownloadURLs[len(s.DownloadURLs)-1]
}

// GetLastDownloadPath returns the last path used for download
func (s *StubOSRepository) GetLastDownloadPath() string {
	if len(s.DownloadPaths) == 0 {
		return ""
	}
	return s.DownloadPaths[len(s.DownloadPaths)-1]
}

// WasMethodCalled checks if any of the specified methods was called
func (s *StubOSRepository) WasMethodCalled(method string) bool {
	switch strings.ToLower(method) {
	case "download":
		return s.DownloadCallCount > 0
	case "move":
		return s.MoveCallCount > 0
	case "remove":
		return s.RemoveCallCount > 0
	case "makeexecutable":
		return s.MakeExecutableCallCount > 0
	default:
		return false
	}
}