package commands

import "github.com/rios0rios0/terra/cmd/terra/domain/entities"

type InstallDependencies interface {
	Execute(dependencies []entities.Dependency)
}
