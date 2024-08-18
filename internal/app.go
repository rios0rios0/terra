package internal

import "github.com/rios0rios0/terra/internal/domain/entities"

type AppCLI struct {
	controllers []entities.Controller
}

func NewAppCLI(controllers []entities.Controller) *AppCLI {
	return &AppCLI{controllers: controllers}
}

func (it AppCLI) GetControllers() []entities.Controller {
	return it.controllers
}
