package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

// ParallelState defines the interface for parallel state manipulation commands.
type ParallelState interface {
	Execute(targetPath string, arguments []string, dependencies []entities.Dependency) error
}
