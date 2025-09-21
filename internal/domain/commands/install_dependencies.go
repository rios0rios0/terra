package commands

import "github.com/rios0rios0/terra/internal/domain/entities"

//nolint:iface // Different semantic purpose than FormatFiles
type InstallDependencies interface {
	Execute(dependencies []entities.Dependency)
}
