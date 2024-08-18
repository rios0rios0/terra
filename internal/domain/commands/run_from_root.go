package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

type RunFromRoot interface {
	Execute(targetDirectory string, dependencies []entities.Dependency)
}
