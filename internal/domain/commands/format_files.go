package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

type FormatFiles interface {
	Execute(dependencies []entities.Dependency)
}
