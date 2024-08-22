package repositories

import "github.com/rios0rios0/terra/internal/domain/entities"

// OSRepository is not totally necessary, but it is rather a good example for other applications
type OSRepository interface {
	ExecuteCommand(command string, arguments []string, directory string) error
	InstallExecutable(sourcePath, destinationPath string, currentOS entities.OS) error
}
