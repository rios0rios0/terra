package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

type RunFromRoot interface {
	Execute(dependencies []entities.Dependency)
}
