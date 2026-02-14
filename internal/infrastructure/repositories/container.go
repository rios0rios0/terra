package repositories

import (
	"github.com/rios0rios0/terra/internal/domain/repositories"
	"go.uber.org/dig"
)

// RegisterProviders registers all repository providers with the DIG container.
func RegisterProviders(container *dig.Container) error {
	if err := container.Provide(NewStdShellRepository); err != nil {
		return err
	}
	if err := container.Provide(NewUpgradeAwareShellRepository); err != nil {
		return err
	}
	if err := container.Provide(NewInteractiveShellRepository); err != nil {
		return err
	}
	// Bind interface to implementation
	if err := container.Provide(func(impl *StdShellRepository) repositories.ShellRepository {
		return impl
	}); err != nil {
		return err
	}
	// Bind UpgradeShellRepository interface to implementation
	if err := container.Provide(func(impl *UpgradeAwareShellRepository) repositories.UpgradeShellRepository {
		return impl
	}); err != nil {
		return err
	}
	// Bind InteractiveShellRepository interface to implementation
	if err := container.Provide(func(impl *InteractiveShellRepository) repositories.InteractiveShellRepository {
		return impl
	}); err != nil {
		return err
	}
	return nil
}
