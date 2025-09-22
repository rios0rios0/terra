package entities

import (
	"go.uber.org/dig"
)

// RegisterProviders registers all entity providers with the DIG container.
func RegisterProviders(container *dig.Container) error {
	if err := container.Provide(NewSettings); err != nil {
		return err
	}
	if err := container.Provide(NewCLI); err != nil {
		return err
	}
	return nil
}
