package internal

import "github.com/rios0rios0/terra/internal/domain/entities"

// AppInternal is not totally necessary, but it is rather a good example for other applications.
type AppInternal struct {
	controllers []entities.Controller
}

func NewAppInternal(controllers *[]entities.Controller) *AppInternal {
	return &AppInternal{controllers: *controllers}
}

func (it *AppInternal) GetControllers() []entities.Controller {
	return it.controllers
}
