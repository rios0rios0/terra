package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

type RunFromRoot interface {
	Execute(targetPath string, arguments []string, dependencies []entities.Dependency)
}
