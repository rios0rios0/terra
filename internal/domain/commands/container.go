package commands

import (
	"go.uber.org/dig"
)

// RegisterProviders registers all command providers with the DIG container.
func RegisterProviders(container *dig.Container) error {
	// Register command constructors
	if err := container.Provide(NewDeleteCacheCommand); err != nil {
		return err
	}
	if err := container.Provide(NewFormatFilesCommand); err != nil {
		return err
	}
	if err := container.Provide(NewInstallDependenciesCommand); err != nil {
		return err
	}
	if err := container.Provide(NewRunAdditionalBeforeCommand); err != nil {
		return err
	}
	if err := container.Provide(NewRunFromRootCommand); err != nil {
		return err
	}
	if err := container.Provide(NewSelfUpdateCommand); err != nil {
		return err
	}
	if err := container.Provide(NewVersionCommand); err != nil {
		return err
	}

	// Bind interfaces to implementations
	if err := container.Provide(func(impl *DeleteCacheCommand) DeleteCache {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *FormatFilesCommand) FormatFiles {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *InstallDependenciesCommand) InstallDependencies {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *RunAdditionalBeforeCommand) RunAdditionalBefore {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *RunFromRootCommand) RunFromRoot {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *SelfUpdateCommand) SelfUpdate {
		return impl
	}); err != nil {
		return err
	}
	if err := container.Provide(func(impl *VersionCommand) Version {
		return impl
	}); err != nil {
		return err
	}

	return nil
}
