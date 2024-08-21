package interfaces

import "github.com/rios0rios0/terra/internal/domain/entities"

type FormatFiles interface {
	Execute(dependencies []entities.Dependency, listeners FormatFilesListeners)
}

type FormatFilesListeners struct {
	OnSuccess func()
	OnError   func(error)
}
